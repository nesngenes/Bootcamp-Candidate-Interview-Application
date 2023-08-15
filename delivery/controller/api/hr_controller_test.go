package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"interview_bootcamp/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type hrUseCaseMock struct {
	mock.Mock
}

func (h *hrUseCaseMock) CreateHRRecruitment(payload model.HRRecruitment) error {
	args := h.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
func (h *hrUseCaseMock) Get(id string) (model.HRRecruitment, error) {
	args := h.Called(id)
	if args.Get(1) != nil {
		return model.HRRecruitment{}, args.Error(1)
	}
	return args.Get(0).(model.HRRecruitment), nil

}

func (h *hrUseCaseMock) ListHRRecruitments() ([]model.HRRecruitment, error) {
	args := h.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.HRRecruitment), args.Error(1)
}

func (h *hrUseCaseMock) UpdateHRRecruitment(payload model.HRRecruitment) error {
	args := h.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (h *hrUseCaseMock) DeleteHRRecruitment(id string) error {
	args := h.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

type HRRecruitmentControllerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	useCaseMock *hrUseCaseMock
}

func (suite *HRRecruitmentControllerTestSuite) SetupTest() {
	suite.router = gin.Default()
	suite.useCaseMock = new(hrUseCaseMock)
	NewHRRecruitmentController(suite.router, suite.useCaseMock)
}

func TestHRControllerTestSuite(t *testing.T) {
	suite.Run(t, new(HRRecruitmentControllerTestSuite))
}

func (suite *HRRecruitmentControllerTestSuite) TestCreateHandler_BadRequest() {
	// request dengan invalid payload
	invalidPayload := "invalid JSON"
	req := httptest.NewRequest("POST", "/api/v1/hr-recruitment", bytes.NewReader([]byte(invalidPayload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

}

func (suite *HRRecruitmentControllerTestSuite) TestCreateHandler_Success() {
	expectedUser := model.HRRecruitment{
		ID:       "user123",
		FullName: "Athy",
		UserID:   "user123",
		User: model.Users{
			Id:       "user123",
			Email:    "athy@example.com",
			UserName: "athy",
			Password: "password",
			UserRole: model.UserRoles{
				Id:   "role123",
				Name: "HR",
			},
		},
	}

	// Set up the mock expectation for Get (user not found) and CreateHRRecruitment
	suite.useCaseMock.On("Get", "user123").Return(model.HRRecruitment{}, errors.New("not found"))

	suite.useCaseMock.On("CreateHRRecruitment", mock.AnythingOfType("model.HRRecruitment")).Return(nil)

	payloadJSON, _ := json.Marshal(expectedUser)

	//HTTP request
	req := httptest.NewRequest("POST", "/api/v1/hr-recruitment", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	suite.useCaseMock.AssertExpectations(suite.T())
}

func (suite *HRRecruitmentControllerTestSuite) TestListHandler_InternalServerError() {
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("ListHRRecruitments").Return(nil, expectedError)

	// request
	req := httptest.NewRequest("GET", "/api/v1/hr-recruitment", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

}

func (suite *HRRecruitmentControllerTestSuite) TestListHandler_Success() {
	// Setup
	expectedUser := []model.HRRecruitment{
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
	suite.useCaseMock.On("ListHRRecruitments").Return(expectedUser, nil)

	req := httptest.NewRequest("GET", "/api/v1/hr-recruitment", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *HRRecruitmentControllerTestSuite) TestGetHandler_Success() {
	// Setup
	expectedUser := model.HRRecruitment{
		ID:       "1",
		FullName: "Ella Ella",
		UserID:   "user123",
		User: model.Users{
			Id:       "user123",
			Email:    "ella@mail.com",
			UserName: "ellaa",
			UserRole: model.UserRoles{
				Id:   "1",
				Name: "HR",
			},
		},
	}

	suite.useCaseMock.On("Get", "1").Return(expectedUser, nil)

	req := httptest.NewRequest("GET", "/api/v1/hr-recruitment/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

}

func (suite *HRRecruitmentControllerTestSuite) TestGetHandler_InternalServerError() {
	// Setup
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("Get", "1").Return(model.HRRecruitment{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/hr-recruitment/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *HRRecruitmentControllerTestSuite) TestUpdateHandler_Success() {
	expectedHRRecruitment := model.HRRecruitment{
		ID:       "user123",
		FullName: "Athanasia",
		UserID:   "user123",
		User: model.Users{
			Id:       "user123",
			Email:    "athanasia@example.com",
			UserName: "athanasia",
			Password: "password",
			UserRole: model.UserRoles{
				Id:   "role123",
				Name: "HR",
			},
		},
	}

	suite.useCaseMock.On("Get", "user123").Return(expectedHRRecruitment, nil)
	suite.useCaseMock.On("UpdateHRRecruitment", expectedHRRecruitment).Return(nil)

	payloadJSON, _ := json.Marshal(expectedHRRecruitment)
	req := httptest.NewRequest("PUT", "/api/v1/hr-recruitment/"+expectedHRRecruitment.ID, bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	suite.useCaseMock.AssertCalled(suite.T(), "Get", expectedHRRecruitment.ID)
	suite.useCaseMock.AssertCalled(suite.T(), "UpdateHRRecruitment", expectedHRRecruitment)
}

func (suite *HRRecruitmentControllerTestSuite) TestUpdateHR_NotFound() {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, "/api/v1/hr-recruiment/1", nil)
	suite.router.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusNotFound, recorder.Code)
}

func (suite *HRRecruitmentControllerTestSuite) TestUpdateHR_InternalServerError() {
	// Setup the mock expectation for Get method call
	expectedHRRecruitment := model.HRRecruitment{
		ID:       "user123",
		FullName: "Athanasia",
		UserID:   "user123",
		User: model.Users{
			Id:       "user123",
			Email:    "athanasia@example.com",
			UserName: "athanasia",
			Password: "password",
			UserRole: model.UserRoles{
				Id:   "role123",
				Name: "HR",
			},
		},
	}
	suite.useCaseMock.On("Get", "user123").Return(expectedHRRecruitment, nil)

	// Setup the mock expectation for UpdateHRRecruitment method call
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("UpdateHRRecruitment", expectedHRRecruitment).Return(expectedError)

	// Create the request payload
	payloadJSON, _ := json.Marshal(expectedHRRecruitment)

	// Make the HTTP request
	req := httptest.NewRequest("PUT", "/api/v1/hr-recruitment/"+expectedHRRecruitment.ID, bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert the response status code
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	// Assert that the mock methods were called
	suite.useCaseMock.AssertCalled(suite.T(), "Get", expectedHRRecruitment.ID)
	suite.useCaseMock.AssertCalled(suite.T(), "UpdateHRRecruitment", expectedHRRecruitment)
}

func (suite *HRRecruitmentControllerTestSuite) TestDeleteHandler_Success() {
	// Setup
	suite.useCaseMock.On("DeleteHRRecruitment", "1").Return(nil)

	//request
	req := httptest.NewRequest("DELETE", "/api/v1/hr-recruitment/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *HRRecruitmentControllerTestSuite) TestDeleteHandler_InternalServerError() {
	// Setup
	expectedError := errors.New("internal server error")
	suite.useCaseMock.On("DeleteHRRecruitment", "1").Return(expectedError)

	//request
	req := httptest.NewRequest("DELETE", "/api/v1/hr-recruitment/1", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
