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

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

type UserRepositoryTestSuite struct {
	suite.Suite
	mockDb    *sql.DB
	mockSql   sqlmock.Sqlmock
	repo      UserRepository
	errMocked error
}

var userDummy = []model.Users{
	{
		Id:       "1a",
		Email:    "ella@gmail.com",
		UserName: "ella1",
		UserRole: model.UserRoles{
			Id:   "1",
			Name: "HR",
		},
	},

	{
		Id:       "1b",
		Email:    "elle@gmail.com",
		UserName: "elle1",
		UserRole: model.UserRoles{
			Id:   "1",
			Name: "HR",
		},
	},
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.mockDb = db
	suite.mockSql = mock
	suite.errMocked = err
	suite.repo = NewUserRepository(db)

}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func (suite *UserRepositoryTestSuite) TestCreate_Success() {
	dummy := userDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO users (.+)").WithArgs(dummy.Id, dummy.Email, dummy.UserName, dummy.Password, dummy.UserRole.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Create(dummy)
	assert.NoError(suite.T(), actualError)
}

func (suite *UserRepositoryTestSuite) TestCreate_Fail() {
	dummy := userDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO users (.+)").WithArgs(dummy.Id, dummy.Email, dummy.UserName, dummy.Password, dummy.UserRole.Id).WillReturnError(suite.errMocked)
	actualError := suite.repo.Create(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *UserRepositoryTestSuite) TestList_Success() {
	rows := sqlmock.NewRows([]string{"id", "email", "username", "id", "name"})
	for _, user := range userDummy {
		rows.AddRow(user.Id, user.Email, user.UserName, user.UserRole.Id, user.UserRole.Name)
	}

	suite.mockSql.ExpectQuery("SELECT (.+) FROM users (.+)").WillReturnRows(rows)
	users, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 2)
	assert.Equal(suite.T(), userDummy[0], userDummy[0])
	assert.Equal(suite.T(), userDummy[1], userDummy[1])
}

func (suite *UserRepositoryTestSuite) TestList_Fail() {
	suite.mockSql.ExpectQuery("SELECT (.+) FROM users (.+)").WillReturnError(suite.errMocked)
	users, err := suite.repo.List()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), users)
}

func (suite *UserRepositoryTestSuite) TestGet_Success() {
	expectedUser := userDummy[0]
	rows := sqlmock.NewRows([]string{"id", "email", "username", "id", "name"})
	rows.AddRow(
		expectedUser.Id, expectedUser.Email, expectedUser.UserName, expectedUser.UserRole.Id, expectedUser.UserRole.Name,
	)
	suite.mockSql.ExpectQuery("SELECT (.+) FROM users (.+) WHERE u.id = ?").WithArgs(expectedUser.Id).WillReturnRows(rows)
	actualUser, actualError := suite.repo.Get(expectedUser.Id)
	assert.NoError(suite.T(), actualError)
	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), expectedUser, actualUser)
}

func (suite *UserRepositoryTestSuite) TestGet_Fail() {
	suite.mockSql.ExpectQuery("SELECT (.+) FROM users (.+) WHERE u.id = ?").WithArgs("1xxx").WillReturnError(fmt.Errorf("error"))
	actualUser, err := suite.repo.Get("1xxx")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Users{}, actualUser)
}

func (suite *UserRepositoryTestSuite) TestUpdate_Success() {
	dummy := userDummy[0]
	suite.mockSql.ExpectExec(`UPDATE users SET email = \$2, username = \$3, role_id = \$4? WHERE id = \$1?`).
		WithArgs(dummy.Id, dummy.Email, dummy.UserName, dummy.UserRole.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Update(dummy)
	assert.NoError(suite.T(), actualError)
}

func (suite *UserRepositoryTestSuite) TestUpdate_Fail() {
	dummy := userDummy[0]
	suite.mockSql.ExpectExec("UPDATE users SET email = ?, username = ?, role_id = ? WHERE id = ?").
		WithArgs(dummy.Id, dummy.Email, dummy.UserName, dummy.UserRole.Id).WillReturnError(suite.errMocked)
	actualError := suite.repo.Update(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *UserRepositoryTestSuite) TestDelete_Success() {
	suite.mockSql.ExpectExec("DELETE FROM users WHERE id=?").WithArgs("1ABC").WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Delete("1ABC")
	assert.Nil(suite.T(), actualError)
}

func (suite *UserRepositoryTestSuite) TestDelete_Fail() {
	suite.mockSql.ExpectExec("DELETE FROM users WHERE id =?").WithArgs("1ABC").WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Delete("1ABC")
	assert.Error(suite.T(), actualError)
}

func (suite *UserRepositoryTestSuite) TestGetByUserName_Success() {
	expectedUser := userDummy[0]

	rows := sqlmock.NewRows([]string{"id", "email", "username", "password", "id", "name"})
	rows.AddRow(
		expectedUser.Id, expectedUser.Email, expectedUser.UserName, expectedUser.Password,
		expectedUser.UserRole.Id, expectedUser.UserRole.Name,
	)

	suite.mockSql.ExpectQuery("SELECT (.+) FROM users (.+) WHERE u.username ILIKE ?").
		WithArgs("%" + expectedUser.UserName + "%").
		WillReturnRows(rows)
	actualUser, err := suite.repo.GetByUserName(expectedUser.UserName)
	assert.NoError(suite.T(), err)

	// Ignore password field comparison
	expectedUser.Password = actualUser.Password
	assert.Equal(suite.T(), expectedUser, actualUser)
}

func (suite *UserRepositoryTestSuite) TestGetByUserName_Fail() {
	expectedUser := userDummy[0]
	suite.mockSql.ExpectQuery("SELECT (.+) FROM users (.+) WHERE u.username ILIKE ?").WithArgs("%" + expectedUser.UserName + "%").
		WillReturnError(suite.errMocked)

	actualUser, err := suite.repo.GetByUserName(expectedUser.UserName)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Users{}, actualUser)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail_Success() {
	expectedUser := userDummy[0]

	rows := sqlmock.NewRows([]string{"id", "email", "username", "id", "name"})
	rows.AddRow(
		expectedUser.Id, expectedUser.Email, expectedUser.UserName,
		expectedUser.UserRole.Id, expectedUser.UserRole.Name,
	)

	suite.mockSql.ExpectQuery("SELECT (.+) FROM users (.+) WHERE u.email ILIKE ?").
		WithArgs("%" + expectedUser.Email + "%").
		WillReturnRows(rows)
	actualUser, err := suite.repo.GetByEmail(expectedUser.Email)
	assert.NoError(suite.T(), err)

	// Ignore password field comparison
	expectedUser.Email = actualUser.Email
	assert.Equal(suite.T(), expectedUser, actualUser)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail_Fail() {
	expectedUser := userDummy[0]
	suite.mockSql.ExpectQuery("SELECT (.+) FROM users (.+) WHERE u.email ILIKE ?").WithArgs("%" + expectedUser.Email + "%").
		WillReturnError(suite.errMocked)

	actualUser, err := suite.repo.GetByEmail(expectedUser.Email)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Users{}, actualUser)
}
