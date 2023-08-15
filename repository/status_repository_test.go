package repository

import (
	"database/sql"
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var statusDummy = []model.Status{
	{
		StatusId: "1",
		Name:     "Status 1",
	},
	{
		StatusId: "2",
		Name:     "Status 2",
	},
	{
		StatusId: "3",
		Name:     "Status 3",
	},
}

type StatusRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    StatusRepository
}

func (suite *StatusRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewStatusRepository(db)
}

func (suite *StatusRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestStatusRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(StatusRepositoryTestSuite))
}

func (suite *StatusRepositoryTestSuite) TestCreateStatus_Success() {
	dummy := statusDummy[0]
	query := "INSERT INTO status (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.StatusId, dummy.Name).WillReturnResult(sqlmock.NewResult(1, 1))

	actualError := suite.repo.Create(dummy)

	assert.Nil(suite.T(), actualError)
}

func (suite *StatusRepositoryTestSuite) TestCreateStatus_Fai() {
	dummy := statusDummy[0]
	query := "INSERT INTO status (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.StatusId, dummy.Name).WillReturnError(fmt.Errorf("error"))

	actualError := suite.repo.Create(dummy)

	expectedError := fmt.Errorf("error")
	assert.NotNil(suite.T(), actualError)
	assert.EqualError(suite.T(), actualError, expectedError.Error())
}

func (suite *StatusRepositoryTestSuite) TestGetStatusByName_Success() {
	dummy := statusDummy[0]
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(dummy.StatusId, dummy.Name)

	query := "SELECT (.+) FROM status WHERE name ILIKE (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs("%" + dummy.Name + "%").WillReturnRows(rows)

	actualStatus, actualError := suite.repo.GetByName(dummy.Name)

	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), actualStatus, dummy)
}

func (suite *StatusRepositoryTestSuite) TestGetStatusByName_NotFound() {
	dummyStatusName := "Nonexistent Status"
	query := "SELECT (.+) FROM status WHERE name ILIKE (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs("%" + dummyStatusName + "%").
		WillReturnError(sql.ErrNoRows)

	actualStatus, actualError := suite.repo.GetByName(dummyStatusName)

	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), model.Status{}, actualStatus)
}

func (suite *StatusRepositoryTestSuite) TestListStatus_Succes() {
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, dummy := range statusDummy {
		rows.AddRow(dummy.StatusId, dummy.Name)
	}

	query := "SELECT (.+) FROM status"
	suite.mockSql.ExpectQuery(query).WillReturnRows(rows)
	statuses, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), statuses, 3)
	assert.Equal(suite.T(), statusDummy[0], statuses[0])
	assert.Equal(suite.T(), statusDummy[1], statuses[1])
	assert.Equal(suite.T(), statusDummy[2], statuses[2])

}

func (suite *StatusRepositoryTestSuite) TestListStatus_Fail() {
	query := "SELECT (.+) FROM status"
	suite.mockSql.ExpectQuery(query).WillReturnError(fmt.Errorf("error"))
	statuses, err := suite.repo.List()
	assert.Nil(suite.T(), statuses)
	assert.Error(suite.T(), err)
}

func (suite *StatusRepositoryTestSuite) TestGetStatusByID_Success() {
	dummy := statusDummy[0]
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(dummy.StatusId, dummy.Name)

	query := "SELECT (.+) FROM status WHERE id= (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs(dummy.StatusId).WillReturnRows(rows)

	actualStatus, actualError := suite.repo.Get(dummy.StatusId)

	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), actualStatus, dummy)
}

func (suite *StatusRepositoryTestSuite) TestGetStatusByID_NotFound() {
	dummyStatusID := "3rujwke"
	query := "SELECT (.+) FROM status WHERE id= (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs(dummyStatusID).
		WillReturnError(sql.ErrNoRows)

	actualStatus, actualError := suite.repo.Get(dummyStatusID)

	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), model.Status{}, actualStatus)
}

func (suite *StatusRepositoryTestSuite) TestUpdateStatus_Success() {
	dummy := statusDummy[0]
	query := "UPDATE status SET name = (.+) WHERE id = (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.StatusId, dummy.Name).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Update(dummy)

	assert.Nil(suite.T(), actualError)
	assert.NoError(suite.T(), actualError)
}

func (suite *StatusRepositoryTestSuite) TestUpdateStatus_Fail() {
	dummy := statusDummy[0]
	query := "UPDATE status SET name = (.+) WHERE id = (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.StatusId, dummy.Name).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Update(dummy)

	assert.Error(suite.T(), actualError)
}

func (suite *StatusRepositoryTestSuite) TestDeleteStatus_Success() {
	dummy := statusDummy[0]
	query := "DELETE FROM status WHERE id=(.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.StatusId).WillReturnResult(sqlmock.NewResult(0, 1))
	actualError := suite.repo.Delete(dummy.StatusId)
	assert.Nil(suite.T(), actualError)
}

func (suite *StatusRepositoryTestSuite) TestDeleteStatus_Fail() {
	dummy := statusDummy[0]
	query := "DELETE FROM status WHERE id=(.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.StatusId).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Delete(dummy.StatusId)
	assert.Error(suite.T(), actualError)
}

func (suite *StatusRepositoryTestSuite) TestPagingStatus_Succes() {
	requestPaging := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	paginationQuery := dto.PaginationQuery{
		Take: 5,
		Skip: 0,
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, dummy := range statusDummy {
		rows.AddRow(dummy.StatusId, dummy.Name)
	}
	selectQuery := "SELECT (.+) FROM status  LIMIT (.+) OFFSET (.+)"
	countQuery := "SELECT COUNT\\(.*\\) FROM status"
	suite.mockSql.ExpectQuery(selectQuery).WithArgs(paginationQuery.Take, paginationQuery.Skip).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(3)
	suite.mockSql.ExpectQuery(countQuery).WillReturnRows(rowCount)
	actualStatus, actualPaging, actualError := suite.repo.Paging(requestPaging)
	assert.Nil(suite.T(), actualError)
	assert.NotNil(suite.T(), actualStatus)
	assert.Equal(suite.T(), actualPaging.TotalRows, 3)
}

func (suite *StatusRepositoryTestSuite) TestPagingStatus_QueryPagingError() {
	selectQuery := "SELECT (.+) FROM status  LIMIT (.+) OFFSET (.+)"
	suite.mockSql.ExpectQuery(selectQuery).WillReturnError(fmt.Errorf("error"))
	actualStatus, actualPaging, actualError := suite.repo.Paging(dto.PaginationParam{})
	assert.Error(suite.T(), actualError)
	assert.Nil(suite.T(), actualStatus)
	assert.Equal(suite.T(), actualPaging.TotalRows, 0)
}

func (suite *StatusRepositoryTestSuite) TestPagingStatus_QueryCountError() {
	requestPaging := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	paginationQuery := dto.PaginationQuery{
		Take: 5,
		Skip: 0,
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, dummy := range statusDummy {
		rows.AddRow(dummy.StatusId, dummy.Name)
	}
	selectQuery := "SELECT (.+) FROM status  LIMIT (.+) OFFSET (.+)"
	countQuery := "SELECT COUNT\\(.*\\) FROM status"
	suite.mockSql.ExpectQuery(selectQuery).WithArgs(paginationQuery.Take, paginationQuery.Skip).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(3)
	suite.mockSql.ExpectQuery(countQuery).WillReturnError(fmt.Errorf("error"))
	_, actualPaging, actualError := suite.repo.Paging(requestPaging)
	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), actualPaging.TotalRows, 0)
}
