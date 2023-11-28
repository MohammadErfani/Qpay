package services

import (
	"Qpay/models"
	"bou.ke/monkey"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

var BA1 = models.BankAccount{
	UserID: 1,
	BankID: 2,
	Status: 1,
	Sheba:  "IR740170000000106748249001",
	Bank:   models.Bank{Name: "pasargad"},
}
var BA2 = models.BankAccount{
	UserID: 2,
	BankID: 1,
	Status: 1,
	Sheba:  "IR750170000000106748249001",
	Bank:   models.Bank{Name: "parsian"},
}

type BankAccountTestSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (suite *BankAccountTestSuite) SetupTest() {
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
func (suite *BankAccountTestSuite) TestGetSpecificBankAccount_UserIDAndBankAccountIDMatch() {

	userID := uint(1)
	bankAccountID := uint(2)
	ba1 := BA1
	ba1.Bank.ID = 2
	expectedBankAccount := ba1

	rows := sqlmock.NewRows([]string{"user_id", "bank_id", "status", "sheba"}).
		AddRow(expectedBankAccount.UserID, expectedBankAccount.BankID, expectedBankAccount.Status,
			expectedBankAccount.Sheba)

	suite.mock.ExpectQuery(`SELECT (.+) FROM "bank_accounts"`).
		WithArgs(bankAccountID, userID).
		WillReturnRows(rows)

	bankRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(2, expectedBankAccount.Bank.Name)

	suite.mock.ExpectQuery(`SELECT (.+) FROM "banks" WHERE "banks"."id" = \$1`).
		WithArgs(expectedBankAccount.Bank.ID).
		WillReturnRows(bankRows)

	bankAccount, err := GetSpecificBankAccount(suite.DB, userID, bankAccountID)

	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.Equal(expectedBankAccount, bankAccount, "Bank accounts should match")
}

func (suite *BankAccountTestSuite) TestGetSpecificBankAccount_UserIDAndBankAccountIDDoNotMatch() {
	userID := uint(1)
	bankAccountID := uint(1)
	ba1 := BA1
	ba1.Bank.ID = 2
	ba2 := BA2
	ba2.Bank.ID = 1

	sqlmock.NewRows([]string{"user_id", "bank_id", "status", "sheba"}).
		AddRow(ba1.UserID, ba1.BankID, ba1.Status,
			ba1.Sheba).AddRow(ba2.UserID, ba2.BankID, ba2.Status,
		ba2.Sheba)

	suite.mock.ExpectQuery(`SELECT (.+) FROM "bank_accounts"`).
		WithArgs(bankAccountID, userID).
		WillReturnError(gorm.ErrRecordNotFound)

	bankRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(2, ba1.Bank.Name).AddRow(1, ba2.Bank.Name)

	suite.mock.ExpectQuery(`SELECT (.+) FROM "banks" WHERE "banks"."id" = \$1`).
		WithArgs(ba1.Bank.ID).
		WillReturnRows(bankRows)

	bankAccount, err := GetSpecificBankAccount(suite.DB, userID, bankAccountID)

	require := suite.Require()
	require.Error(err, "expect an  error")
	require.Equal(models.BankAccount{}, bankAccount, "Bank account should be empty")
	require.Equal("bank Account Not found", err.Error(), "Error message should match")
}
func (suite *BankAccountTestSuite) TestCreateBankAccount() {
	userID := uint(1)
	sheba := "IR123456789012345678901234"

	// Patch SetUserAndBankForBankAccount
	monkey.Patch(SetUserAndBankForBankAccount, func(db *gorm.DB, userID uint,
		bankAccount *models.BankAccount) error {
		return nil
	})
	defer monkey.Unpatch(SetUserAndBankForBankAccount)
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "bank_accounts"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	suite.mock.ExpectCommit()

	createdBankAccount, err := CreateBankAccount(suite.DB, userID, sheba)
	createdBankAccount.ID = 1
	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(createdBankAccount, "Expected a non-nil bank account")
	require.NotZero(createdBankAccount.ID, "Expected a non-zero bank account ID")
	require.NoError(suite.mock.ExpectationsWereMet())
}

func TestBankAccountSuite(t *testing.T) {
	suite.Run(t, new(BankAccountTestSuite))
}
