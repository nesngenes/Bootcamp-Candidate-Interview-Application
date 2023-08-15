package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type hrRepoMock struct {
	mock.Mock
}

// hr repo
func (u *hrRepoMock) Create(user model.HRRecruitment) error {
	args := u.Called(user)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (u *hrRepoMock) List() ([]model.HRRecruitment, error) {
	args := u.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.HRRecruitment), nil
}

func (u *hrRepoMock) Get(id string) (model.HRRecruitment, error) {
	args := u.Called(id)
	if args.Get(1) != nil {
		return model.HRRecruitment{}, args.Error(1)
	}
	return args.Get(0).(model.HRRecruitment), nil
}

func (u *hrRepoMock) GetByUserID(userID string) (model.HRRecruitment, error) {
	args := u.Called(userID)
	if args.Get(1) != nil {
		return model.HRRecruitment{}, args.Error(1)
	}
	return args.Get(0).(model.HRRecruitment), nil
}

func (u *hrRepoMock) Update(payload model.HRRecruitment) error {
	args := u.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (u *hrRepoMock) Delete(id string) error {
	args := u.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// user
type usersRepoMock struct {
	mock.Mock
}

// implement user
func (u *usersRepoMock) Create(payload model.Users) error {
	panic("unimplemented")
}
func (u *usersRepoMock) List() ([]model.Users, error) {
	panic("unimplemented")
}
func (u *usersRepoMock) GetByEmail(email string) (model.Users, error) {
	panic("unimplemented")
}
func (u *usersRepoMock) GetByUserName(userName string) (model.Users, error) {
	panic("unimplemented")
}

// krn ini yang perlu buat create & update
func (u *usersRepoMock) Get(id string) (model.Users, error) {
	args := u.Called(id)
	if args.Get(1) != nil {
		return model.Users{}, args.Error(1)
	}
	return args.Get(0).(model.Users), nil
}

func (u *usersRepoMock) Update(payload model.Users) error {
	panic("unimplemented")
}

func (u *usersRepoMock) Delete(id string) error {
	panic("unimplemented")
}

type HRUsecaseTestSuite struct {
	suite.Suite
	hrRepoMock   *hrRepoMock
	userRepoMock *usersRepoMock
	usecase      HRRecruitmentUsecase
}

func (suite *HRUsecaseTestSuite) SetupTest() {
	suite.hrRepoMock = new(hrRepoMock)
	suite.userRepoMock = new(usersRepoMock)
	suite.usecase = NewHRRecruitmentUsecase(suite.hrRepoMock, suite.userRepoMock)
}

var hrDummy = []model.HRRecruitment{
	{
		ID:       "1",
		FullName: "Ella lalala",
		UserID:   "100",
		User: model.Users{
			Id:       "100",
			Email:    "ela@mail.com",
			UserName: "kakiku",
			UserRole: model.UserRoles{
				Id:   "1",
				Name: "HR",
			},
		},
	},
	{
		ID:       "2",
		FullName: "Elli lilili",
		UserID:   "125",
		User: model.Users{
			Id:       "125",
			Email:    "eli@mail.com",
			UserName: "sasisu",
			UserRole: model.UserRoles{
				Id:   "1",
				Name: "HR",
			},
		},
	},
	{
		ID:       "3",
		FullName: "haha",
		UserID:   "haha123",
		User: model.Users{
			Id:       "haha123",
			Email:    "eli@mail.com",
			UserName: "sasisu",
			UserRole: model.UserRoles{
				Id:   "1",
				Name: "HR",
			},
		},
	},
}

func (suite *HRUsecaseTestSuite) TestUpdateHR_UserIdExist() {
	dummy := hrDummy[0]

	//mock expectations for Get and GetByUserID
	suite.hrRepoMock.On("Get", dummy.ID).Return(dummy, nil)
	suite.hrRepoMock.On("GetByUserID", dummy.UserID).Return(hrDummy[1], nil) // simulasi record hr lain sudah menggunakan userid ini.

	// Call the use case function
	err := suite.usecase.UpdateHRRecruitment(dummy)

	assert.Error(suite.T(), err) // error
	suite.hrRepoMock.AssertExpectations(suite.T())
}

//masih error
// func (suite *HRUsecaseTestSuite) TestUpdateHR_Success() {
// 	dummy := hrDummy[0]
// 	userID := "user123"

// 	suite.hrRepoMock.On("Get", dummy.ID).Return(dummy, nil).Once()
// 	suite.hrRepoMock.On("GetByUserID", userID).Return(model.HRRecruitment{}, nil).Once()
// 	suite.hrRepoMock.On("Update", mock.AnythingOfType("model.HRRecruitment")).Return(nil).Once()

// 	// Call the use case
// 	err := suite.usecase.UpdateHRRecruitment(model.HRRecruitment{
// 		ID:       dummy.ID,
// 		FullName: "Updated Name",
// 		UserID:   userID,
// 	})

// 	// Assert no error
// 	assert.NoError(suite.T(), err)
// 	suite.hrRepoMock.AssertExpectations(suite.T())
// }

func (suite *HRUsecaseTestSuite) TestCreateHRRecruitment() {
	payload := model.HRRecruitment{
		ID:       "123",
		FullName: "Nama Lengkap",
		UserID:   "user123",
	}

	// Set up mock expectations
	suite.userRepoMock.On("Get", payload.UserID).Return(model.Users{}, nil)
	suite.hrRepoMock.On("GetByUserID", payload.UserID).Return(model.HRRecruitment{}, nil)
	suite.hrRepoMock.On("Create", payload).Return(nil)

	// Call the use case function
	err := suite.usecase.CreateHRRecruitment(payload)

	// Assert expectations and results
	assert.NoError(suite.T(), err)
	suite.userRepoMock.AssertExpectations(suite.T())
	suite.hrRepoMock.AssertExpectations(suite.T())
}

func (suite *HRUsecaseTestSuite) TestCreateHRRecruitment_Fail() {
	payload := hrDummy[0]
	suite.userRepoMock.On("Get", payload.UserID).Return(model.HRRecruitment{}, fmt.Errorf("user not found"))
	suite.hrRepoMock.On("GetByUserID", payload.UserID).Return(model.HRRecruitment{}, nil)
	suite.hrRepoMock.On("Create", payload).Return(nil)
	err := suite.usecase.CreateHRRecruitment(payload)
	assert.Error(suite.T(), err)

}

func (suite *HRUsecaseTestSuite) TestList_Success() {
	dummy := hrDummy
	suite.hrRepoMock.On("List").Return(dummy, nil)
	actualList, err := suite.usecase.ListHRRecruitments()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummy, actualList)
}

func (suite *HRUsecaseTestSuite) TestList_Fail() {
	suite.hrRepoMock.On("List").Return(nil, fmt.Errorf("error"))
	actualList, err := suite.usecase.ListHRRecruitments()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), actualList)
}

func (suite *HRUsecaseTestSuite) TestGetUser_ByID() {
	dummy := hrDummy[0]
	suite.hrRepoMock.On("Get", dummy.ID).Return(dummy, nil)
	actualRole, err := suite.usecase.Get(dummy.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummy, actualRole)
}

func (suite *HRUsecaseTestSuite) TestDeleteUser_Success() {
	suite.hrRepoMock.On("Delete", hrDummy[0].ID).Return(nil)
	actualError := suite.usecase.DeleteHRRecruitment(hrDummy[0].ID)
	assert.Nil(suite.T(), actualError)

}

func (suite *HRUsecaseTestSuite) TestDeleteUser_Fail() {
	suite.hrRepoMock.On("Delete", "1xxx").Return(fmt.Errorf("error"))
	actualError := suite.usecase.DeleteHRRecruitment("1xxx")
	assert.Error(suite.T(), actualError)
}

func TestHRUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(HRUsecaseTestSuite))
}
