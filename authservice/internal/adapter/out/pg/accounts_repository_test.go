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
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"

	"github.com/stretchr/testify/suite"
)

type AccountsRepoSuite struct {
	suite.Suite
	db          *sql.DB
	repo        *pg.AccountsRepository
	ctx         context.Context
	dsn         string
	migrate     *migrate.Migrate
	testAccount *model.Account
}

func TestAccountRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(AccountsRepoSuite))
}

func (s *AccountsRepoSuite) setupDatabase() {
	sourceDriver, err := iofs.New(migrations.FS, ".")
	s.Require().NoError(err, "failed to create iofs driver")

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		sourceDriver,
		s.dsn,
	)
	s.Require().NoError(err, "failed to crate migration instance")

	s.migrate = m

	err = m.Up()

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

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		s.Require().NoError(err, "failed to migrate up after recovery")
	}
}

func (s *AccountsRepoSuite) SetupSuite() {
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
	s.repo = pg.NewAccountsRepository(queries)

	s.ctx = context.Background()

	s.testAccount = model.RestoreAccount(
		uuid.New(),
		"new@email.com",
		"hashed-secret-pass",
		model.AccountActive,
		false,
		time.Now(),
		time.Now(),
		nil,
	)
}

func (s *AccountsRepoSuite) TearDownSuite() {
	if s.migrate != nil {
		if err := s.migrate.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			s.Require().NoError(err, "failed to migrate down")
		}
	}
	err := s.db.Close()
	s.Require().NoError(err, "failed to close db connection")
}

func (s *AccountsRepoSuite) SetupTest() {
	_, err := s.db.Exec("TRUNCATE TABLE accounts CASCADE")
	s.Require().NoError(err)
}

func (s *AccountsRepoSuite) TestCreateGetByID() {
	// Check create first
	err := s.repo.Create(s.ctx, s.testAccount)
	s.Require().NoError(err)

	// And then get
	acc, err := s.repo.GetByID(s.ctx, s.testAccount.ID())
	s.Require().NoError(err)
	s.Require().Exactlyf(s.testAccount.Email(), acc.Email(),
		"Expected email %v, got %v", s.testAccount.Email(), acc.Email())
	s.Require().Exactlyf(s.testAccount.PasswordHash(), acc.PasswordHash(),
		"Expected pass has %v, got %v", s.testAccount.PasswordHash(), acc.PasswordHash())
}

func (s *AccountsRepoSuite) TestCreate_DuplicateEmail() {
	// Create an account at first
	_ = s.repo.Create(s.ctx, s.testAccount)

	// Trying to create an account with the same email
	var newAcc = model.RestoreAccount(
		uuid.New(), s.testAccount.Email(),
		"hashed-pass", model.AccountActive, true,
		time.Now(), time.Now(), nil,
	)
	err := s.repo.Create(s.ctx, newAcc)
	s.Require().Error(err)
}

func (s *AccountsRepoSuite) TestGetByEmail() {
	// Create an account in advance
	_ = s.repo.Create(s.ctx, s.testAccount)

	// Get by email
	acc, err := s.repo.GetByEmail(s.ctx, s.testAccount.Email())
	s.Require().NoError(err)
	s.Require().Exactlyf(s.testAccount.ID(), acc.ID(),
		"Expected id %v, got %v", s.testAccount.ID(), acc.ID())
	s.Require().Exactlyf(s.testAccount.PasswordHash(), acc.PasswordHash(),
		"Expected pass has %v, got %v", s.testAccount.PasswordHash(), acc.PasswordHash())
}

func (s *AccountsRepoSuite) TestGetByEmail_CaseInsensitive() {
	// Create an account in advance
	_ = s.repo.Create(s.ctx, s.testAccount)

	// Trying to get by the same email, but in upper case
	var upperEmail = strings.ToUpper(s.testAccount.Email())
	acc, err := s.repo.GetByEmail(s.ctx, upperEmail)

	s.Require().NoError(err)
	s.Require().Equalf(s.testAccount.ID(), acc.ID(),
		"expected id %v, got %v", s.testAccount.ID(), acc.ID())
}

func (s *AccountsRepoSuite) TestGetByEmail_NotFound() {
	// Trying to get non-existing account
	var unexistingEmail = "unexist@gmail.com"
	_, err := s.repo.GetByEmail(s.ctx, unexistingEmail)
	s.Require().ErrorIsf(err, errs.ErrObjectNotFound,
		"Expected error \"ErrObjectNotFound\", got %v", err)
}

func (s *AccountsRepoSuite) TestMarkLogin() {
	// Create account at first
	_ = s.repo.Create(s.ctx, s.testAccount)

	// Mark as logged in
	err := s.repo.MarkLogin(s.ctx, s.testAccount)
	s.Require().NoError(err)

	// Check if the account is marked
	acc, _ := s.repo.GetByEmail(s.ctx, s.testAccount.Email())
	s.Require().NotNil(acc.LastLoginAt(), "expected not null last log in time")

	// Check update time
	s.Require().NotEqual(s.testAccount.UpdatedAt(), acc.UpdatedAt(),
		"expected new update time")
}

func (s *AccountsRepoSuite) TestVerifyEmail() {
	// Create account at first
	_ = s.repo.Create(s.ctx, s.testAccount)

	// Verify its email
	err := s.repo.VerifyEmail(s.ctx, s.testAccount)
	s.Require().NoError(err)

	// Check if the account is marked
	acc, _ := s.repo.GetByID(s.ctx, s.testAccount.ID())
	s.Require().Truef(acc.EmailVerified(), "expected verified email, got non-verified")

	// Check update time
	s.Require().NotEqual(s.testAccount.UpdatedAt(), acc.UpdatedAt(),
		"expected new update time")
}
