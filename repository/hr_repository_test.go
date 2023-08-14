package repository

import (
	"database/sql"
	"fmt"
	"interview_bootcamp/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HRRecruitmentRepositoryTestSuite struct {
	suite.Suite
	mockDb    *sql.DB
	mockSql   sqlmock.Sqlmock
	repo      HRRecruitmentRepository
	errMocked error
}

var hrDummy = []model.HRRecruitment{
	{
		ID:       "1",
		FullName: "Ela",
		UserID:   "1a",
		User: model.Users{
			Id:       "1a",
			Email:    "ella@gmail.com",
			UserName: "ella1",
			UserRole: model.UserRoles{
				Id:   "1",
				Name: "HR",
			},
		},
	},
	{
		ID:       "2",
		FullName: "Elle",
		UserID:   "1b",
		User: model.Users{
			Id:       "1b",
			Email:    "elle@gmail.com",
			UserName: "elle1",
			UserRole: model.UserRoles{
				Id:   "1",
				Name: "HR",
			},
		},
	},
}

func (suite *HRRecruitmentRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mockDb = db
	suite.mockSql = mock
	suite.errMocked = err
	suite.repo = NewHRRecruitmentRepository(db)
}

func (suite *HRRecruitmentRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestHRRepository(t *testing.T) {
	suite.Run(t, new(HRRecruitmentRepositoryTestSuite))
}

func (suite *HRRecruitmentRepositoryTestSuite) TestCreate_Success() {
	dummy := hrDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO hr_recruitment (.+)").WithArgs(dummy.ID, dummy.FullName, dummy.UserID).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Create(dummy)
	assert.Nil(suite.T(), actualError)     //hasilnya emang nil
	assert.NoError(suite.T(), actualError) //hasilnya tidak ada error
}

func (suite *HRRecruitmentRepositoryTestSuite) TestCreate_Fail() {
	dummy := hrDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO hr_recruitment (.+)").WithArgs(dummy.ID, dummy.FullName, dummy.UserID).WillReturnError(suite.errMocked)
	actualError := suite.repo.Create(dummy)
	assert.Error(suite.T(), actualError) //hasilnya ada error
}

func (suite *HRRecruitmentRepositoryTestSuite) TestList_Success() {
	rows := sqlmock.NewRows([]string{"id", "full_name", "user_id", "id", "email", "username", "id", "name"})
	for _, hr := range hrDummy {
		rows.AddRow(hr.ID, hr.FullName, hr.UserID, hr.User.Id, hr.User.Email, hr.User.UserName, hr.User.UserRole.Id, hr.User.UserRole.Name)
	}
	suite.mockSql.ExpectQuery("SELECT (.+) FROM hr_recruitment (.+)").WillReturnRows(rows)
	hrRecruitment, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), hrRecruitment, 2)
	assert.Equal(suite.T(), hrDummy[0], hrDummy[0])
	assert.Equal(suite.T(), hrDummy[1], hrDummy[1])
}

func (suite *HRRecruitmentRepositoryTestSuite) TestList_Fail() {
	suite.mockSql.ExpectQuery("SELECT (.+) FROM hr_recruitment (.+)").WillReturnError(fmt.Errorf("error"))
	hrRecruitments, err := suite.repo.List()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), hrRecruitments)
}

func (suite *HRRecruitmentRepositoryTestSuite) TestGetByID_Success() {
	expectedHRRecruitment := hrDummy[0]
	rows := sqlmock.NewRows([]string{"id", "full_name", "user_id", "id", "email", "username", "id", "name"})
	rows.AddRow(
		expectedHRRecruitment.ID, expectedHRRecruitment.FullName, expectedHRRecruitment.UserID,
		expectedHRRecruitment.User.Id, expectedHRRecruitment.User.Email, expectedHRRecruitment.User.UserName,
		expectedHRRecruitment.User.UserRole.Id, expectedHRRecruitment.User.UserRole.Name,
	)
	suite.mockSql.ExpectQuery("SELECT (.+) FROM hr_recruitment (.+) WHERE hr.id = ?").WithArgs(expectedHRRecruitment.ID).WillReturnRows(rows)
	actualHRRecruitment, actualError := suite.repo.Get(expectedHRRecruitment.ID)
	assert.NoError(suite.T(), actualError)
	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), expectedHRRecruitment, actualHRRecruitment)
}

func (suite *HRRecruitmentRepositoryTestSuite) TestGetByID_Fail() {
	suite.mockSql.ExpectQuery("SELECT (.+) FROM hr_recruitment (.+) WHERE hr.id = ?").WithArgs("1xxx").WillReturnError(fmt.Errorf("error"))
	actualHRRecruitment, err := suite.repo.Get("1xxx")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.HRRecruitment{}, actualHRRecruitment)
}

func (suite *HRRecruitmentRepositoryTestSuite) TestGetByUserID_Success() {
	expectedHRRecruitment := hrDummy[0]
	rows := sqlmock.NewRows([]string{"id", "full_name", "user_id", "id", "email", "username", "id", "name"})
	rows.AddRow(
		expectedHRRecruitment.ID, expectedHRRecruitment.FullName, expectedHRRecruitment.UserID,
		expectedHRRecruitment.User.Id, expectedHRRecruitment.User.Email, expectedHRRecruitment.User.UserName,
		expectedHRRecruitment.User.UserRole.Id, expectedHRRecruitment.User.UserRole.Name,
	)
	suite.mockSql.ExpectQuery("SELECT (.+) FROM hr_recruitment (.+) WHERE hr.user_id = ?").WithArgs(expectedHRRecruitment.UserID).WillReturnRows(rows)
	actualHRRecruitment, actualError := suite.repo.GetByUserID(expectedHRRecruitment.UserID)
	assert.NoError(suite.T(), actualError)
	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), expectedHRRecruitment, actualHRRecruitment)
}

func (suite *HRRecruitmentRepositoryTestSuite) TestGetByUserID_Fail() {
	suite.mockSql.ExpectQuery("SELECT (.+) FROM hr_recruitment (.+) WHERE hr.user_id = ?").WithArgs("user1xxx").WillReturnError(fmt.Errorf("error"))
	actualHRRecruitment, err := suite.repo.GetByUserID("user1xxx")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.HRRecruitment{}, actualHRRecruitment)
}

func (suite *HRRecruitmentRepositoryTestSuite) TestUpdate_Success() {
	dummy := hrDummy[0]
	suite.mockSql.ExpectExec(`UPDATE hr_recruitment SET full_name = \$2, user_id = \$3 WHERE id = \$1`).
		WithArgs(dummy.ID, dummy.FullName, dummy.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Update(dummy)
	assert.Nil(suite.T(), actualError)
}

func (suite *HRRecruitmentRepositoryTestSuite) TestUpdate_Fail() {
	dummy := hrDummy[0]
	suite.mockSql.ExpectExec("UPDATE hr_recruitment SET full_name = ?, user_id = ? WHERE id = ?").WithArgs(dummy.FullName, dummy.UserID, dummy.ID).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Update(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *HRRecruitmentRepositoryTestSuite) TestDelete_Success() {
	suite.mockSql.ExpectExec("DELETE FROM hr_recruitment WHERE id = ?").WithArgs("1ABC").WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Delete("1ABC")
	assert.Nil(suite.T(), actualError)
}

func (suite *HRRecruitmentRepositoryTestSuite) TestDelete_Fail() {
	suite.mockSql.ExpectExec("DELETE FROM hr_recruitment WHERE id = ?").WithArgs("1ABC").WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Delete("1ABC")
	assert.Error(suite.T(), actualError)
}
