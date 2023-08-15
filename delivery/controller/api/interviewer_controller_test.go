package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"interview_bootcamp/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type interviewerUsecaseMock struct {
	mock.Mock
}

// DeleteInterviewer implements usecase.InterviewerUseCase.
func (i *interviewerUsecaseMock) DeleteInterviewer(id string) error {
	args := i.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// FindAllInterviewer implements usecase.InterviewerUseCase.
func (i *interviewerUsecaseMock) FindAllInterviewer() ([]model.Interviewer, error) {
	args := i.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Interviewer), nil
}

// FindByIdInterviewer implements usecase.InterviewerUseCase.
func (i *interviewerUsecaseMock) FindByIdInterviewer(id string) (model.Interviewer, error) {
	args := i.Called(id)
	if args.Get(1) != nil {
		return model.Interviewer{}, args.Error(1)
	}
	return args.Get(0).(model.Interviewer), nil
}

// RegisterNewInterviewer implements usecase.InterviewerUseCase.
func (i *interviewerUsecaseMock) RegisterNewInterviewer(payload model.Interviewer) error {
	args := i.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// UpdateInterviewer implements usecase.InterviewerUseCase.
func (i *interviewerUsecaseMock) UpdateInterviewer(payload model.Interviewer) error {
	args := i.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

type InterviewerControllerTestSuite struct {
	suite.Suite
	interviewerUsecaseMock *interviewerUsecaseMock
	router                 *gin.Engine
}

func (suite *InterviewerControllerTestSuite) SetupTest() {
	suite.interviewerUsecaseMock = new(interviewerUsecaseMock)
	suite.router = gin.Default()
	NewInterviewerController(suite.router, suite.interviewerUsecaseMock)
}

func TestInterviewerControllerTestSuite(t *testing.T) {
	suite.Run(t, new(InterviewerControllerTestSuite))
}

func (suite *InterviewerControllerTestSuite) TestCreateHandler_Success() {
	payload := model.Interviewer{
		InterviewerID: "1",
		FullName:      "John Doe",
		UserID:        "1",
	}
	suite.interviewerUsecaseMock.On("RegisterNewInterviewer", mock.MatchedBy(func(args model.Interviewer) bool {
		return args.FullName == payload.FullName
	})).Return(nil)
	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/interviewers", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
}

func (suite *InterviewerControllerTestSuite) TestCreateHandler_BadRequest() {
	invalidPayload := "invalid JSON"

	req := httptest.NewRequest("POST", "/api/v1/interviewers", strings.NewReader(invalidPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "err")
}

func (suite *InterviewerControllerTestSuite) TestCreateHandler_InternalServerError() {
	payload := model.Interviewer{
		InterviewerID: "1",
		FullName:      "John Doe",
		UserID:        "1",
	}

	expectedError := fmt.Errorf("internal server error")
	suite.interviewerUsecaseMock.On("RegisterNewInterviewer", payload).Return(expectedError)

	requestBody, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/interviewers", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *InterviewerControllerTestSuite) TestListHandler_Success() {
	dummyInterviewers := []model.Interviewer{
		{
			InterviewerID: "1",
			FullName:      "John Doe",
			UserID:        "123",
		},
		{
			InterviewerID: "2",
			FullName:      "Jane Smith",
			UserID:        "456",
		},
	}
	suite.interviewerUsecaseMock.On("FindAllInterviewer").Return(dummyInterviewers, nil)

	req := httptest.NewRequest("GET", "/api/v1/interviewers", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *InterviewerControllerTestSuite) TestListHandler_Error() {
	expectedError := fmt.Errorf("failed to retrieve interviewers")
	suite.interviewerUsecaseMock.On("FindAllInterviewer").Return(nil, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/interviewers", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

}

func (suite *InterviewerControllerTestSuite) TestGetHandler_Success() {
	expectedInterviewer := model.Interviewer{
		InterviewerID: "1",
		FullName:      "John Doe",
		UserID:        "123",
	}
	suite.interviewerUsecaseMock.On("FindByIdInterviewer", "1").Return(expectedInterviewer, nil)

	req := httptest.NewRequest("GET", "/api/v1/interviewers/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *InterviewerControllerTestSuite) TestGetHandler_Error() {
	expectedError := fmt.Errorf("interviewer with id 2 not found")
	suite.interviewerUsecaseMock.On("FindByIdInterviewer", "2").Return(model.Interviewer{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/interviewers/2", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
func (suite *InterviewerControllerTestSuite) TestUpdateHandler_Success() {
	payload := model.Interviewer{
		InterviewerID: "1",
		FullName:      "Updated Name",
		UserID:        "123",
	}
	suite.interviewerUsecaseMock.On("UpdateInterviewer", payload).Return(nil)
	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/api/v1/interviewers", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *InterviewerControllerTestSuite) TestUpdateHandler_MissingID() {
	// Set up mock behavior to simulate an error when InterviewerID is missing
	payload := model.Interviewer{
		FullName: "Updated Name",
		UserID:   "123",
	}
	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/api/v1/interviewers", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *InterviewerControllerTestSuite) TestUpdateHandler_Error() {
	payload := model.Interviewer{
		InterviewerID: "1",
		FullName:      "Updated Name",
		UserID:        "123",
	}
	suite.interviewerUsecaseMock.On("UpdateInterviewer", payload).Return(fmt.Errorf("error"))
	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/api/v1/interviewers", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *InterviewerControllerTestSuite) TestUpdateHandler_BadRequest() {
	requestBody := []byte("invalid_json")
	req := httptest.NewRequest("PUT", "/api/v1/interviewers", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *InterviewerControllerTestSuite) TestDeleteHandler_Success() {
	id := "1"
	suite.interviewerUsecaseMock.On("DeleteInterviewer", id).Return(nil)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/interviewers/%s", id), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *InterviewerControllerTestSuite) TestDeleteHandler_InternalServerError() {
	id := "2"
	suite.interviewerUsecaseMock.On("DeleteInterviewer", id).Return(fmt.Errorf("error"))

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/interviewers/%s", id), nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
