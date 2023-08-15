package repository

import (
	"database/sql"
	"fmt"
	"interview_bootcamp/model"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var dummyInterviewers = []model.Interviewer{
	{
		InterviewerID: "1",
		FullName:      "First",
		UserID:        "abcdefg",
	},
	{
		InterviewerID: "2",
		FullName:      "Second",
		UserID:        "abcdefh",
	},
	{
		InterviewerID: "3",
		FullName:      "Three",
		UserID:        "abcdefi",
	},
}

type InterviewerRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    InterviewerRepository
}

func (suite *InterviewerRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewInterviewerRepository(db)
}

func (suite *InterviewerRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestInterviewerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(InterviewerRepositoryTestSuite))
}

func (suite *InterviewerRepositoryTestSuite) TestCreateInterviewer_Success() {
	dummy := dummyInterviewers[0]
	suite.mockSql.ExpectExec("INSERT INTO interviewer (.+)").WithArgs(dummy.InterviewerID, dummy.FullName, dummy.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

	actualError := suite.repo.Create(dummy)

	assert.Nil(suite.T(), actualError)
}

func (suite *InterviewerRepositoryTestSuite) TestCreateInterviewer_Failure() {
	dummyInterviewer := dummyInterviewers[0]
	suite.mockSql.ExpectExec("INSERT INTO interviewer (.+)").
		WithArgs(dummyInterviewer.UserID, dummyInterviewer.FullName, dummyInterviewer.UserID).
		WillReturnError(fmt.Errorf("error"))

	err := suite.repo.Create(dummyInterviewer)

	assert.NotNil(suite.T(), err)
}

func (suite *InterviewerRepositoryTestSuite) TestListInterviewers_Success() {
	rows := sqlmock.NewRows([]string{"id", "full_name", "id"})
	for _, dummy := range dummyInterviewers {
		rows.AddRow(dummy.InterviewerID, dummy.FullName, dummy.UserID)
	}

	query := "SELECT i.id, i.full_name, u.id FROM interviewer i INNER JOIN users u ON u.id = i.user_id"
	suite.mockSql.ExpectQuery(query).
		WillReturnRows(rows)

	interviewers, err := suite.repo.List()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyInterviewers, interviewers)
}

func (suite *InterviewerRepositoryTestSuite) TestListInterviewers_Fail() {
	query := "SELECT i.id, i.full_name, u.id FROM interviewer i INNER JOIN users u ON u.id = i.user_id"
	suite.mockSql.ExpectQuery(query).
		WillReturnError(fmt.Errorf("error"))
	interviewers, err := suite.repo.List()

	assert.Nil(suite.T(), interviewers)
	assert.Error(suite.T(), err)
}

func (suite *InterviewerRepositoryTestSuite) TestGetInterviewer_Success() {
	dummy := dummyInterviewers[0]
	rows := sqlmock.NewRows([]string{"id", "full_name", "id"})

	query := "SELECT i.id, i.full_name, u.id FROM interviewer i INNER JOIN users u ON u.id = i.user_id WHERE i.id = $1"
	rows.AddRow(dummy.InterviewerID, dummy.FullName, dummy.UserID)
	suite.mockSql.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(dummy.InterviewerID).
		WillReturnRows(rows)

	interviewer, err := suite.repo.Get(dummy.InterviewerID)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummy, interviewer)
}

func (suite *InterviewerRepositoryTestSuite) TestGetInterviewer_NotFound() {
	query := "SELECT i.id, i.full_name, u.id FROM interviewer i INNER JOIN users u ON u.id = i.user_id WHERE i.id = $1"
	suite.mockSql.ExpectQuery(query).WithArgs("2").
		WillReturnError(fmt.Errorf("error"))

	interviewer, err := suite.repo.Get("2")

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), model.Interviewer{}, interviewer)
}

func (suite *InterviewerRepositoryTestSuite) TestUpdateInterviewer_Success() {
	dummy := dummyInterviewers[0]
	query := "UPDATE interviewer (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.InterviewerID, dummy.FullName, dummy.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := suite.repo.Update(dummy)

	assert.Nil(suite.T(), err)
}

func (suite *InterviewerRepositoryTestSuite) TestUpdateInterviewer_Failure() {
	dummy := dummyInterviewers[0]
	query := "UPDATE interviewer (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.InterviewerID, dummy.FullName, dummy.UserID).
		WillReturnError(fmt.Errorf("error"))

	err := suite.repo.Update(dummy)

	assert.NotNil(suite.T(), err)
}

func (suite *InterviewerRepositoryTestSuite) TestDeleteInterviewer_Success() {
	dummy := dummyInterviewers[0]
	query := "DELETE FROM interviewer (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.InterviewerID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := suite.repo.Delete(dummy.InterviewerID)

	assert.Nil(suite.T(), err)
}

func (suite *InterviewerRepositoryTestSuite) TestDeleteInterviewer_Failure() {
	dummy := dummyInterviewers[0]
	query := "DELETE FROM interviewer (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.InterviewerID).
		WillReturnError(fmt.Errorf("error"))

	err := suite.repo.Delete(dummy.InterviewerID)

	assert.NotNil(suite.T(), err)
}
