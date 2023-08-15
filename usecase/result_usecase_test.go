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

type resultRepoMock struct {
	mock.Mock
}

func (r *resultRepoMock) Create(payload model.Result) error {
	args := r.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *resultRepoMock) Delete(id string) error {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *resultRepoMock) Get(id string) (model.Result, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.Result{}, args.Error(1)
	}
	return args.Get(0).(model.Result), nil
}

func (r *resultRepoMock) GetByName(name string) (model.Result, error) {
	args := r.Called(name)
	if args.Get(1) != nil {
		return model.Result{}, args.Error(1)
	}
	return args.Get(0).(model.Result), nil

}

func (r *resultRepoMock) List() ([]model.Result, error) {
	panic("unimplemented")
}

func (r *resultRepoMock) Paging(requestPaging dto.PaginationParam) ([]model.Result, dto.Paging, error) {
	args := r.Called(requestPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Result), args.Get(1).(dto.Paging), nil

}

func (r *resultRepoMock) Update(payload model.Result) error {
	args := r.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

type ResultUseCaseTestSuite struct {
	suite.Suite
	resultRepoMock *resultRepoMock
	usecase        ResultUseCase
}

func (suite *ResultUseCaseTestSuite) SetupTest() {
	suite.resultRepoMock = new(resultRepoMock)
	suite.usecase = NewResultUseCase(suite.resultRepoMock)
}

var resultDummy = []model.Result {
	{
		ResultId: "1", 
		Name: "Name 1",     
	},
	{
		ResultId: "2", 
		Name: "Name 2",
	},
	{
		ResultId: "3", 
		Name: "Name 3",
	},
}

func (suite *ResultUseCaseTestSuite) TestRegisterNewResult_Succes() {
	newResult := resultDummy[0]
	suite.resultRepoMock.On("GetByName", newResult.Name).Return(model.Result{}, fmt.Errorf("error"))
	suite.resultRepoMock.On("Create", newResult).Return(nil)

	err := suite.usecase.RegisterNewResult(newResult)

	assert.Nil(suite.T(), err)
}

func (suite *ResultUseCaseTestSuite) TestRegisterNewResult_EmptyName() {
	emptyResult := model.Result{Name: ""}

	err := suite.usecase.RegisterNewResult(emptyResult)

	expectedError := "name  required fields"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *ResultUseCaseTestSuite) TestRegisterNewResult_DuplicateName() {
	duplicateResult := model.Result{Name: "Duplicate Result"}
	suite.resultRepoMock.On("GetByName", duplicateResult.Name).Return(duplicateResult, nil)

	err := suite.usecase.RegisterNewResult(duplicateResult)

	expectedError := fmt.Sprintf("ERR result with name %s exits", duplicateResult.Name)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *ResultUseCaseTestSuite) TestRegisterNewResult_Failure() {
	newResult := model.Result{Name: "New Result"}
	suite.resultRepoMock.On("GetByName", newResult.Name).Return(model.Result{}, fmt.Errorf("error"))
	suite.resultRepoMock.On("Create", newResult).Return(fmt.Errorf("error"))

	err := suite.usecase.RegisterNewResult(newResult)

	expectedError := fmt.Sprintf("failed to create new result: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)

}

func (suite *ResultUseCaseTestSuite) TestFindAllResult_Success() {
	dummy := resultDummy
	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   3,
		TotalPages:  1,
	}
	requestPaging := dto.PaginationParam{
		Page: 1,
	}
	suite.resultRepoMock.On("Paging", requestPaging).Return(dummy, expectedPaging, nil)
	results, paging, err := suite.usecase.FindAllResult(requestPaging)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummy, results)
	assert.Equal(suite.T(), expectedPaging, paging)
}

func (suite *ResultUseCaseTestSuite) TestFindAllResult_Failure() {
	requestPaging := dto.PaginationParam{
		Page: 1,
	}
	suite.resultRepoMock.On("Paging", requestPaging).Return(nil, dto.Paging{}, fmt.Errorf("error"))

	results, paging, err := suite.usecase.FindAllResult(requestPaging)

	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "error")
	assert.Empty(suite.T(), results)
	assert.Empty(suite.T(), paging)

}

func (suite *ResultUseCaseTestSuite) TestUpdateResult_Success() {
	dummy := resultDummy[0]
	suite.resultRepoMock.On("Update", dummy).Return(nil)

	err := suite.usecase.UpdateResult(dummy)

	assert.Nil(suite.T(), err)

	suite.resultRepoMock.AssertCalled(suite.T(), "Update", dummy)
}

func (suite *ResultUseCaseTestSuite) TestUpdateResult_EmptyName() {
	emptyResult := model.Result{Name: ""}

	err := suite.usecase.UpdateResult(emptyResult)

	expectedError := "name  required fields"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *ResultUseCaseTestSuite) TestUpdateResult_Failure() {
	dummy := resultDummy[0]
	suite.resultRepoMock.On("Update", dummy).Return(fmt.Errorf("error"))

	err := suite.usecase.UpdateResult(dummy)

	expectedError := fmt.Sprintf("failed to update result: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *ResultUseCaseTestSuite) TestFindByIdResult_Success() {
	dummyResult := resultDummy[0]
	suite.resultRepoMock.On("Get", dummyResult.ResultId).Return(dummyResult, nil)

	result, err := suite.usecase.FindByIdResult(dummyResult.ResultId)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyResult, result)
}

func (suite *ResultUseCaseTestSuite) TestFindByIdResult_NotFound() {
	suite.resultRepoMock.On("Get", "1234").Return(model.Result{}, fmt.Errorf("error"))

	result, err := suite.usecase.FindByIdResult("1234")

	expectedError := "result with id 1234 not found"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
	assert.Equal(suite.T(), model.Result{}, result)
}

func (suite *ResultUseCaseTestSuite) TestDeleteResult_Success() {
	dummyResult := resultDummy[0]
	suite.resultRepoMock.On("Get", dummyResult.ResultId).Return(dummyResult, nil)
	suite.resultRepoMock.On("Delete", dummyResult.ResultId).Return(nil)

	err := suite.usecase.DeleteResult(dummyResult.ResultId)

	assert.Nil(suite.T(), err)
}

func (suite *ResultUseCaseTestSuite) TestDeleteResult_ResultNotFound() {
	nonExistentResultID := "1234"
	suite.resultRepoMock.On("Get", nonExistentResultID).Return(model.Result{}, fmt.Errorf("error"))

	err := suite.usecase.DeleteResult(nonExistentResultID)

	expectedError := fmt.Sprintf("result with ID %s not found", nonExistentResultID)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *ResultUseCaseTestSuite) TestDeleteResult_Failure() {
	dummyResult := resultDummy[0]
	suite.resultRepoMock.On("Get", dummyResult.ResultId).Return(dummyResult, nil)
	suite.resultRepoMock.On("Delete", dummyResult.ResultId).Return(fmt.Errorf("error"))

	err := suite.usecase.DeleteResult(dummyResult.ResultId)

	expectedError := fmt.Sprintf("failed to delete result: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func TestResultUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ResultUseCaseTestSuite))
}

