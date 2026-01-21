package postgres_test

import (
	"ads/authservice/internal/adapter/out/postgres"
	"ads/authservice/internal/domain/model"
	"ads/authservice/migrations"
	"ads/pkg/errs"
	"ads/pkg/postgres"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AccountRolesRepoSuite struct {
	suite.Suite
	dbClient *postgres.Client
	repo     *postgres.AccountRoleRepository
	ctx      context.Context
	migrate  *migrate.Migrate
	testRole *model.AccountRole
}

func TestAccountRolesRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(AccountRolesRepoSuite))
}

func (s *AccountRolesRepoSuite) setupDatabase() {
	dbConfig := postgres.NewConfig(
		"localhost", 5432,
		"test", "test", "testdb",
		"disable", 25, 25, time.Minute*5,
	)
	dsn := "postgres://test:test@localhost:5432/testdb?sslmode=disable"

	dbClient, err := postgres.NewClient(dbConfig)
	s.Require().NoError(err)
	s.dbClient = dbClient

	sourceDriver, err := iofs.New(migrations.FS, ".")
	s.Require().NoError(err, "failed to create iofs driver")

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		sourceDriver,
		dsn,
	)
	s.Require().NoError(err, "failed to create migration instance")

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

func (s *AccountRolesRepoSuite) SetupSuite() {
	s.setupDatabase()
	s.repo = postgres.NewAccountRolesRepository(s.dbClient)
	s.ctx = context.Background()
	testAccount, _ := model.NewAccount("new@email.com", "hashed-secret-pass")

	// Create an account in the main table
	accountsRepo := postgres.NewAccountsRepository(s.dbClient)
	_ = accountsRepo.Create(s.ctx, testAccount)

	s.testRole, _ = model.NewAccountRole(testAccount.ID())
}

func (s *AccountRolesRepoSuite) TearDownSuite() {
	if s.migrate != nil {
		if err := s.migrate.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			s.Require().NoError(err, "failed to migrate down")
		}
	}
	err := s.dbClient.Close()
	s.Require().NoError(err, "failed to close db connection")
}

func (s *AccountRolesRepoSuite) SetupTest() {
	_, err := s.dbClient.DB.Exec("TRUNCATE TABLE account_roles CASCADE")
	s.Require().NoError(err)
}

func (s *AccountRolesRepoSuite) TestCreateGet() {
	// Create at first
	err := s.repo.Create(s.ctx, s.testRole)
	s.Require().NoError(err)

	// Get by account id
	role, err := s.repo.Get(s.ctx, s.testRole.AccountID())
	s.Require().NoError(err)
	s.Require().Equal(s.testRole.Role(), role.Role())
}

func (s *AccountRolesRepoSuite) TestCreate_NonExistingAccount() {
	// Create an account role for unexisting account
	newRole, _ := model.NewAccountRole(uuid.New())
	err := s.repo.Create(s.ctx, newRole)
	s.Require().Error(err)
}

func (s *AccountRolesRepoSuite) TestGet_NotFound() {
	// Try to get non-existing account role
	_, err := s.repo.Get(s.ctx, s.testRole.AccountID())
	s.Require().Error(err)
	s.Require().ErrorIs(err, errs.ErrObjectNotFound)
}

func (s *AccountRolesRepoSuite) TestUpdate() {
	// Create at first
	_ = s.repo.Create(s.ctx, s.testRole)

	// Copy value and assigned to not change test data
	assignedRole := *s.testRole
	_ = assignedRole.Assign("admin")

	err := s.repo.Update(s.ctx, &assignedRole)
	s.Require().NoError(err)

	// Ensure update was correct
	acc, _ := s.repo.Get(s.ctx, s.testRole.AccountID())
	s.Require().Equal(model.RoleAdmin, acc.Role())
}

func (s *AccountRolesRepoSuite) TestDelete() {
	// Create at first
	_ = s.repo.Create(s.ctx, s.testRole)

	// Delete
	err := s.repo.Delete(s.ctx, s.testRole.AccountID())
	s.Require().NoError(err)

	// Ensure deletion was successful
	_, err = s.repo.Get(s.ctx, s.testRole.AccountID())
	s.Require().Error(err)
	s.Require().ErrorIs(err, errs.ErrObjectNotFound)
}
