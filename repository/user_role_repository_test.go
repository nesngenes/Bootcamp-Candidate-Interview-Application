package repository

import (
	"database/sql"
	"interview_bootcamp/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUserRolesRepository(t *testing.T) {
	suite.Run(t, new(UserRolesRepositoryTestSuite))
}

type UserRolesRepositoryTestSuite struct {
	suite.Suite
	mockDb    *sql.DB
	mockSql   sqlmock.Sqlmock
	repo      UserRolesRepository
	errMocked error
}

var userRolesDummy = []model.UserRoles{
	{
		Id:   "1",
		Name: "Role A",
	},
	{
		Id:   "2",
		Name: "Role B",
	},
}

func (suite *UserRolesRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.mockDb = db
	suite.mockSql = mock
	suite.errMocked = err
	suite.repo = NewUserRolesRepository(db)
}

func (suite *UserRolesRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func (suite *UserRolesRepositoryTestSuite) TestCreate_Success() {
	dummy := userRolesDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO user_roles (.+)").WithArgs(dummy.Id, dummy.Name).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Create(dummy)
	assert.NoError(suite.T(), actualError)
}

func (suite *UserRolesRepositoryTestSuite) TestCreate_Fail() {
	dummy := userRolesDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO user_roles (.+)").WithArgs(dummy.Id, dummy.Name).WillReturnError(suite.errMocked)
	actualError := suite.repo.Create(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *UserRolesRepositoryTestSuite) TestList_Success() {
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("1", "Admin").
		AddRow("2", "User")

	suite.mockSql.ExpectQuery("SELECT id, name FROM user_roles").WillReturnRows(rows)
	userRoles, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), userRoles, 2)
}

func (suite *UserRolesRepositoryTestSuite) TestList_Fail() {
	suite.mockSql.ExpectQuery("SELECT id, name FROM user_roles").WillReturnError(suite.errMocked)
	userRoles, err := suite.repo.List()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), userRoles)
}

func (suite *UserRolesRepositoryTestSuite) TestGet_Success() {
	expectedUserRole := model.UserRoles{
		Id:   "1",
		Name: "Admin",
	}

	suite.mockSql.ExpectQuery("SELECT id, name FROM user_roles WHERE id = ?").WithArgs(expectedUserRole.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUserRole.Id, expectedUserRole.Name))
	actualUserRole, err := suite.repo.Get(expectedUserRole.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUserRole, actualUserRole)
}

func (suite *UserRolesRepositoryTestSuite) TestGet_Fail() {
	suite.mockSql.ExpectQuery("SELECT id, name FROM user_roles WHERE id = ?").WithArgs("999").
		WillReturnError(suite.errMocked)
	actualUserRole, err := suite.repo.Get("999")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.UserRoles{}, actualUserRole)
}

func (suite *UserRolesRepositoryTestSuite) TestGetByName_Success() {
	expectedUserRole := userRolesDummy[0]

	suite.mockSql.ExpectQuery("SELECT id, name FROM user_roles WHERE name ILIKE ?").WithArgs("%Admin%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUserRole.Id, expectedUserRole.Name))
	actualUserRole, err := suite.repo.GetByName("Admin")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUserRole, actualUserRole)
}

func (suite *UserRolesRepositoryTestSuite) TestGetByName_Fail() {
	suite.mockSql.ExpectQuery("SELECT id, name FROM user_roles WHERE name ILIKE ?").WithArgs("%Admin%").
		WillReturnError(suite.errMocked)
	actualUserRole, err := suite.repo.GetByName("Admin")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.UserRoles{}, actualUserRole)
}

func (suite *UserRolesRepositoryTestSuite) TestUpdate_Success() {
	dummy := userRolesDummy[0]
	suite.mockSql.ExpectExec(`UPDATE user_roles SET name = \$2 WHERE id = \$1`).
		WithArgs(dummy.Id, dummy.Name).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Update(dummy)
	assert.NoError(suite.T(), actualError)
}

func (suite *UserRolesRepositoryTestSuite) TestUpdate_Fail() {
	dummy := userRolesDummy[0]

	suite.mockSql.ExpectExec("UPDATE user_roles SET name = ? WHERE id = ?").
		WithArgs(dummy.Name, dummy.Id).WillReturnError(suite.errMocked)
	actualError := suite.repo.Update(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *UserRolesRepositoryTestSuite) TestDelete_Success() {
	suite.mockSql.ExpectExec("DELETE FROM user_roles WHERE id = ?").
		WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Delete("1")
	assert.NoError(suite.T(), actualError)
}

func (suite *UserRolesRepositoryTestSuite) TestDelete_Fail() {
	suite.mockSql.ExpectExec("DELETE FROM user_roles WHERE id = ?").
		WithArgs("1").WillReturnError(suite.errMocked)
	actualError := suite.repo.Delete("1")
	assert.Error(suite.T(), actualError)
}
