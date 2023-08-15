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


var formDummy = []model.Form {
	{
		FormID: "1",
		FormLink: "test_link_1.com",
	},
	{
		FormID: "2",
		FormLink: "test_link_2.com",
	},
	{
		FormID: "1",
		FormLink: "test_link_1.com",
	},
}

type FormRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    FormRepository
}

func (suite *FormRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewFormRepository(db)
}

func (suite *FormRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestFormRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(FormRepositoryTestSuite))
}

func (suite *FormRepositoryTestSuite) TestCreateNewForm_Success() {
	dummy := formDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO form (.+)").WithArgs(dummy.FormID, dummy.FormLink).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Create(dummy)
	assert.Nil(suite.T(), actualError)
	assert.NoError(suite.T(), actualError)
}

func (suite *FormRepositoryTestSuite) TestCreateNewForm_Fail() {
	dummy := formDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO form (.+)").WithArgs(dummy.FormID, dummy.FormLink).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Create(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *FormRepositoryTestSuite) TestUpdateForm_Success() {
	dummy := formDummy[0]
	query := "UPDATE form SET form_link = (.+) WHERE id = (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.FormID, dummy.FormLink).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Update(dummy)

	assert.Nil(suite.T(), actualError)
	assert.NoError(suite.T(), actualError)
}

func (suite *FormRepositoryTestSuite) TestUpdateForm_Fail() {
	dummy := formDummy[0]
	query := "UPDATE form SET form_link = (.+) WHERE id = (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.FormID, dummy.FormLink).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Update(dummy)

	assert.Error(suite.T(), actualError)
}


func (suite *FormRepositoryTestSuite) TestDeleteForm_Success() {
	dummy := formDummy[0]
	query := "DELETE FROM form WHERE id=(.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.FormID).WillReturnResult(sqlmock.NewResult(0, 1))
	actualError := suite.repo.Delete(dummy.FormID)
	assert.Nil(suite.T(), actualError)
}

func (suite *FormRepositoryTestSuite) TestDeleteForm_Fail() {
	dummy := formDummy[0]
	query := "DELETE FROM form WHERE id=(.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.FormID).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Delete(dummy.FormID)
	assert.Error(suite.T(), actualError)
}

func (suite *FormRepositoryTestSuite) TestPagingForm_QueryPagingError() {
	selectQuery := "SELECT (.+) FROM form  LIMIT (.+) OFFSET (.+)"
	suite.mockSql.ExpectQuery(selectQuery).WillReturnError(fmt.Errorf("error"))
	actualForm, actualPaging, actualError := suite.repo.Paging(dto.PaginationParam{})
	assert.Error(suite.T(), actualError)
	assert.Nil(suite.T(), actualForm)
	assert.Equal(suite.T(), actualPaging.TotalRows, 0)
}

func (suite *FormRepositoryTestSuite) TestPagingForm_QueryCountError() {
	requestPaging := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	paginationQuery := dto.PaginationQuery{
		Take: 5,
		Skip: 0,
	}
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, dummy := range formDummy {
		rows.AddRow(dummy.FormID, dummy.FormLink)
	}
	selectQuery := "SELECT (.+) FROM form  LIMIT (.+) OFFSET (.+)"
	countQuery := "SELECT COUNT\\(.*\\) FROM form"
	suite.mockSql.ExpectQuery(selectQuery).WithArgs(paginationQuery.Take, paginationQuery.Skip).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(3)
	suite.mockSql.ExpectQuery(countQuery).WillReturnError(fmt.Errorf("error"))
	_, actualPaging, actualError := suite.repo.Paging(requestPaging)
	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), actualPaging.TotalRows, 0)
}
