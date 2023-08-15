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

type statusUseCaseMock struct {
	mock.Mock
}

// DeleteStatus implements usecase.StatusUseCase.
func (s *statusUseCaseMock) DeleteStatus(id string) error {
	args := s.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// FindAllStatus implements usecase.StatusUseCase.
func (s *statusUseCaseMock) FindAllStatus(requesPaging dto.PaginationParam) ([]model.Status, dto.Paging, error) {
	args := s.Called(requesPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Status), args.Get(1).(dto.Paging), nil

}

// FindByIdStatus implements usecase.StatusUseCase.
func (s *statusUseCaseMock) FindByIdStatus(id string) (model.Status, error) {
	args := s.Called(id)
	if args.Get(1) != nil {
		return model.Status{}, args.Error(1)
	}
	return args.Get(0).(model.Status), nil

}

// RegisterNewStatus implements usecase.StatusUseCase.
func (s *statusUseCaseMock) RegisterNewStatus(payload model.Status) error {
	args := s.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

// UpdateStatus implements usecase.StatusUseCase.
func (s *statusUseCaseMock) UpdateStatus(payload model.Status) error {
	args := s.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

type StatusControllerTestSuite struct {
	suite.Suite
	usecaseMock *statusUseCaseMock
	router      *gin.Engine
}

func (suite *StatusControllerTestSuite) SetupTest() {
	suite.usecaseMock = new(statusUseCaseMock)
	suite.router = gin.Default()
	NewStatusController(suite.router, suite.usecaseMock)
}

func TestStatusControllerTestSuite(t *testing.T) {
	suite.Run(t, new(StatusControllerTestSuite))
}

func (suite *StatusControllerTestSuite) TestCreateHandler_Success() {
	payload := model.Status{
		Name: "New Status",
	}
	suite.usecaseMock.On("RegisterNewStatus", mock.MatchedBy(func(arg model.Status) bool {
		return arg.Name == payload.Name
	})).Return(nil)
	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/statuss", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

func (suite *StatusControllerTestSuite) TestCreateHandler_BadRequest() {
	invalidPayload := "invalid JSON"

	req := httptest.NewRequest("POST", "/api/v1/statuss", strings.NewReader(invalidPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "err")
}

func (suite *StatusControllerTestSuite) TestCreateHandler_InternalServerError() {
	payload := model.Status{
		StatusId: "ID1",
		Name:     "New Status",
	}

	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("RegisterNewStatus", payload).Return(expectedError)

	requestBody, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/statuss", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *StatusControllerTestSuite) TestListHandler_Success() {
	// Set up mock behavior
	paginationParam := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}
	expectedStatuses := []model.Status{
		{StatusId: "1", Name: "Status 1"},
		{StatusId: "2", Name: "Status 2"},
	}
	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   len(expectedStatuses),
		TotalPages:  1,
	}
	suite.usecaseMock.On("FindAllStatus", paginationParam).Return(expectedStatuses, expectedPaging, nil)

	req := httptest.NewRequest("GET", "/api/v1/statuss?page=1&limit=5", nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *StatusControllerTestSuite) TestListHandler_InternalServerError() {
	paginationParam := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("FindAllStatus", paginationParam).Return(nil, dto.Paging{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/statuss?page=1&limit=5", nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *StatusControllerTestSuite) TestGetHandler_Success() {
	statusID := "1"
	expectedStatus := model.Status{StatusId: statusID, Name: "Status 1"}
	suite.usecaseMock.On("FindByIdStatus", statusID).Return(expectedStatus, nil)

	req := httptest.NewRequest("GET", "/api/v1/statuss/"+statusID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *StatusControllerTestSuite) TestGetHandler_InternalServerError() {
	statusID := "1"
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("FindByIdStatus", statusID).Return(model.Status{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/statuss/"+statusID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *StatusControllerTestSuite) TestUpdateHandler_Success() {
	expectedStatus := model.Status{StatusId: "1", Name: "Updated Status"}
	suite.usecaseMock.On("UpdateStatus", expectedStatus).Return(nil)

	requestBody := []byte(`{"id": "1", "name": "Updated Status"}`)
	req := httptest.NewRequest("PUT", "/api/v1/statuss", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *StatusControllerTestSuite) TestUpdateHandler_BadRequest() {
	req := httptest.NewRequest("PUT", "/api/v1/statuss", strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *StatusControllerTestSuite) TestUpdateHandler_InternalServerError() {
	expectedStatus := model.Status{StatusId: "1", Name: "Updated Status"}
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("UpdateStatus", expectedStatus).Return(expectedError)

	requestBody := []byte(`{"id": "1", "name": "Updated Status"}`)
	req := httptest.NewRequest("PUT", "/api/v1/statuss", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *StatusControllerTestSuite) TestDeleteHandler_Success() {
	statusID := "1"
	suite.usecaseMock.On("DeleteStatus", statusID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/v1/statuss/"+statusID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *StatusControllerTestSuite) TestDeleteHandler_InternalServerError() {
	statusID := "1"
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("DeleteStatus", statusID).Return(expectedError)

	req := httptest.NewRequest("DELETE", "/api/v1/statuss/"+statusID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
