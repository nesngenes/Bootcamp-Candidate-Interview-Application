package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"interview_bootcamp/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type userUsecaseMock struct {
	mock.Mock
}

func (m *userUsecaseMock) RegisterNewUser(payload model.Users) error {
	args := m.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *userUsecaseMock) List() ([]model.Users, error) {
	args := m.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Users), args.Error(1)
}

func (m *userUsecaseMock) GetUserByID(id string) (model.Users, error) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return model.Users{}, args.Error(1)
	}
	return args.Get(0).(model.Users), args.Error(1)
}
func (m *userUsecaseMock) GetUserByEmail(email string) (model.Users, error) {
	args := m.Called(email)
	if args.Get(1) != nil {
		return model.Users{}, args.Error(1)
	}
	return args.Get(0).(model.Users), args.Error(1)
}
func (m *userUsecaseMock) GetUserByUserName(username string) (model.Users, error) {
	args := m.Called(username)
	if args.Get(1) != nil {
		return model.Users{}, args.Error(1)
	}
	return args.Get(0).(model.Users), args.Error(1)
}

func (m *userUsecaseMock) FindByUsernamePassword(username string, password string) (model.Users, error) {
	args := m.Called(username, password)
	return args.Get(0).(model.Users), args.Error(1)
}

func (m *userUsecaseMock) UpdateUser(payload model.Users) error {
	args := m.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

func (m *userUsecaseMock) DeleteUser(id string) error {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

type UserControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	useCaseMock *userUsecaseMock
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.router = gin.Default()
	suite.useCaseMock = new(userUsecaseMock)
	NewUserController(suite.router, suite.useCaseMock)
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (suite *UserControllerTestSuite) TestCreateHandler_Success() {
	expectedUser := model.Users{
		Id:       "1",
		Email:    "ella@mail.com",
		UserName: "ella-updated",
		UserRole: model.UserRoles{
			Id:   "3",
			Name: "CEO",
		},
	}

	suite.useCaseMock.On("RegisterNewUser", mock.MatchedBy(func(arg model.Users) bool {
		return arg.UserName == expectedUser.UserName
	})).Return(nil) //

	requestBody, _ := json.Marshal(expectedUser)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// menangkap respon
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	//buat pastiin expectednya emg payload yg bener
	suite.useCaseMock.AssertCalled(suite.T(), "RegisterNewUser", mock.MatchedBy(func(arg model.Users) bool {
		return arg.UserName == expectedUser.UserName
	}))
}

func (suite *UserControllerTestSuite) TestCreateHandler_BadRequest() {
	// request dengan invalid payload
	invalidPayload := "invalid JSON"
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader([]byte(invalidPayload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "err")
}

func (suite *UserControllerTestSuite) TestCreateHandler_InternalServerError() {
	payload := model.Users{
		Email:    "ella@mail.com",
	}

	expectedError := fmt.Errorf("internal server error")
	suite.useCaseMock.On("RegisterNewUser", mock.MatchedBy(func(arg model.Users) bool {
		return arg.Email == payload.Email
	})).Return(expectedError)

	requestBody, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	// Verify the mock method call with any ResultId value
	suite.useCaseMock.AssertCalled(suite.T(), "RegisterNewUser", mock.MatchedBy(func(arg model.Users) bool {
		return arg.Email == payload.Email
	}))
}


func (suite *UserControllerTestSuite) TestListHandler_Success() {
	// Setup
	expectedUser := []model.Users{
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
	suite.useCaseMock.On("List").Return(expectedUser, nil)

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *UserControllerTestSuite) TestListHandler_InternalServerError() {
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("List").Return(nil, expectedError)

	// request
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

}

func (suite *UserControllerTestSuite) TestGetHandler_Success() {
	// Setup
	expectedUser := model.Users{
		Id:       "1",
		Email:    "ella@mail.com",
		UserName: "ella-updated",
		UserRole: model.UserRoles{
			Id:   "3",
			Name: "CEO",
		},
	}
	suite.useCaseMock.On("GetUserByID", "1").Return(expectedUser, nil)

	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

}

// get id 500
func (suite *UserControllerTestSuite) TestGetHandler_InternalServerError() {
	// Setup
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("GetUserByID", "1").Return(model.Users{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

// func (suite *UserControllerTestSuite) TestGetUserByUserNameHandler_InternalServerError() {
// 	// Setup
// 	expectedError := errors.New("internal server error")
// 	suite.useCaseMock.On("GetUserByUserName", "lia").Return(model.Users{}, expectedError)

// 	req := httptest.NewRequest("GET", "/api/v1/users/by-username/lia", nil)
// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
// }

func (suite *UserControllerTestSuite) TestUpdateHandler_Success() {
	// Setup
	// Persiapan: Menyiapkan data payload untuk pengujian
	payload := model.Users{
		Id:       "1",
		Email:    "ella@mail.com",
		UserName: "ella-updated",
		UserRole: model.UserRoles{
			Id:   "3",
			Name: "CEO",
		},
	}
	// Menyiapkan ekspektasi pemanggilan metode UpdateUser dengan payload yang diberikan
	suite.useCaseMock.On("UpdateUser", payload).Return(nil)

	// Request
	// Membuat permintaan HTTP PUT dengan payload sebagai isi permintaan
	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/api/v1/users/1", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json") // dengan body permintaan yang diambil dari requestBody (format JSON).
	w := httptest.NewRecorder()
	// Menjalankan permintaan HTTP menggunakan router Gin
	suite.router.ServeHTTP(w, req) // Respons dari permintaan akan direkam oleh w.

	// Assertion
	// Memastikan kode status HTTP yang diterima sesuai dengan yang diharapkan (http.StatusOK)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// mau   aku bikin buat bad request
func (suite *UserControllerTestSuite) TestUpdateUser_FailBind() {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, "/api/v1/users/1", nil)
	suite.router.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *UserControllerTestSuite) TestUpdateUser_InternalServerError() {
    expectedResult := model.Users{Id: "3", Email: "ella@mail.com"}
    expectedError := fmt.Errorf("internal server error")
    
    // Change "UpdateUserRole" to "UpdateUser" in the following line
    suite.useCaseMock.On("UpdateUser", expectedResult).Return(expectedError)

    requestBody := []byte(`{"id": "3", "email": "ella@mail.com"}`)
    req := httptest.NewRequest("PUT", "/api/v1/users/3", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *UserControllerTestSuite) TestDeleteHandler_Success() {
	// Setup
	suite.useCaseMock.On("DeleteUser", "1").Return(nil)

	//request
	req := httptest.NewRequest("DELETE", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *UserControllerTestSuite) TestDeleteHandler_InternalServerError() {
	// Setup
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("DeleteUser", "1").Return(expectedError)

	//request
	req := httptest.NewRequest("DELETE", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
