package services

import (
	"Qpay/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type UserSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

var UserModel = models.User{
	Name:        "John Doe",
	Username:    "johndoe",
	Email:       "john@example.com",
	Password:    "hashedpassword",
	PhoneNumber: "09125330680",
	Identity:    "1234567890",
	Address:     "123 Main St, City",
	Role:        1,
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
	user := UserModel

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
	expectedUser := &UserModel

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
func (suite *UserSuite) TestGetUserByPhoneNumber_UserExist() {
	phoneNumber := "09125330680"
	expectedUser := &UserModel

	rows := sqlmock.NewRows([]string{"name", "username", "email", "password", "phone_number", "identity", "address", "role"}).
		AddRow(expectedUser.Name, expectedUser.Username, expectedUser.Email,
			expectedUser.Password, expectedUser.PhoneNumber, expectedUser.Identity,
			expectedUser.Address, expectedUser.Role)
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE phone_number=\$1`).
		WithArgs(phoneNumber).
		WillReturnRows(rows)

	user, err := GetUser(suite.DB, "phone_number", phoneNumber)

	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(user, "User should not be nil")
	require.Equal(expectedUser, user, "Users should match")
}
func (suite *UserSuite) TestGetUserByPhoneNumber_UserNotFound() {
	phoneNumber := "09125330680"

	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE phone_number=\$1`).
		WithArgs(phoneNumber).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := GetUser(suite.DB, "phone_number", phoneNumber)

	require := suite.Require()
	require.Error(err, "Expected an error")
	require.Nil(user, "User should be nil")
	require.Equal("user not found", err.Error(), "Error message should match")
}
func (suite *UserSuite) TestGetUserByUserName_UserExist() {
	userName := "johndoe"
	expectedUser := &UserModel

	rows := sqlmock.NewRows([]string{"name", "username", "email", "password", "phone_number", "identity", "address", "role"}).
		AddRow(expectedUser.Name, expectedUser.Username, expectedUser.Email,
			expectedUser.Password, expectedUser.PhoneNumber, expectedUser.Identity,
			expectedUser.Address, expectedUser.Role)
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE username=\$1`).
		WithArgs(userName).
		WillReturnRows(rows)

	user, err := GetUser(suite.DB, "username", userName)

	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(user, "User should not be nil")
	require.Equal(expectedUser, user, "Users should match")
}
func (suite *UserSuite) TestGetUserByUserName_UserNotFound() {
	userName := "johndoe"

	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE username=\$1`).
		WithArgs(userName).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := GetUser(suite.DB, "username", userName)

	require := suite.Require()
	require.Error(err, "Expected an error")
	require.Nil(user, "User should be nil")
	require.Equal("user not found", err.Error(), "Error message should match")
}
func (suite *UserSuite) TestGetUserByIdentity_UserExist() {
	Identity := "1234567890"
	expectedUser := &UserModel

	rows := sqlmock.NewRows([]string{"name", "username", "email", "password", "phone_number", "identity", "address", "role"}).
		AddRow(expectedUser.Name, expectedUser.Username, expectedUser.Email,
			expectedUser.Password, expectedUser.PhoneNumber, expectedUser.Identity,
			expectedUser.Address, expectedUser.Role)
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE identity=\$1`).
		WithArgs(Identity).
		WillReturnRows(rows)

	user, err := GetUser(suite.DB, "identity", Identity)

	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(user, "User should not be nil")
	require.Equal(expectedUser, user, "Users should match")
}
func (suite *UserSuite) TestGetUserByIdentity_UserNotFound() {
	Identity := "1234567890"

	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE identity=\$1`).
		WithArgs(Identity).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := GetUser(suite.DB, "identity", Identity)

	require := suite.Require()
	require.Error(err, "Expected an error")
	require.Nil(user, "User should be nil")
	require.Equal("user not found", err.Error(), "Error message should match")
}
func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
