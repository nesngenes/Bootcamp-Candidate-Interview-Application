package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type interviewerRepoMock struct {
	mock.Mock
}

// Create implements repository.InterviewerRepository.
func (i *interviewerRepoMock) Create(payload model.Interviewer) error {
	args := i.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// Delete implements repository.InterviewerRepository.
func (i *interviewerRepoMock) Delete(id string) error {
	args := i.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// Get implements repository.InterviewerRepository.
func (i *interviewerRepoMock) Get(id string) (model.Interviewer, error) {
	args := i.Called(id)
	if args.Get(1) != nil {
		return model.Interviewer{}, args.Error(1)
	}
	return args.Get(0).(model.Interviewer), nil
}

// List implements repository.InterviewerRepository.
func (i *interviewerRepoMock) List() ([]model.Interviewer, error) {
	args := i.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Interviewer), nil
}

// Update implements repository.InterviewerRepository.
func (i *interviewerRepoMock) Update(payload model.Interviewer) error {
	args := i.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

type InterviewerUseCaseTestSuite struct {
	suite.Suite
	interviewerRepoMock *interviewerRepoMock
	usecase             InterviewerUseCase
}

func (suite *InterviewerUseCaseTestSuite) SetupTest() {
	suite.interviewerRepoMock = new(interviewerRepoMock)
	suite.usecase = NewInterviewerUseCase(suite.interviewerRepoMock)
}

func TestInterviewerUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(InterviewerUseCaseTestSuite))
}

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

func (suite *InterviewerUseCaseTestSuite) TestRegisterNewInterviewer_Success() {
	dummyInterviewer := dummyInterviewers[0]
	suite.interviewerRepoMock.On("Create", dummyInterviewer).Return(nil)

	err := suite.usecase.RegisterNewInterviewer(dummyInterviewer)

	assert.Nil(suite.T(), err)
}

func (suite *InterviewerUseCaseTestSuite) TestRegisterNewInterviewer_MissingFields() {
	dummyInterviewer := model.Interviewer{}

	err := suite.usecase.RegisterNewInterviewer(dummyInterviewer)

	assert.NotNil(suite.T(), err)
}

func (suite *InterviewerUseCaseTestSuite) TestRegisterNewInterviewer_Failure() {
	dummyInterviewer := dummyInterviewers[0]
	suite.interviewerRepoMock.On("Create", dummyInterviewer).Return(fmt.Errorf("error"))

	err := suite.usecase.RegisterNewInterviewer(dummyInterviewer)

	assert.NotNil(suite.T(), err)
}

func (suite *InterviewerUseCaseTestSuite) TestFindAllInterviewer_Success() {
	dummy := dummyInterviewers
	suite.interviewerRepoMock.On("List").Return(dummy, nil)

	interviewers, err := suite.usecase.FindAllInterviewer()

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummy, interviewers)
}

func (suite *InterviewerUseCaseTestSuite) TestFindAllInterviewer_Failure() {
	suite.interviewerRepoMock.On("List").Return(nil, fmt.Errorf("error"))

	interviewers, err := suite.usecase.FindAllInterviewer()

	assert.NotNil(suite.T(), err)
	assert.Empty(suite.T(), interviewers)
}

func (suite *InterviewerUseCaseTestSuite) TestFindByIdInterviewer_Success() {
	dummyInterviewer := dummyInterviewers[0]
	suite.interviewerRepoMock.On("Get", "1").Return(dummyInterviewer, nil)

	interviewer, err := suite.usecase.FindByIdInterviewer("1")

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyInterviewer, interviewer)
}

func (suite *InterviewerUseCaseTestSuite) TestFindByIdInterviewer_NotFound() {
	suite.interviewerRepoMock.On("Get", "2384").Return(model.Interviewer{}, fmt.Errorf("error"))

	interviewer, err := suite.usecase.FindByIdInterviewer("2384")

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), model.Interviewer{}, interviewer)
}

func (suite *InterviewerUseCaseTestSuite) TestDeleteInterviewer_Success() {
	dummyInterviewer := dummyInterviewers[0]
	suite.interviewerRepoMock.On("Get", "1").Return(dummyInterviewer, nil)
	suite.interviewerRepoMock.On("Delete", "1").Return(nil)

	err := suite.usecase.DeleteInterviewer("1")

	assert.Nil(suite.T(), err)
}

func (suite *InterviewerUseCaseTestSuite) TestDeleteInterviewer_NotFound() {
	suite.interviewerRepoMock.On("Get", "2").Return(model.Interviewer{}, fmt.Errorf("error"))

	err := suite.usecase.DeleteInterviewer("2")

	assert.NotNil(suite.T(), err)
}

func (suite *InterviewerUseCaseTestSuite) TestDeleteInterviewer_DeleteError() {
	dummyInterviewer := dummyInterviewers[0]
	suite.interviewerRepoMock.On("Get", "1").Return(dummyInterviewer, nil)

	suite.interviewerRepoMock.On("Delete", "1").Return(fmt.Errorf("error"))

	err := suite.usecase.DeleteInterviewer("1")

	assert.NotNil(suite.T(), err)
}

func (suite *InterviewerUseCaseTestSuite) TestUpdateInterviewer_Success() {
	dummyInterviewer := dummyInterviewers[0]
	suite.interviewerRepoMock.On("Update", dummyInterviewer).Return(nil)

	err := suite.usecase.UpdateInterviewer(dummyInterviewer)

	assert.Nil(suite.T(), err)
}

func (suite *InterviewerUseCaseTestSuite) TestUpdateInterviewer_InvalidFields() {
	invalidInterviewer := model.Interviewer{
		InterviewerID: "",
		FullName:      "",
		UserID:        "",
	}

	err := suite.usecase.UpdateInterviewer(invalidInterviewer)

	assert.NotNil(suite.T(), err)
}

func (suite *InterviewerUseCaseTestSuite) TestUpdateInterviewer_Fail() {
	dummyInterviewer := dummyInterviewers[0]
	suite.interviewerRepoMock.On("Update", dummyInterviewer).Return(fmt.Errorf("error"))

	err := suite.usecase.UpdateInterviewer(dummyInterviewer)

	assert.NotNil(suite.T(), err)
}
