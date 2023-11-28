package services

import (
	"Qpay/models"
	"bou.ke/monkey"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type GatewaySuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (suite *GatewaySuite) SetupTest() {
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
func (suite *GatewaySuite) TestGetSpecificGateway_Success() {
	userID := uint(1)
	gatewayID := uint(2)
	expectedGateway := models.Gateway{

		UserID: userID,
		User:   models.User{Name: "John Doe"},
	}
	expectedGateway.User.ID = 1
	expectedGateway.ID = 2
	userRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "John Doe") // Assuming the user with ID 1 exists

	rows := sqlmock.NewRows([]string{"id", "user_id"}).
		AddRow(expectedGateway.ID, expectedGateway.UserID)

	suite.mock.ExpectQuery(`SELECT (.+) FROM "gateways" WHERE \(id=\$1 AND user_id=\$2\)`).
		WithArgs(gatewayID, userID).
		WillReturnRows(rows)
	suite.mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE "users"."id" = \$1`).
		WithArgs(userID).
		WillReturnRows(userRows)
	resultGateway, err := GetSpecificGateway(suite.DB, userID, gatewayID)

	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.Equal(expectedGateway, resultGateway, "Gateways should match")
}
func (suite *GatewaySuite) TestGetSpecificGateway_Fail() {
	userID := uint(1)
	gatewayID := uint(1)

	suite.mock.ExpectQuery(`SELECT (.+) FROM "gateways" WHERE \(id=\$1 AND user_id=\$2\)`).
		WithArgs(gatewayID, userID).
		WillReturnError(gorm.ErrRecordNotFound)

	resultGateway, err := GetSpecificGateway(suite.DB, userID, gatewayID)

	require := suite.Require()
	require.Error(err, "Expect an error")
	require.Equal(models.Gateway{}, resultGateway, "Gateway should be empty")
	require.Equal("gateway Not found", err.Error(), "Error message should match")
}
func (suite *GatewaySuite) TestCreateGateway_Success() {
	userID := uint(1)
	name := "Gateway1"
	logo := "Logo1"
	bankAccountID := uint(2)
	commissionID := uint(3)
	isPersonal := false

	monkey.Patch(CheckUserAndBankAccountForGateway, func(db *gorm.DB, userID, bankAccountID uint) (*models.User, error) {
		return &models.User{}, nil
	})
	defer monkey.Unpatch(CheckUserAndBankAccountForGateway)

	monkey.Patch(CheckPersonal, func(db *gorm.DB, user *models.User) bool {
		return false
	})
	defer monkey.Unpatch(CheckPersonal)

	monkey.Patch(CheckCommission, func(db *gorm.DB, commissionID uint) error {
		return nil
	})
	defer monkey.Unpatch(CheckCommission)

	monkey.Patch(SetDefaultRoute, func(db *gorm.DB, user *models.User, gateway *models.Gateway) {
	})
	defer monkey.Unpatch(SetDefaultRoute)

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "gateways"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.mock.ExpectCommit()

	createdGateway, err := CreateGateway(suite.DB, userID, name, logo, bankAccountID, commissionID, isPersonal)
	createdGateway.ID = 1
	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(createdGateway, "Expected a non-nil gateway")
	require.NotZero(createdGateway.ID, "Expected a non-zero gateway ID")
	require.NoError(suite.mock.ExpectationsWereMet())
}
func (suite *GatewaySuite) TestCreateGateway_Fail1() {
	userID := uint(1)
	name := "Gateway1"
	logo := "Logo1"
	bankAccountID := uint(2)
	commissionID := uint(3)
	isPersonal := false

	monkey.Patch(CheckUserAndBankAccountForGateway, func(db *gorm.DB, userID, bankAccountID uint) (*models.User, error) {
		return nil, errors.New("UnAuthorize")
	})
	defer monkey.Unpatch(CheckUserAndBankAccountForGateway)

	createdGateway, err := CreateGateway(suite.DB, userID, name, logo, bankAccountID, commissionID, isPersonal)
	require := suite.Require()
	require.Error(err, "Expected an error")
	require.Nil(createdGateway, "Expected a nil gateway")
	require.EqualError(err, "UnAuthorize", "Error message should match")
}
func (suite *GatewaySuite) TestCreateGateway_Fail2() {
	userID := uint(1)
	name := "Gateway1"
	logo := "Logo1"
	bankAccountID := uint(2)
	commissionID := uint(3)
	isPersonal := false
	monkey.Patch(CheckUserAndBankAccountForGateway, func(db *gorm.DB, userID, bankAccountID uint) (*models.User, error) {
		return &models.User{}, nil
	})
	defer monkey.Unpatch(CheckUserAndBankAccountForGateway)

	monkey.Patch(CheckPersonal, func(db *gorm.DB, user *models.User) bool {
		return false
	})
	defer monkey.Unpatch(CheckPersonal)

	monkey.Patch(CheckCommission, func(db *gorm.DB, commissionID uint) error {
		return nil
	})
	defer monkey.Unpatch(CheckCommission)

	monkey.Patch(SetDefaultRoute, func(db *gorm.DB, user *models.User, gateway *models.Gateway) {
	})
	defer monkey.Unpatch(SetDefaultRoute)
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "gateways"`)).
		WillReturnError(errors.New("internal server error"))
	suite.mock.ExpectRollback()

	createdGateway, err := CreateGateway(suite.DB, userID, name, logo, bankAccountID, commissionID, isPersonal)
	require := suite.Require()
	require.Error(err, "Expected an error")
	require.Nil(createdGateway, "Expected a nil gateway")
	require.EqualError(err, "internal server error", "Error message should match")
}
func (suite *GatewaySuite) TestPurchaseAddress_Success() {
	userID := uint(1)
	gatewayID := uint(2)
	route := "new_route"
	gateway := models.Gateway{
		UserID: userID,
		// set other fields as needed
	}
	// Monkey patch GetSpecificGateway
	monkey.Patch(GetSpecificGateway, func(db *gorm.DB, userID, gatewayID uint) (models.Gateway, error) {
		return gateway, nil
	})
	defer monkey.Unpatch(GetSpecificGateway)

	monkey.Patch(GetGateway, func(db *gorm.DB, fieldName, fieldValue string) (*models.Gateway, error) {
		return nil, errors.New("not found")
	})
	defer monkey.Unpatch(GetGateway)
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(`^UPDATE "gateways" SET (.+)`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mock.ExpectCommit()
	resultGateway, err := PurchaseAddress(suite.DB, userID, gatewayID, route)
	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(resultGateway, "Expected a non-nil gateway")

}

func TestGatewaySuite(t *testing.T) {
	suite.Run(t, new(GatewaySuite))
}
