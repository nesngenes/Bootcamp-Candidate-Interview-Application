package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	// "interview_bootcamp/model/dto"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"interview_bootcamp/model"
)

type userRolesUseCaseMock struct {
	mock.Mock
}

func (m *userRolesUseCaseMock) RegisterNewUserRole(payload model.UserRoles) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *userRolesUseCaseMock) GetAllUserRoles() ([]model.UserRoles, error) {
	args := m.Called()
	return args.Get(0).([]model.UserRoles), args.Error(1)
}

func (m *userRolesUseCaseMock) GetUserRoleByID(id string) (model.UserRoles, error) {
	args := m.Called(id)
	return args.Get(0).(model.UserRoles), args.Error(1)
}

func (m *userRolesUseCaseMock) UpdateUserRole(payload model.UserRoles) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *userRolesUseCaseMock) DeleteUserRole(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type UserRoleControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	useCaseMock *userRolesUseCaseMock
}

func (suite *UserRoleControllerTestSuite) SetupTest() {
	suite.router = gin.Default()
	suite.useCaseMock = new(userRolesUseCaseMock)
	NewUserRoleController(suite.router, suite.useCaseMock)
}

func TestUserRoleControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserRoleControllerTestSuite))
}

func (suite *UserRoleControllerTestSuite) TestCreateHandler_Success() {
	expectedUserRole := model.UserRoles{
		Name: "New Role",
	}
	suite.useCaseMock.On("RegisterNewUserRole", mock.MatchedBy(func(arg model.UserRoles) bool {
		return arg.Name == expectedUserRole.Name
	})).Return(nil) //

	requestBody, _ := json.Marshal(expectedUserRole)

	req := httptest.NewRequest("POST", "/api/v1/user-roles", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// menangkap respon
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	//buat pastiin expectednya emg payload yg bener
	suite.useCaseMock.AssertCalled(suite.T(), "RegisterNewUserRole", mock.MatchedBy(func(arg model.UserRoles) bool {
		return arg.Name == expectedUserRole.Name
	}))
}

func (suite *UserRoleControllerTestSuite) TestCreateHandler_BadRequest() {
	// request dengan invalid payload
	invalidPayload := "invalid JSON"
	req := httptest.NewRequest("POST", "/api/v1/user-roles", bytes.NewReader([]byte(invalidPayload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "err")
}

func (suite *UserRoleControllerTestSuite) TestCreateHandler_InternalServerError() {
	payload := model.UserRoles{
		Name: "New UserRoles",
	}

	expectedError := fmt.Errorf("internal server error")
	suite.useCaseMock.On("RegisterNewUserRole", mock.MatchedBy(func(arg model.UserRoles) bool {
		return arg.Name == payload.Name
	})).Return(expectedError)

	requestBody, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/user-roles", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	// Verify the mock method call with any UserRoleId value
	suite.useCaseMock.AssertCalled(suite.T(), "RegisterNewUserRole", mock.MatchedBy(func(arg model.UserRoles) bool {
		return arg.Name == payload.Name
	}))
}
func (suite *UserRoleControllerTestSuite) TestListHandler_Success() {
	// Setup
	expectedUserRoles := []model.UserRoles{
		{Id: "1", Name: "Role 1"},
		{Id: "2", Name: "Role 2"},
	}
	suite.useCaseMock.On("GetAllUserRoles").Return(expectedUserRoles, nil)

	req := httptest.NewRequest("GET", "/api/v1/user-roles", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserRoleControllerTestSuite) TestListHandler_InternalServerError() {
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("GetAllUserRoles").Return([]model.UserRoles{}, expectedError)

	// request
	req := httptest.NewRequest("GET", "/api/v1/user-roles", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *UserRoleControllerTestSuite) TestGetHandler_Success() {
	// Setup
	expectedUserRole := model.UserRoles{Id: "1", Name: "Role 1"}
	suite.useCaseMock.On("GetUserRoleByID", "1").Return(expectedUserRole, nil)

	req := httptest.NewRequest("GET", "/api/v1/user-roles/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

}

func (suite *UserRoleControllerTestSuite) TestGetHandler_InternalServerError() {
	// Setup
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("GetUserRoleByID", "1").Return(model.UserRoles{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/user-roles/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *UserRoleControllerTestSuite) TestUpdateHandler_Success() {
	// Setup
	payload := model.UserRoles{
		Id:   "1",
		Name: "Updated Role",
	}
	suite.useCaseMock.On("UpdateUserRole", payload).Return(nil)

	// request
	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/api/v1/user-roles", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserRoleControllerTestSuite) TestUpdateHandler_BadRequest() {
	// Make request with invalid JSON
	req := httptest.NewRequest("PUT", "/api/v1/user-roles", strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *UserRoleControllerTestSuite) TestUpdateHandler_InternalServerError() {
	// Setup
	expectedError := errors.New("internal server error")
	payload := model.UserRoles{
		Id:   "1",
		Name: "Updated Role",
	}
	suite.useCaseMock.On("UpdateUserRole", payload).Return(expectedError)

	//request
	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/api/v1/user-roles", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *UserRoleControllerTestSuite) TestDeleteHandler_Success() {
	// Setup
	suite.useCaseMock.On("DeleteUserRole", "1").Return(nil)

	//request
	req := httptest.NewRequest("DELETE", "/api/v1/user-roles/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *UserRoleControllerTestSuite) TestDeleteHandler_InternalServerError() {
	// Setup
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("DeleteUserRole", "1").Return(expectedError)

	//request
	req := httptest.NewRequest("DELETE", "/api/v1/user-roles/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
