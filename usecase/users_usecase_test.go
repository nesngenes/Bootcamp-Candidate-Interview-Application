package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type userRepoMock struct {
	mock.Mock
}

// GetUsernamePassword implements repository.UserRepository.
func (u *userRepoMock) GetUsernamePassword(username string, password string) (model.Users, error) {
	args := u.Called(username, password)
	if args.Get(1) != nil {
		return model.Users{}, args.Error(1)
	}
	return args.Get(0).(model.Users), nil
}

func (u *userRepoMock) Create(user model.Users) error {
	args := u.Called(user)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (u *userRepoMock) List() ([]model.Users, error) {
	args := u.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Users), nil
}

func (u *userRepoMock) Get(id string) (model.Users, error) {
	args := u.Called(id)
	if args.Get(1) != nil {
		return model.Users{}, args.Error(1)
	}
	return args.Get(0).(model.Users), nil
}

func (u *userRepoMock) GetByEmail(email string) (model.Users, error) {
	args := u.Called(email)
	if args.Get(1) != nil {
		return model.Users{}, args.Error(1)
	}
	return args.Get(0).(model.Users), nil
}

func (u *userRepoMock) GetByUserName(username string) (model.Users, error) {
	args := u.Called(username)
	if args.Get(1) != nil {
		return model.Users{}, args.Error(1)
	}
	return args.Get(0).(model.Users), nil
}

func (u *userRepoMock) Update(payload model.Users) error {
	args := u.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (u *userRepoMock) Delete(id string) error {
	args := u.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

type userUsecaseMock struct {
	mock.Mock
}

type UserUsecaseTestSuite struct {
	suite.Suite
	repoMock    *userRepoMock
	usecaseMock *userUsecaseMock
	usecase     UserUsecase
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.repoMock = new(userRepoMock)
	suite.usecaseMock = new(userUsecaseMock)
	suite.usecase = NewUserUsecase(suite.repoMock)
}

var userDummy = []model.Users{
	{
		Id:       "1",
		Email:    "ella@mail.com",
		UserName: "ella",
		UserRole: model.UserRoles{
			Id:   "1",
			Name: "HR",
		},
	},
	{
		Id:       "2",
		Email:    "robin@mail.com",
		UserName: "robin",
		UserRole: model.UserRoles{
			Id:   "2",
			Name: "Arkeolog",
		},
	},
}

func (suite *UserUsecaseTestSuite) TestRegisterNewUser_Success() {
	dummy := model.Users{
		Id:       "1",
		Email:    "ella@mail.com",
		UserName: "ella",
		Password: "password",
		UserRole: model.UserRoles{
			Id:   "1",
			Name: "HR",
		},
	}
	suite.repoMock.On("GetByUserName", dummy.UserName).Return(model.Users{}, nil)
	suite.repoMock.On("GetByEmail", dummy.Email).Return(model.Users{}, nil)

	// skip password field bc of hashing
	suite.repoMock.On("Create", mock.AnythingOfType("model.Users")).Return(nil)

	err := suite.usecase.RegisterNewUser(dummy)
	assert.Nil(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) TestRegisterNewUser_EmptyField() {
	suite.repoMock.On("Create", model.Users{}).Return(fmt.Errorf("field required"))
	err := suite.usecase.RegisterNewUser(model.Users{})
	assert.Error(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) TestRegisterNewUser_Fail() {
	suite.repoMock.On("Create", userDummy[0]).Return(fmt.Errorf("failed register"))
	err := suite.usecase.RegisterNewUser(userDummy[0])
	assert.Error(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) TestRegisterNewUser_UsernameExist() {
	suite.usecaseMock.On("GetByUserName", "ella").Return(userDummy, nil)
	suite.repoMock.On("Create", userDummy[0]).Return(fmt.Errorf("failed register"))
	err := suite.usecase.RegisterNewUser(userDummy[0])
	assert.Error(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) TestRegisterNewUser_EmailExist() {
	suite.usecaseMock.On("GetByEmail", "ella@mail.com").Return(userDummy, nil)
	suite.repoMock.On("Create", userDummy[0]).Return(fmt.Errorf("failed register"))
	err := suite.usecase.RegisterNewUser(userDummy[0])
	assert.Error(suite.T(), err)
}

func (suite *UserUsecaseTestSuite) TestList_Success() {
	dummy := userDummy
	suite.repoMock.On("List").Return(dummy, nil)
	actualList, err := suite.usecase.List()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummy, actualList)
}

func (suite *UserUsecaseTestSuite) TestList_Fail() {
	suite.repoMock.On("List").Return(nil, fmt.Errorf("error"))
	actualList, err := suite.usecase.List()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), actualList)
}

func (suite *UserUsecaseTestSuite) TestGetUserByID() {
	dummy := userDummy[0]
	suite.repoMock.On("Get", dummy.Id).Return(dummy, nil)

	actualRole, err := suite.usecase.GetUserByID(dummy.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummy, actualRole)
}
func (suite *UserUsecaseTestSuite) TestGetUserByEmail() {
	dummy := userDummy[0]
	suite.repoMock.On("GetByEmail", dummy.Email).Return(dummy, nil)

	actualRole, err := suite.usecase.GetUserByEmail(dummy.Email)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummy, actualRole)
}
func (suite *UserUsecaseTestSuite) TestGetUserByUserName() {
	dummy := userDummy[0]
	suite.repoMock.On("GetByUserName", dummy.UserName).Return(dummy, nil)

	actualRole, err := suite.usecase.GetUserByUserName(dummy.UserName)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummy, actualRole)
}

func (suite *UserUsecaseTestSuite) TestUpdateUser_Success() {
	dummy := userDummy[0]
	suite.repoMock.On("GetByUserName", dummy.UserName).Return(model.Users{}, nil)
	suite.repoMock.On("GetByEmail", dummy.Email).Return(model.Users{}, nil)
	suite.repoMock.On("Update", mock.AnythingOfType("model.Users")).Return(nil)
}

func (suite *UserUsecaseTestSuite) TestUpdateUser_EmailExist() {
	dummy := userDummy[0]

	existingUser := model.Users{
		Id:       "2",
		Email:    "ada@mail.com",
		UserName: "ella",
		UserRole: model.UserRoles{
			Id:   "1",
			Name: "HR",
		},
	}

	// Set up expectations for repository mock
	suite.repoMock.On("Get", dummy.Id).Return(existingUser, nil)
	suite.repoMock.On("GetByUserName", dummy.UserName).Return(existingUser, nil)
	suite.repoMock.On("GetByEmail", dummy.Email).Return(existingUser, nil)
	actualError := suite.usecase.UpdateUser(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *UserUsecaseTestSuite) TestUpdateUser_Fail() {
	suite.repoMock.On("Get", mock.AnythingOfType("string")).Return(model.Users{}, nil) //mock.Anything --> supaya valuenya lebih simple & ga terlalu ribet
	suite.repoMock.On("GetByUserName", mock.AnythingOfType("string")).Return(model.Users{}, nil)
	suite.repoMock.On("GetByEmail", mock.AnythingOfType("string")).Return(model.Users{}, nil)
	suite.repoMock.On("Update", mock.AnythingOfType("model.Users")).Return(fmt.Errorf("error"))

	actualError := suite.usecase.UpdateUser(model.Users{})
	assert.Error(suite.T(), actualError)
}

func (suite *UserUsecaseTestSuite) TestDeleteUser() {
	suite.repoMock.On("Delete", userDummy[0].Id).Return(nil)
	actualError := suite.usecase.DeleteUser(userDummy[0].Id)
	assert.Nil(suite.T(), actualError)

}

func (suite *UserUsecaseTestSuite) TestUpdateUser_UsernameExist() {
	dummy := model.Users{
		Id:       "1",
		Email:    "ella@mail.com",
		UserName: "newusername",
		UserRole: model.UserRoles{
			Id:   "1",
			Name: "HR",
		},
	}

	existingUser := model.Users{
		Id:       "2",
		Email:    "ellania@mail.com",
		UserName: "ella",
		UserRole: model.UserRoles{
			Id:   "1",
			Name: "HR",
		},
	}

	suite.repoMock.On("Get", dummy.Id).Return(existingUser, nil)
	suite.repoMock.On("GetByEmail", dummy.Email).Return(existingUser, nil)
	suite.repoMock.On("GetByUserName", dummy.UserName).Return(existingUser, nil)

	actualError := suite.usecase.UpdateUser(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *UserUsecaseTestSuite) TestDeleteUser_Fail() {
	suite.repoMock.On("Delete", "1xxx").Return(fmt.Errorf("error"))
	actualError := suite.usecase.DeleteUser("1xxx")
	assert.Error(suite.T(), actualError)
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
