package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type statusRepoMock struct {
	mock.Mock
}

func (s *statusRepoMock) Create(payload model.Status) error {
	args := s.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *statusRepoMock) Delete(id string) error {
	args := s.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *statusRepoMock) Get(id string) (model.Status, error) {
	args := s.Called(id)
	if args.Get(1) != nil {
		return model.Status{}, args.Error(1)
	}
	return args.Get(0).(model.Status), nil
}

func (s *statusRepoMock) GetByName(name string) (model.Status, error) {
	args := s.Called(name)
	if args.Get(1) != nil {
		return model.Status{}, args.Error(1)
	}
	return args.Get(0).(model.Status), nil

}

func (s *statusRepoMock) List() ([]model.Status, error) {
	panic("unimplemented")
}

func (s *statusRepoMock) Paging(requestPaging dto.PaginationParam) ([]model.Status, dto.Paging, error) {
	args := s.Called(requestPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Status), args.Get(1).(dto.Paging), nil

}

func (s *statusRepoMock) Update(payload model.Status) error {
	args := s.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

type StatusUseCaseTestSuite struct {
	suite.Suite
	statusRepoMock *statusRepoMock
	usecase        StatusUseCase
}

func (suite *StatusUseCaseTestSuite) SetupTest() {
	suite.statusRepoMock = new(statusRepoMock)
	suite.usecase = NewStatusUseCase(suite.statusRepoMock)
}

var statusDummy = []model.Status{
	{
		StatusId: "1",
		Name:     "Status 1",
	},
	{
		StatusId: "2",
		Name:     "Status 2",
	},
	{
		StatusId: "3",
		Name:     "Status 3",
	},
}

func (suite *StatusUseCaseTestSuite) TestRegisterNewStatus_Succes() {
	newStatus := statusDummy[0]
	suite.statusRepoMock.On("GetByName", newStatus.Name).Return(model.Status{}, fmt.Errorf("error"))
	suite.statusRepoMock.On("Create", newStatus).Return(nil)

	err := suite.usecase.RegisterNewStatus(newStatus)

	assert.Nil(suite.T(), err)
}

func (suite *StatusUseCaseTestSuite) TestRegisterNewStatus_EmptyName() {
	emptyStatus := model.Status{Name: ""}

	err := suite.usecase.RegisterNewStatus(emptyStatus)

	expectedError := "name  required fields"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *StatusUseCaseTestSuite) TestRegisterNewStatus_DuplicateName() {
	duplicateStatus := model.Status{Name: "Duplicate Status"}
	suite.statusRepoMock.On("GetByName", duplicateStatus.Name).Return(duplicateStatus, nil)

	err := suite.usecase.RegisterNewStatus(duplicateStatus)

	expectedError := fmt.Sprintf("ERR status with name %s exits", duplicateStatus.Name)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *StatusUseCaseTestSuite) TestRegisterNewStatus_Failure() {
	newStatus := model.Status{Name: "New Status"}
	suite.statusRepoMock.On("GetByName", newStatus.Name).Return(model.Status{}, fmt.Errorf("error"))
	suite.statusRepoMock.On("Create", newStatus).Return(fmt.Errorf("error"))

	err := suite.usecase.RegisterNewStatus(newStatus)

	expectedError := fmt.Sprintf("failed to create new status: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)

}

func (suite *StatusUseCaseTestSuite) TestFindAllStatus_Success() {
	dummy := statusDummy
	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   3,
		TotalPages:  1,
	}
	requestPaging := dto.PaginationParam{
		Page: 1,
	}
	suite.statusRepoMock.On("Paging", requestPaging).Return(dummy, expectedPaging, nil)
	statuses, paging, err := suite.usecase.FindAllStatus(requestPaging)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummy, statuses)
	assert.Equal(suite.T(), expectedPaging, paging)
}
func (suite *StatusUseCaseTestSuite) TestFindAllStatus_Failure() {
	requestPaging := dto.PaginationParam{
		Page: 1,
	}
	suite.statusRepoMock.On("Paging", requestPaging).Return(nil, dto.Paging{}, fmt.Errorf("error"))

	statuses, paging, err := suite.usecase.FindAllStatus(requestPaging)

	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "error")
	assert.Empty(suite.T(), statuses)
	assert.Empty(suite.T(), paging)

}
func (suite *StatusUseCaseTestSuite) TestUpdateStatus_Success() {
	dummy := statusDummy[0]
	suite.statusRepoMock.On("Update", dummy).Return(nil)

	err := suite.usecase.UpdateStatus(dummy)

	assert.Nil(suite.T(), err)

	suite.statusRepoMock.AssertCalled(suite.T(), "Update", dummy)
}

func (suite *StatusUseCaseTestSuite) TestUpdateStatus_EmptyName() {
	emptyStatus := model.Status{Name: ""}

	err := suite.usecase.UpdateStatus(emptyStatus)

	expectedError := "name  required fields"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *StatusUseCaseTestSuite) TestUpdateStatus_Failure() {
	dummy := statusDummy[0]
	suite.statusRepoMock.On("Update", dummy).Return(fmt.Errorf("error"))

	err := suite.usecase.UpdateStatus(dummy)

	expectedError := fmt.Sprintf("failed to update status: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)

}

func (suite *StatusUseCaseTestSuite) TestFindByIdStatus_Success() {
	dummyStatus := statusDummy[0]
	suite.statusRepoMock.On("Get", dummyStatus.StatusId).Return(dummyStatus, nil)

	status, err := suite.usecase.FindByIdStatus(dummyStatus.StatusId)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyStatus, status)
}

func (suite *StatusUseCaseTestSuite) TestFindByIdStatus_NotFound() {
	suite.statusRepoMock.On("Get", "1234").Return(model.Status{}, fmt.Errorf("error"))

	status, err := suite.usecase.FindByIdStatus("1234")

	expectedError := "status with id 1234 not found"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
	assert.Equal(suite.T(), model.Status{}, status)
}
func (suite *StatusUseCaseTestSuite) TestDeleteStatus_Success() {
	dummyStatus := statusDummy[0]
	suite.statusRepoMock.On("Get", dummyStatus.StatusId).Return(dummyStatus, nil)
	suite.statusRepoMock.On("Delete", dummyStatus.StatusId).Return(nil)

	err := suite.usecase.DeleteStatus(dummyStatus.StatusId)

	assert.Nil(suite.T(), err)
}

func (suite *StatusUseCaseTestSuite) TestDeleteStatus_StatusNotFound() {
	nonExistentStatusID := "1234"
	suite.statusRepoMock.On("Get", nonExistentStatusID).Return(model.Status{}, fmt.Errorf("error"))

	err := suite.usecase.DeleteStatus(nonExistentStatusID)

	expectedError := fmt.Sprintf("status with ID %s not found", nonExistentStatusID)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *StatusUseCaseTestSuite) TestDeleteStatus_Failure() {
	dummyStatus := statusDummy[0]
	suite.statusRepoMock.On("Get", dummyStatus.StatusId).Return(dummyStatus, nil)
	suite.statusRepoMock.On("Delete", dummyStatus.StatusId).Return(fmt.Errorf("error"))

	err := suite.usecase.DeleteStatus(dummyStatus.StatusId)

	expectedError := fmt.Sprintf("failed to delete status: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func TestStatusUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(StatusUseCaseTestSuite))
}
