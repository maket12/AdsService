package pg_test

import (
	"ads/authservice/internal/adapter/out/pg"
	"ads/authservice/internal/adapter/out/pg/sqlc"
	"ads/authservice/internal/domain/model"
	"ads/authservice/pkg/errs"
	"ads/migrations"
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type RefreshSessionsRepoSuite struct {
	suite.Suite
	db          *sql.DB
	repo        *pg.RefreshSessionsRepository
	ctx         context.Context
	dsn         string
	migrate     *migrate.Migrate
	testSession *model.RefreshSession
}

func TestRefreshSessionRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(RefreshSessionsRepoSuite))
}

func (s *RefreshSessionsRepoSuite) setupDatabase() {
	const targetVersion = 3

	sourceDriver, err := iofs.New(migrations.FS, ".")
	s.Require().NoError(err, "failed to create iofs driver")

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		sourceDriver,
		s.dsn,
	)
	s.Require().NoError(err, "failed to create migration instance")

	s.migrate = m

	err = m.Migrate(targetVersion)

	// If migration is correct - setup has done
	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		return
	}

	// Except dirty db as a normal scenario
	var dirtyErr migrate.ErrDirty
	if !errors.As(err, &dirtyErr) {
		s.FailNowf("failed to migrate up", "unexpected error: %v", err)
	}

	// ================ Restore dirty database ================
	_ = m.Force(dirtyErr.Version)

	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err, "failed to migrate down during recovery")
	}

	err = m.Migrate(targetVersion)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err, "failed to migrate up after recovery")
	}
}

func (s *RefreshSessionsRepoSuite) SetupSuite() {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "postgres://test:test@localhost:5432/testdb?sslmode=disable"
	}
	s.dsn = dsn

	pgClient, err := pg.NewPostgresClient(s.dsn, nil)
	s.Require().NoError(err)
	s.db = pgClient.DB

	s.setupDatabase()

	queries := sqlc.New(s.db)
	s.repo = pg.NewRefreshSessionsRepository(queries)

	s.ctx = context.Background()

	s.testSession = model.RestoreRefreshSession(
		uuid.New(),
		uuid.New(),
		"hashed-secret-token",
		time.Now(),
		time.Now().Add(time.Hour),
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	var testAccount = model.RestoreAccount(
		s.testSession.AccountID(),
		"new@email.com",
		"hashed-secret-pass",
		model.AccountActive,
		false,
		time.Now(),
		time.Now(),
		nil,
	)

	// Create an account in the main table
	accountsRepo := pg.NewAccountsRepository(queries)
	_ = accountsRepo.Create(s.ctx, testAccount)
}

func (s *RefreshSessionsRepoSuite) TearDownSuite() {
	if s.migrate != nil {
		if err := s.migrate.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			s.Require().NoError(err, "failed to migrate down")
		}
	}
	err := s.db.Close()
	s.Require().NoError(err, "failed to close db connection")
}

func (s *RefreshSessionsRepoSuite) SetupTest() {
	_, err := s.db.Exec("TRUNCATE TABLE refresh_sessions CASCADE")
	s.Require().NoError(err)
}

func (s *RefreshSessionsRepoSuite) TestCreateGetByID() {
	// Create at first
	err := s.repo.Create(s.ctx, s.testSession)
	s.Require().NoError(err)

	// Get by id
	session, err := s.repo.GetByID(s.ctx, s.testSession.ID())
	s.Require().NoError(err)
	s.Require().Equalf(s.testSession.AccountID(), session.AccountID(),
		"expected account id %v, got %v", s.testSession.AccountID(), session.AccountID())
	s.Require().Equalf(s.testSession.RefreshTokenHash(), session.RefreshTokenHash(),
		"expected token hash %v, got %v", s.testSession.RefreshTokenHash(), session.RefreshTokenHash())
}

func (s *RefreshSessionsRepoSuite) TestCreate_NonExistingAccount() {
	// Trying to create a session for non-existing account
	var anotherSession = model.RestoreRefreshSession(
		uuid.New(),
		uuid.New(),
		s.testSession.RefreshTokenHash(),
		s.testSession.CreatedAt(),
		s.testSession.ExpiresAt(),
		s.testSession.RevokedAt(),
		s.testSession.RevokeReason(),
		s.testSession.RotatedFrom(),
		s.testSession.IP(),
		s.testSession.UserAgent(),
	)
	err := s.repo.Create(s.ctx, anotherSession)
	s.Require().Error(err)
}

func (s *RefreshSessionsRepoSuite) TestCreate_DuplicateHash() {
	// Create a session
	_ = s.repo.Create(s.ctx, s.testSession)

	// Trying to create a session with the same token hash
	var anotherSession = model.RestoreRefreshSession(
		uuid.New(),
		s.testSession.AccountID(),
		s.testSession.RefreshTokenHash(),
		s.testSession.CreatedAt(),
		s.testSession.ExpiresAt(),
		s.testSession.RevokedAt(),
		s.testSession.RevokeReason(),
		s.testSession.RotatedFrom(),
		s.testSession.IP(),
		s.testSession.UserAgent(),
	)
	err := s.repo.Create(s.ctx, anotherSession)
	s.Require().Error(err)
}

func (s *RefreshSessionsRepoSuite) TestGetByHash() {
	// Create in advance
	_ = s.repo.Create(s.ctx, s.testSession)

	// Get by hash
	session, err := s.repo.GetByHash(s.ctx, s.testSession.RefreshTokenHash())
	s.Require().NoError(err)
	s.Require().Equalf(s.testSession.ID(), session.ID(),
		"expected token hash %v, got %v", s.testSession.ID(), session.ID())
	s.Require().Equalf(s.testSession.AccountID(), session.AccountID(),
		"expected account id %v, got %v", s.testSession.AccountID(), session.AccountID())
}

func (s *RefreshSessionsRepoSuite) TestGetByHash_NotFound() {
	// Get a non-existing session
	_, err := s.repo.GetByHash(s.ctx, s.testSession.RefreshTokenHash())
	s.Require().Error(err)
	s.Require().ErrorIsf(err, errs.ErrObjectNotFound,
		"expected \"ErrObjectNotFound\", got %v", err)
}

func (s *RefreshSessionsRepoSuite) TestRevoke() {
	// Create in advance
	_ = s.repo.Create(s.ctx, s.testSession)

	var (
		revokedSession = *s.testSession
		reason         = "account is blocked"
	)
	_ = revokedSession.Revoke(&reason)

	// Revoke the session
	err := s.repo.Revoke(s.ctx, &revokedSession)
	s.Require().NoError(err)

	// Ensure the session has been revoked
	session, _ := s.repo.GetByID(s.ctx, s.testSession.ID())
	s.Require().Equalf(revokedSession.RevokeReason(), session.RevokeReason(),
		"expected revoke reason %v, got %v", revokedSession.RevokeReason(), session.RevokeReason())
}

func (s *RefreshSessionsRepoSuite) TestRevokeAllForAccount() {
	var anotherSession = model.RestoreRefreshSession(
		uuid.New(),
		s.testSession.AccountID(),
		"hashed",
		s.testSession.CreatedAt(),
		s.testSession.ExpiresAt(),
		s.testSession.RevokedAt(),
		s.testSession.RevokeReason(),
		s.testSession.RotatedFrom(),
		s.testSession.IP(),
		s.testSession.UserAgent(),
	)

	// Create some sessions for the same account
	_ = s.repo.Create(s.ctx, s.testSession)
	_ = s.repo.Create(s.ctx, anotherSession)

	var reason = "tests"

	err := s.repo.RevokeAllForAccount(s.ctx, s.testSession.AccountID(), &reason)
	s.Require().NoError(err)

	// Ensure all sessions have been revoked
	sess, _ := s.repo.GetByID(s.ctx, s.testSession.ID())
	s.Require().Equalf(reason, *sess.RevokeReason(),
		"expected revoke reason %v, got %v", reason, *sess.RevokeReason())

	sess, _ = s.repo.GetByID(s.ctx, anotherSession.ID())
	s.Require().Equalf(reason, *sess.RevokeReason(),
		"expected revoke reason %v, got %v", reason, *sess.RevokeReason())
}

func (s *RefreshSessionsRepoSuite) TestRevokeDescendants() {
	// Create sessions - one is the descendant of the second
	var (
		rotatedID      = s.testSession.ID()
		anotherSession = model.RestoreRefreshSession(
			uuid.New(),
			s.testSession.AccountID(),
			"hashed",
			s.testSession.CreatedAt(),
			s.testSession.ExpiresAt(),
			s.testSession.RevokedAt(),
			s.testSession.RevokeReason(),
			&rotatedID,
			s.testSession.IP(),
			s.testSession.UserAgent(),
		)
		reason = "test revoke"
	)
	_ = s.repo.Create(s.ctx, s.testSession)
	_ = s.repo.Create(s.ctx, anotherSession)

	err := s.repo.RevokeDescendants(s.ctx, s.testSession.ID(), &reason)
	s.Require().NoError(err)

	// Ensure the session has been revoked
	session, _ := s.repo.GetByHash(s.ctx, anotherSession.RefreshTokenHash())
	s.Require().Equalf(reason, *session.RevokeReason(),
		"expected revoke reason %v, got %v", reason, *session.RevokeReason())
}

func (s *RefreshSessionsRepoSuite) TestDeleteExpired() {
	// Create a session
	_ = s.repo.Create(s.ctx, s.testSession)

	// Delete expired (set time what is much later)
	var expiresAt = time.Now().Add(time.Hour)
	err := s.repo.DeleteExpired(s.ctx, expiresAt)
	s.Require().NoError(err)

	// Ensure it was deleted
	_, err = s.repo.GetByID(s.ctx, s.testSession.ID())
	s.Require().Error(err)
	s.Require().ErrorIs(err, errs.ErrObjectNotFound)
}

func (s *RefreshSessionsRepoSuite) TestListActiveForAccount() {
	const sessionsAmount = 2

	var anotherSession = model.RestoreRefreshSession(
		uuid.New(),
		s.testSession.AccountID(),
		"hashed",
		s.testSession.CreatedAt(),
		s.testSession.ExpiresAt(),
		s.testSession.RevokedAt(),
		s.testSession.RevokeReason(),
		s.testSession.RotatedFrom(),
		s.testSession.IP(),
		s.testSession.UserAgent(),
	)

	// Create sessions
	_ = s.repo.Create(s.ctx, s.testSession)
	_ = s.repo.Create(s.ctx, anotherSession)

	// List of active
	sessions, err := s.repo.ListActiveForAccount(s.ctx, s.testSession.AccountID())
	s.Require().NoError(err)
	s.Require().Len(sessions, sessionsAmount)

	var fstFound, sndFound bool
	for i := range sessions {
		value := *sessions[i]
		if value.ID() == s.testSession.ID() {
			fstFound = true
		}
		if value.ID() == anotherSession.ID() {
			sndFound = true
		}
	}

	s.Require().Truef(fstFound, "expected %v\n in %v",
		s.testSession, sessions)
	s.Require().Truef(sndFound, "expected %v\n in %v",
		anotherSession, sessions)
}
