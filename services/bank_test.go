package services

import (
	"Qpay/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type BankSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (suite *BankSuite) SetupTest() {
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
func (suite *BankSuite) TestGetBank_Success() {
	fieldName := "name"
	fieldValue := "Example Bank"
	expectedBank := models.Bank{
		Name: fieldValue,
	}
	expectedBank.ID = 1
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(expectedBank.ID, expectedBank.Name)

	suite.mock.ExpectQuery(`SELECT (.+) FROM "banks" WHERE name=\$1`).
		WithArgs(fieldValue).
		WillReturnRows(rows)

	resultBank, err := GetBank(suite.DB, fieldName, fieldValue)

	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.Equal(&expectedBank, resultBank, "Banks should match")

	require.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *BankSuite) TestGetBank_NotFound() {
	fieldName := "name"
	fieldValue := "Nonexistent Bank"

	suite.mock.ExpectQuery(`SELECT (.+) FROM "banks" WHERE name=\$1`).
		WithArgs(fieldValue).
		WillReturnError(gorm.ErrRecordNotFound)

	resultBank, err := GetBank(suite.DB, fieldName, fieldValue)

	require := suite.Require()
	require.Error(err, "Expected an error")
	require.Nil(resultBank, "Bank should be nil")
	require.Equal("bank not found", err.Error(), "Error message should match")
	require.NoError(suite.mock.ExpectationsWereMet())
}
func TestBankSuite(t *testing.T) {
	suite.Run(t, new(BankSuite))
}
