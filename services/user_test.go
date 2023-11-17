package services

import (
	"Qpay/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

// ReturnResult(sqlmock.NewResult(1, 1)) is saying that when the SQL statement
// specified in ExpectExec is executed, it should return a result indicating that
// it affected 1 row and the last insert ID is 1. This is a way to simulate the
// result of a database operation for testing purposes.
type UserSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (suite *UserSuite) SetupTest() {
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
func (suite *UserSuite) TestCreateUser() {
	// Set up expectations for the SQL mock
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(`INSERT`).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.mock.ExpectCommit()

	// Set up user data for testing
	user := models.User{
		// Populate with necessary user fields
		Name:        "John Doe",
		Username:    "johndoe",
		Email:       "john@example.com",
		Password:    "hashedpassword",
		PhoneNumber: "123456789",
		Identity:    "1234567890",
		Address:     "123 Main St, City",
		Role:        1,
	}

	// Call the CreateUser function
	createdUser, err := CreateUser(suite.DB, user)

	// Assertions
	suite.NoError(err)
	suite.NotNil(createdUser)
	err = suite.mock.ExpectationsWereMet()
	suite.NoError(err)
}
func (suite *UserSuite) TestGetUserByEmail_UserExists() {
	email := "john@example.com"
	expectedUser := &models.User{
		Name:        "John Doe",
		Username:    "johndoe",
		Email:       "john@example.com",
		Password:    "hashedpassword",
		PhoneNumber: "123456789",
		Identity:    "1234567890",
		Address:     "123 Main St, City",
		Role:        1,
	}

	rows := sqlmock.NewRows([]string{"name", "username", "email", "password", "phone_number", "identity", "address", "role"}).
		AddRow(expectedUser.Name, expectedUser.Username, expectedUser.Email,
			expectedUser.Password, expectedUser.PhoneNumber, expectedUser.Identity,
			expectedUser.Address, expectedUser.Role)

	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE email=\$1`).
		WithArgs(email).
		WillReturnRows(rows)

	// Call the function being tested
	user, err := GetUserByEmail(suite.DB, email)

	// Assertions using require
	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(user, "User should not be nil")
	require.Equal(expectedUser, user, "Users should match")
}
func (suite *UserSuite) TestGetUserByEmail_UserNotFound() {
	email := "nonexistent@example.com"

	// Mock database response for a user not found
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE email=\$1`).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)

	// Call the function being tested
	user, err := GetUserByEmail(suite.DB, email)

	// Assertions using require
	require := suite.Require()
	require.Error(err, "Expected an error")
	require.Nil(user, "User should be nil")
	require.Equal("user not found", err.Error(), "Error message should match")
}
func TestSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
