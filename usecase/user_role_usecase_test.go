package usecase

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"interview_bootcamp/model"
)

type userRepoMock struct {
	mock.Mock
}

func (r *userRepoMock) Create(userRole model.UserRoles) error {
	args := r.Called(userRole)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *userRepoMock) List() ([]model.UserRoles, error) {
	args := r.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.UserRoles), nil
}

func (r *userRepoMock) Get(id string) (model.UserRoles, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.UserRoles{}, args.Error(1)
	}
	return args.Get(0).(model.UserRoles), nil
}

func (r *userRepoMock) GetByName(name string) (model.UserRoles, error) {
	args := r.Called(name)
	if args.Get(1) != nil {
		return model.UserRoles{}, args.Error(1)
	}
	return args.Get(0).(model.UserRoles), nil
}

func (r *userRepoMock) Update(userRole model.UserRoles) error {
	args := r.Called(userRole)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *userRepoMock) Delete(id string) error {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

type UserRoleUseCaseTestSuite struct {
	suite.Suite
	repoMock *userRepoMock
	usecase  UserRolesUseCase
}

func (suite *UserRoleUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(userRepoMock)
	suite.usecase = NewUserRolesUseCase(suite.repoMock)
}

func TestUserRoleUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserRoleUseCaseTestSuite))
}

func (suite *UserRoleUseCaseTestSuite) TestRegisterNewUserRole_Success() {
	dummyRole := model.UserRoles{
		Id:   "1",
		Name: "Admin",
	}

	suite.repoMock.On("GetByName", dummyRole.Name).Return(model.UserRoles{}, nil)
	suite.repoMock.On("Create", dummyRole).Return(nil)

	err := suite.usecase.RegisterNewUserRole(dummyRole)
	assert.NoError(suite.T(), err)
}

func (suite *UserRoleUseCaseTestSuite) TestRegisterNewUserRole_ExistingRole() {
	dummyRole := model.UserRoles{
		Id:   "1",
		Name: "Admin",
	}

	existingRole := model.UserRoles{
		Id:   "2",
		Name: "Admin",
	}

	suite.repoMock.On("GetByName", dummyRole.Name).Return(existingRole, nil)

	err := suite.usecase.RegisterNewUserRole(dummyRole)
	assert.Error(suite.T(), err)
}

func (suite *UserRoleUseCaseTestSuite) TestRegisterNewUserRole_EmptyName() {
	dummyRole := model.UserRoles{
		Id:   "1",
		Name: "",
	}

	err := suite.usecase.RegisterNewUserRole(dummyRole)
	assert.Error(suite.T(), err)
}

// tes buat list daftar nya
func (suite *UserRoleUseCaseTestSuite) TestGetAllUserRoles_Success() {
	dummyRoles := []model.UserRoles{
		{Id: "1", Name: "Admin"},
		{Id: "2", Name: "User"},
	}

	suite.repoMock.On("List").Return(dummyRoles, nil)

	actualRoles, err := suite.usecase.GetAllUserRoles()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummyRoles, actualRoles)
}

func (suite *UserRoleUseCaseTestSuite) TestGetAllUserRoles_Fail() {
	suite.repoMock.On("List").Return(nil, fmt.Errorf("error"))

	actualRoles, err := suite.usecase.GetAllUserRoles()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), actualRoles)
}

func (suite *UserRoleUseCaseTestSuite) TestGetUserRoleByID_Success() {
	dummyRole := model.UserRoles{
		Id:   "1",
		Name: "Admin",
	}

	suite.repoMock.On("Get", dummyRole.Id).Return(dummyRole, nil)

	actualRole, err := suite.usecase.GetUserRoleByID(dummyRole.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummyRole, actualRole)
}

func (suite *UserRoleUseCaseTestSuite) TestGetUserRoleByID_Fail() {
	suite.repoMock.On("Get", "1").Return(model.UserRoles{}, fmt.Errorf("error"))

	actualRole, err := suite.usecase.GetUserRoleByID("1")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.UserRoles{}, actualRole)
}

func (suite *UserRoleUseCaseTestSuite) TestUpdateUserRole_Success() {
	dummyRole := model.UserRoles{
		Id:   "1",
		Name: "Admin",
	}

	suite.repoMock.On("GetByName", dummyRole.Name).Return(model.UserRoles{}, nil)
	suite.repoMock.On("Update", dummyRole).Return(nil)

	err := suite.usecase.UpdateUserRole(dummyRole)
	assert.NoError(suite.T(), err)
}

// test buat nama role yang udah ada.
func (suite *UserRoleUseCaseTestSuite) TestUpdateUserRole_ExistingRole() {
	dummyRole := model.UserRoles{
		Id:   "1",
		Name: "Admin",
	}

	existingRole := model.UserRoles{
		Id:   "2",
		Name: "Admin",
	}

	suite.repoMock.On("GetByName", dummyRole.Name).Return(existingRole, nil)

	err := suite.usecase.UpdateUserRole(dummyRole)
	assert.Error(suite.T(), err)
}

func (suite *UserRoleUseCaseTestSuite) TestUpdateUserRole_EmptyName() {
	dummyRole := model.UserRoles{
		Id:   "1",
		Name: "",
	}

	err := suite.usecase.UpdateUserRole(dummyRole)
	assert.Error(suite.T(), err)
}

func (suite *UserRoleUseCaseTestSuite) TestDeleteUserRole_Success() {
	dummyRole := model.UserRoles{
		Id:   "1",
		Name: "Admin",
	}

	suite.repoMock.On("Get", dummyRole.Id).Return(dummyRole, nil)
	suite.repoMock.On("Delete", dummyRole.Id).Return(nil)

	err := suite.usecase.DeleteUserRole(dummyRole.Id)
	assert.NoError(suite.T(), err)
}

func (suite *UserRoleUseCaseTestSuite) TestDeleteUserRole_NotFound() {
	suite.repoMock.On("Get", "1").Return(model.UserRoles{}, fmt.Errorf("error"))

	err := suite.usecase.DeleteUserRole("1")
	assert.Error(suite.T(), err)
}

func (suite *UserRoleUseCaseTestSuite) TestDeleteUserRole_Fail() {
	dummyRole := model.UserRoles{
		Id:   "1",
		Name: "Admin",
	}

	suite.repoMock.On("Get", dummyRole.Id).Return(dummyRole, nil)
	suite.repoMock.On("Delete", dummyRole.Id).Return(fmt.Errorf("error"))

	err := suite.usecase.DeleteUserRole(dummyRole.Id)
	assert.Error(suite.T(), err)
}

func TestUserRolesUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserRoleUseCaseTestSuite))
}
