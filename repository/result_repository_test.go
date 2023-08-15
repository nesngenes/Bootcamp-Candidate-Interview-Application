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

var resultDummy = []model.Result{
	{
		ResultId: "1", 
		Name: "Result 1",     
	},
	{
		ResultId: "2", 
		Name: "Result 2",     
	},
	{
		ResultId: "3", 
		Name: "Result 3",     
	},
}


type ResultRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    ResultRepository
}

func (suite *ResultRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewResultRepository(db)
}

func (suite *ResultRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestResultRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ResultRepositoryTestSuite))
}

func (suite *ResultRepositoryTestSuite) TestCreateResult_Success() {
	dummy := resultDummy[0]
	query := "INSERT INTO result (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.ResultId, dummy.Name).WillReturnResult(sqlmock.NewResult(1, 1))

	actualError := suite.repo.Create(dummy)

	assert.Nil(suite.T(), actualError)
}

func (suite *ResultRepositoryTestSuite) TestCreateResult_Fai() {
	dummy := resultDummy[0]
	query := "INSERT INTO result (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.ResultId, dummy.Name).WillReturnError(fmt.Errorf("error"))

	actualError := suite.repo.Create(dummy)

	expectedError := fmt.Errorf("error")
	assert.NotNil(suite.T(), actualError)
	assert.EqualError(suite.T(), actualError, expectedError.Error())
}

func (suite *ResultRepositoryTestSuite) TestGetResultByName_Success() {
	dummy := resultDummy[0]
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(dummy.ResultId, dummy.Name)

	query := "SELECT (.+) FROM result WHERE name ILIKE (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs("%" + dummy.Name + "%").WillReturnRows(rows)

	actualResult, actualError := suite.repo.GetByName(dummy.Name)

	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), actualResult, dummy)
}

func (suite *ResultRepositoryTestSuite) TestGetResultByName_NotFound() {
	dummyResultName := "Nonexistent Result"
	query := "SELECT (.+) FROM result WHERE name ILIKE (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs("%" + dummyResultName + "%").
		WillReturnError(sql.ErrNoRows)

	actualResult, actualError := suite.repo.GetByName(dummyResultName)

	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), model.Result{}, actualResult)
}

func (suite *ResultRepositoryTestSuite) TestListResult_Succes() {
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, dummy := range resultDummy {
		rows.AddRow(dummy.ResultId, dummy.Name)
	}

	query := "SELECT (.+) FROM result"
	suite.mockSql.ExpectQuery(query).WillReturnRows(rows)
	results, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 3)
	assert.Equal(suite.T(), resultDummy[0], results[0])
	assert.Equal(suite.T(), resultDummy[1], results[1])
	assert.Equal(suite.T(), resultDummy[2], results[2])

}

func (suite *ResultRepositoryTestSuite) TestListResult_Fail() {
	query := "SELECT (.+) FROM result"
	suite.mockSql.ExpectQuery(query).WillReturnError(fmt.Errorf("error"))
	results, err := suite.repo.List()
	assert.Nil(suite.T(), results)
	assert.Error(suite.T(), err)
}

func (suite *ResultRepositoryTestSuite) TestGetResultByID_Success() {
	dummy := resultDummy[0]
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(dummy.ResultId, dummy.Name)

	query := "SELECT (.+) FROM result WHERE id= (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs(dummy.ResultId).WillReturnRows(rows)

	actualResult, actualError := suite.repo.Get(dummy.ResultId)

	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), actualResult, dummy)
}

func (suite *ResultRepositoryTestSuite) TestGetResultByID_NotFound() {
	dummyResultID := "3rujwke"
	query := "SELECT (.+) FROM result WHERE id= (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs(dummyResultID).
		WillReturnError(sql.ErrNoRows)

	actualResult, actualError := suite.repo.Get(dummyResultID)

	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), model.Result{}, actualResult)
}

func (suite *ResultRepositoryTestSuite) TestUpdateResult_Success() {
	dummy := resultDummy[0]
	query := "UPDATE result SET name = (.+) WHERE id = (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.ResultId, dummy.Name).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Update(dummy)

	assert.Nil(suite.T(), actualError)
	assert.NoError(suite.T(), actualError)
}

func (suite *ResultRepositoryTestSuite) TestUpdateResult_Fail() {
	dummy := resultDummy[0]
	query := "UPDATE result SET name = (.+) WHERE id = (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.ResultId, dummy.Name).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Update(dummy)

	assert.Error(suite.T(), actualError)
}

func (suite *ResultRepositoryTestSuite) TestDeleteResult_Success() {
	dummy := resultDummy[0]
	query := "DELETE FROM result WHERE id=(.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.ResultId).WillReturnResult(sqlmock.NewResult(0, 1))
	actualError := suite.repo.Delete(dummy.ResultId)
	assert.Nil(suite.T(), actualError)
}

func (suite *ResultRepositoryTestSuite) TestDeleteResult_Fail() {
	dummy := resultDummy[0]
	query := "DELETE FROM result WHERE id=(.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.ResultId).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Delete(dummy.ResultId)
	assert.Error(suite.T(), actualError)
}

func (suite *ResultRepositoryTestSuite) TestPagingResult_Succes() {
	requestPaging := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	paginationQuery := dto.PaginationQuery{
		Take: 5,
		Skip: 0,
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, dummy := range resultDummy {
		rows.AddRow(dummy.ResultId, dummy.Name)
	}
	selectQuery := "SELECT (.+) FROM result  LIMIT (.+) OFFSET (.+)"
	countQuery := "SELECT COUNT\\(.*\\) FROM result"
	suite.mockSql.ExpectQuery(selectQuery).WithArgs(paginationQuery.Take, paginationQuery.Skip).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(3)
	suite.mockSql.ExpectQuery(countQuery).WillReturnRows(rowCount)
	actualResult, actualPaging, actualError := suite.repo.Paging(requestPaging)
	assert.Nil(suite.T(), actualError)
	assert.NotNil(suite.T(), actualResult)
	assert.Equal(suite.T(), actualPaging.TotalRows, 3)
}

func (suite *ResultRepositoryTestSuite) TestPagingResult_QueryPagingError() {
	selectQuery := "SELECT (.+) FROM result  LIMIT (.+) OFFSET (.+)"
	suite.mockSql.ExpectQuery(selectQuery).WillReturnError(fmt.Errorf("error"))
	actualResult, actualPaging, actualError := suite.repo.Paging(dto.PaginationParam{})
	assert.Error(suite.T(), actualError)
	assert.Nil(suite.T(), actualResult)
	assert.Equal(suite.T(), actualPaging.TotalRows, 0)
}

func (suite *ResultRepositoryTestSuite) TestPagingResult_QueryCountError() {
	requestPaging := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	paginationQuery := dto.PaginationQuery{
		Take: 5,
		Skip: 0,
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, dummy := range resultDummy {
		rows.AddRow(dummy.ResultId, dummy.Name)
	}
	selectQuery := "SELECT (.+) FROM result  LIMIT (.+) OFFSET (.+)"
	countQuery := "SELECT COUNT\\(.*\\) FROM result"
	suite.mockSql.ExpectQuery(selectQuery).WithArgs(paginationQuery.Take, paginationQuery.Skip).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(3)
	suite.mockSql.ExpectQuery(countQuery).WillReturnError(fmt.Errorf("error"))
	_, actualPaging, actualError := suite.repo.Paging(requestPaging)
	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), actualPaging.TotalRows, 0)
}

