package services

import (
	"Qpay/models"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type CommSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (suite *CommSuite) SetupTest() {
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
func (suite *CommSuite) TestCreateCommission_Success() {
	amountPerTrans := 10.0
	percentPerTrans := 2.5
	expectedCommission := &models.Commission{
		AmountPerTrans:     amountPerTrans,
		PercentagePerTrans: percentPerTrans,
		Status:             models.CommIsActive,
	}
	expectedCommission.ID = 1
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "commissions"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()
	createdCommission, err := CreateCommission(suite.DB, amountPerTrans, percentPerTrans)
	expectedCommission.ID = 1
	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(createdCommission, "Expected a non-nil commission")
	require.NoError(suite.mock.ExpectationsWereMet())
}
func (suite *CommSuite) TestCreateCommission_Error() {
	amountPerTrans := 10.0
	percentPerTrans := 2.5

	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "commissions"`)).
		WillReturnError(errors.New("some database error"))
	suite.mock.ExpectRollback()
	createdCommission, err := CreateCommission(suite.DB, amountPerTrans, percentPerTrans)
	require := suite.Require()
	require.Error(err, "Expect an error")
	require.Nil(createdCommission, "Expected a nil commission")
	require.EqualError(err, "some database error", "Error message should match")
	require.NoError(suite.mock.ExpectationsWereMet())
}
func (suite *CommSuite) TestGetCommission_Success() {
	fieldName := "someField"
	fieldValue := "someValue"

	expectedCommission := &models.Commission{
		AmountPerTrans:     10.0,
		PercentagePerTrans: 2.5,
		Status:             models.CommIsActive,
	}
	expectedCommission.ID = 1

	commissionRows := sqlmock.NewRows([]string{"id", "amount_per_trans", "percentage_per_trans", "status"}).
		AddRow(expectedCommission.ID, expectedCommission.AmountPerTrans, expectedCommission.PercentagePerTrans, expectedCommission.Status)
	suite.mock.ExpectQuery(`SELECT (.+) FROM "commissions" WHERE someField=\$1`).
		WithArgs(fieldValue).
		WillReturnRows(commissionRows)

	resultCommission, err := GetCommission(suite.DB, fieldName, fieldValue)

	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.NotNil(resultCommission, "Expected a non-nil commission")
	require.Equal(expectedCommission, resultCommission, "Commissions should match")

	require.NoError(suite.mock.ExpectationsWereMet())
}
func (suite *CommSuite) TestGetCommission_Error() {
	fieldName := "someField"
	fieldValue := "someValue"

	suite.mock.ExpectQuery(`SELECT (.+) FROM "commissions" WHERE someField=\$1`).
		WithArgs(fieldValue).
		WillReturnError(errors.New("commission not found"))

	resultCommission, err := GetCommission(suite.DB, fieldName, fieldValue)

	require := suite.Require()
	require.Error(err, "Expect an error")
	require.Nil(resultCommission, "Expected a nil commission")
	require.EqualError(err, "commission not found", "Error message should match")

	require.NoError(suite.mock.ExpectationsWereMet())
}
func (suite *CommSuite) TestListActiveCommission_Success() {
	expectedCommissions := []models.Commission{
		{AmountPerTrans: 10.0, PercentagePerTrans: 2.5, Status: models.CommIsActive},
		{AmountPerTrans: 15.0, PercentagePerTrans: 3.0, Status: models.CommIsActive},
	}
	expectedCommissions[0].ID = 1
	expectedCommissions[1].ID = 2

	commissionRows := sqlmock.NewRows([]string{"id", "amount_per_trans", "percentage_per_trans", "status"}).
		AddRow(expectedCommissions[0].ID, expectedCommissions[0].AmountPerTrans, expectedCommissions[0].PercentagePerTrans, expectedCommissions[0].Status).
		AddRow(expectedCommissions[1].ID, expectedCommissions[1].AmountPerTrans, expectedCommissions[1].PercentagePerTrans, expectedCommissions[1].Status)
	suite.mock.ExpectQuery(`SELECT (.+) FROM "commissions" WHERE status=\$1`).
		WithArgs(models.CommIsActive).
		WillReturnRows(commissionRows)

	resultCommissions, err := ListActiveCommission(suite.DB)

	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.Equal(expectedCommissions, resultCommissions, "Commissions should match")

	require.NoError(suite.mock.ExpectationsWereMet())
}
func (suite *CommSuite) TestListActiveCommission_Error() {

	suite.mock.ExpectQuery(`SELECT (.+) FROM "commissions" WHERE status=\$1`).
		WithArgs(models.CommIsActive).
		WillReturnError(errors.New("error getting commissions"))
	resultCommissions, err := ListActiveCommission(suite.DB)
	require := suite.Require()
	require.Error(err, "Expect an error")
	require.Nil(resultCommissions, "Expected a nil commission slice")
	require.EqualError(err, "error getting commissions", "Error message should match")

	require.NoError(suite.mock.ExpectationsWereMet())
}
func (suite *CommSuite) TestListAllCommission_Success() {
	expectedCommissions := []models.Commission{
		{
			AmountPerTrans:     10.0,
			PercentagePerTrans: 2.5,
			Status:             models.CommIsActive,
			Gateways: []models.Gateway{
				{CommissionID: 1},
				{CommissionID: 1},
			},
		},
		{
			AmountPerTrans:     15.0,
			PercentagePerTrans: 3.0,
			Status:             models.CommIsActive,
			Gateways: []models.Gateway{
				{CommissionID: 2},
			},
		},
	}
	expectedCommissions[0].Gateways[0].ID = 1
	expectedCommissions[0].Gateways[1].ID = 2
	expectedCommissions[1].Gateways[0].ID = 3
	expectedCommissions[0].ID = 1
	expectedCommissions[1].ID = 2

	commissionRows := sqlmock.NewRows([]string{"id", "amount_per_trans", "percentage_per_trans", "status"}).
		AddRow(1, 10.0, 2.5, models.CommIsActive).
		AddRow(2, 15.0, 3.0, models.CommIsActive)

	gatewayRows := sqlmock.NewRows([]string{"id", "commission_id"}).
		AddRow(1, 1).
		AddRow(2, 1).
		AddRow(3, 2)
	suite.mock.ExpectQuery(`SELECT (.+) FROM "commissions"`).
		WillReturnRows(commissionRows)
	suite.mock.ExpectQuery(`SELECT (.+) FROM "gateways" WHERE "gateways"."commission_id" IN \(\$1,\$2\)`).
		WithArgs(1, 2).
		WillReturnRows(gatewayRows)

	resultCommissions, err := ListAllCommission(suite.DB)

	require := suite.Require()
	require.NoError(err, "Unexpected error")
	require.Equal(expectedCommissions, resultCommissions, "Commissions should match")

	require.NoError(suite.mock.ExpectationsWereMet())

}

func (suite *CommSuite) TestListAllCommission_Error() {
	suite.mock.ExpectQuery(`SELECT (.+) FROM "commissions"`).
		WillReturnError(errors.New("error getting commissions"))
	resultCommissions, err := ListAllCommission(suite.DB)
	require := suite.Require()
	require.Error(err, "Expect an error")
	require.Nil(resultCommissions, "Expected a nil commission slice")
	require.EqualError(err, "error getting commissions", "Error message should match")

	require.NoError(suite.mock.ExpectationsWereMet())
}

func TestCommissionSuite(t *testing.T) {
	suite.Run(t, new(CommSuite))
}
