package services

import (
	"Qpay/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type AdminSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (suite *AdminSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.FailNow("Failed to create mock database", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	}))
	if err != nil {
		suite.FailNow("Failed to open GORM DB", err)
	}
	suite.DB = gormDB
	suite.mock = mock
}
func (suite AdminSuite) TestCreateAdmin_Success() {
	// Mock the successful case where admin is created
	expectedAdmin := &models.User{
		Name:     "John Doe",
		Username: "john_doe",
		Email:    "john.doe@example.com",
		Role:     models.IsAdmin,
	}
	expectedAdmin.ID = 1

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.mock.ExpectCommit()

	createdAdmin, err := CreateAdmin(suite.DB, expectedAdmin.Name, expectedAdmin.Username, expectedAdmin.Email, "password123")
	createdAdmin.ID = 1
	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(createdAdmin, "Expected a non-nil admin")

	require.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *AdminSuite) TestCheckIsAdmin_Success() {
	userID := uint(1)

	// Mock the successful case where user is an admin
	suite.mock.ExpectQuery(`SELECT "role" FROM "users" WHERE "id" = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow(models.IsAdmin))

	err := CheckIsAdmin(suite.DB, userID)
	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NoError(suite.mock.ExpectationsWereMet())
}
func (suite *AdminSuite) TestCheckIsAdmin_Failure() {
	userID := uint(1)

	suite.mock.ExpectQuery(`SELECT "role" FROM "users" WHERE "id" = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow(models.IsNaturalPerson))

	err := CheckIsAdmin(suite.DB, userID)

	require := suite.Require()
	require.Error(err, "Expect an error")
	require.Equal("unAuthorize", err.Error(), "Error message should match")

	require.NoError(suite.mock.ExpectationsWereMet())
}
func TestAdminSuite(t *testing.T) {
	suite.Run(t, new(AdminSuite))
}
