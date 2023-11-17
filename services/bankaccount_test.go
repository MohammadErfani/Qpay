package services

import (
	"Qpay/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	BankID: 2,
	Status: 1,
	Sheba:  "IR750170000000106748249001",
	Bank:   models.Bank{Name: "pasargad"},
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
	
}

func TestBankAccountSuite(t *testing.T) {
	suite.Run(t, new(BankAccountTestSuite))
}
