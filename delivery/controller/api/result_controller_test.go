package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type resultUseCaseMock struct {
	mock.Mock
}

// DeleteResult implements usecase.ResultUseCase.
func (r *resultUseCaseMock) DeleteResult(id string) error {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// FindAllResult implements usecase.ResultUseCase.
func (r *resultUseCaseMock) FindAllResult(requesPaging dto.PaginationParam) ([]model.Result, dto.Paging, error) {
	args := r.Called(requesPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Result), args.Get(1).(dto.Paging), nil

}

// FindByIdResult implements usecase.ResultUseCase.
func (r *resultUseCaseMock) FindByIdResult(id string) (model.Result, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.Result{}, args.Error(1)
	}
	return args.Get(0).(model.Result), nil

}

// RegisterNewResult implements usecase.ResultUseCase.
func (r *resultUseCaseMock) RegisterNewResult(payload model.Result) error {
	args := r.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

// UpdateResult implements usecase.ResultUseCase.
func (r *resultUseCaseMock) UpdateResult(payload model.Result) error {
	args := r.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

type ResultControllerTestSuite struct {
	suite.Suite
	usecaseMock *resultUseCaseMock
	router      *gin.Engine
}

func (suite *ResultControllerTestSuite) SetupTest() {
	suite.usecaseMock = new(resultUseCaseMock)
	suite.router = gin.Default()
	NewResultController(suite.router, suite.usecaseMock)
}

func TestResultControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ResultControllerTestSuite))
}

func (suite *ResultControllerTestSuite) TestCreateHandler_Success() {
	payload := model.Result{
		Name: "New Result",
	}
	suite.usecaseMock.On("RegisterNewResult", mock.MatchedBy(func(arg model.Result) bool {
		return arg.Name == payload.Name
	})).Return(nil)
	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/results", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

func (suite *ResultControllerTestSuite) TestCreateHandler_BadRequest() {
	invalidPayload := "invalid JSON"

	req := httptest.NewRequest("POST", "/api/v1/results", strings.NewReader(invalidPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "err")
}

func (suite *ResultControllerTestSuite) TestCreateHandler_InternalServerError() {
	payload := model.Result{
		Name: "New Result",
	}

	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("RegisterNewResult", mock.MatchedBy(func(arg model.Result) bool {
		return arg.Name == payload.Name
	})).Return(expectedError)

	requestBody, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/results", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	// Verify the mock method call with any ResultId value
	suite.usecaseMock.AssertCalled(suite.T(), "RegisterNewResult", mock.MatchedBy(func(arg model.Result) bool {
		return arg.Name == payload.Name
	}))
}


func (suite *ResultControllerTestSuite) TestListHandler_Success() {
	// Set up mock behavior
	paginationParam := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}
	expectedResultes := []model.Result{
		{ResultId: "1", Name: "Result 1"},
		{ResultId: "2", Name: "Result 2"},
	}
	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   len(expectedResultes),
		TotalPages:  1,
	}
	suite.usecaseMock.On("FindAllResult", paginationParam).Return(expectedResultes, expectedPaging, nil)

	req := httptest.NewRequest("GET", "/api/v1/results?page=1&limit=5", nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *ResultControllerTestSuite) TestListHandler_InternalServerError() {
	paginationParam := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("FindAllResult", paginationParam).Return(nil, dto.Paging{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/results?page=1&limit=5", nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *ResultControllerTestSuite) TestGetHandler_Success() {
	resultID := "1"
	expectedResult := model.Result{ResultId: resultID, Name: "Result 1"}
	suite.usecaseMock.On("FindByIdResult", resultID).Return(expectedResult, nil)

	req := httptest.NewRequest("GET", "/api/v1/results/"+resultID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *ResultControllerTestSuite) TestGetHandler_InternalServerError() {
	resultID := "1"
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("FindByIdResult", resultID).Return(model.Result{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/results/"+resultID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *ResultControllerTestSuite) TestUpdateHandler_Success() {
	expectedResult := model.Result{ResultId: "1", Name: "Updated Result"}
	suite.usecaseMock.On("UpdateResult", expectedResult).Return(nil)

	requestBody := []byte(`{"id": "1", "name": "Updated Result"}`)
	req := httptest.NewRequest("PUT", "/api/v1/results", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *ResultControllerTestSuite) TestUpdateHandler_BadRequest() {
	req := httptest.NewRequest("PUT", "/api/v1/results", strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ResultControllerTestSuite) TestUpdateHandler_InternalServerError() {
	expectedResult := model.Result{ResultId: "1", Name: "Updated Result"}
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("UpdateResult", expectedResult).Return(expectedError)

	requestBody := []byte(`{"id": "1", "name": "Updated Result"}`)
	req := httptest.NewRequest("PUT", "/api/v1/results", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *ResultControllerTestSuite) TestDeleteHandler_Success() {
	resultID := "1"
	suite.usecaseMock.On("DeleteResult", resultID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/v1/results/"+resultID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *ResultControllerTestSuite) TestDeleteHandler_InternalServerError() {
	resultID := "1"
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("DeleteResult", resultID).Return(expectedError)

	req := httptest.NewRequest("DELETE", "/api/v1/results/"+resultID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
