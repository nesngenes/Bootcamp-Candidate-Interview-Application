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
	"time"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type bootcampUseCaseMock struct {
	mock.Mock
}

func parseTime(timeStr string) time.Time {
    parsedTime, err := time.Parse(time.RFC3339, timeStr)
    if err != nil {
        panic(err) // Handle the error appropriately in your code
    }
    return parsedTime
}

// DeleteBootcamp implements usecase.BootcampUseCase.
func (b *bootcampUseCaseMock) DeleteBootcamp(id string) error {
	args := b.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// FindAllBootcamp implements usecase.BootcampUseCase.
func (b *bootcampUseCaseMock) FindAllBootcamp(requesPaging dto.PaginationParam) ([]model.Bootcamp, dto.Paging, error) {
	args := b.Called(requesPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Bootcamp), args.Get(1).(dto.Paging), nil

}

// FindByIdBootcamp implements usecase.BootcampUseCase.
func (b *bootcampUseCaseMock) FindByIdBootcamp(id string) (model.Bootcamp, error) {
	args := b.Called(id)
	if args.Get(1) != nil {
		return model.Bootcamp{}, args.Error(1)
	}
	return args.Get(0).(model.Bootcamp), nil

}

// FindByIdBootcamp implements usecase.BootcampUseCase.
func (b *bootcampUseCaseMock) GetBootcampByID(id string) (model.Bootcamp, error) {
	args := b.Called(id)
	if args.Get(1) != nil {
		return model.Bootcamp{}, args.Error(1)
	}
	return args.Get(0).(model.Bootcamp), nil

}

// RegisterNewBootcamp implements usecase.BootcampUseCase.
func (b *bootcampUseCaseMock) RegisterNewBootcamp(payload model.Bootcamp) error {
	args := b.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

// UpdateBootcamp implements usecase.BootcampUseCase.
func (b *bootcampUseCaseMock) UpdateBootcamp(payload model.Bootcamp) error {
	args := b.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

type BootcampControllerTestSuite struct {
	suite.Suite
	usecaseMock *bootcampUseCaseMock
	router      *gin.Engine
}

func (suite *BootcampControllerTestSuite) SetupTest() {
	suite.usecaseMock = new(bootcampUseCaseMock)
	suite.router = gin.Default()
	NewBootcampController(suite.router, suite.usecaseMock)
}

func TestBootcampControllerTestSuite(t *testing.T) {
	suite.Run(t, new(BootcampControllerTestSuite))
}


func (suite *BootcampControllerTestSuite) TestCreateHandler_Success() {
	payload := model.Bootcamp{
		BootcampId: "1",
        Name:       "Bootcamp 1",
        StartDate:  parseTime("2023-01-02T00:00:00Z"),
        EndDate:    parseTime("2023-01-15T00:00:00Z"),
        Location:   "LA",
	}
	
	suite.usecaseMock.On("RegisterNewBootcamp", mock.AnythingOfType("model.Bootcamp")).Return(nil)

	payloadJSON, _ := json.Marshal(payload)

	//HTTP request
	req := httptest.NewRequest("POST", "/api/v1/bootcamps", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	suite.usecaseMock.AssertExpectations(suite.T())
}

func (suite *BootcampControllerTestSuite) TestCreateHandler_BadRequest() {
	invalidPayload := "invalid JSON"

	req := httptest.NewRequest("POST", "/api/v1/bootcamps", strings.NewReader(invalidPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "err")
}

func (suite *BootcampControllerTestSuite) TestCreateHandler_InternalServerError() {
	payload := model.Bootcamp{
		BootcampId: "1",
        Name:       "Bootcamp 1",
        StartDate:  parseTime("2023-01-02T00:00:00Z"),
        EndDate:    parseTime("2023-01-15T00:00:00Z"),
        Location:   "LA",
	}

	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("RegisterNewBootcamp", mock.MatchedBy(func(arg model.Bootcamp) bool {
		return arg.Name == payload.Name
	})).Return(expectedError)

	requestBody, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/bootcamps", strings.NewReader(string(requestBody)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	// Verify the mock method call with any BootcampId value
	suite.usecaseMock.AssertCalled(suite.T(), "RegisterNewBootcamp", mock.MatchedBy(func(arg model.Bootcamp) bool {
		return arg.Name == payload.Name
	}))
}

func (suite *BootcampControllerTestSuite) TestListHandler_Success() {
	// Set up mock behavior
	paginationParam := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}
	expectedResultes := []model.Bootcamp{
		{
			BootcampId: "1",
			Name:       "Bootcamp 1",
			StartDate:  parseTime("2023-01-02T00:00:00Z"),
			EndDate:    parseTime("2023-01-15T00:00:00Z"),
			Location:   "LA",
		},
		{
			BootcampId: "2",
			Name:       "Bootcamp 2",
			StartDate:  parseTime("2023-02-02T00:00:00Z"),
			EndDate:    parseTime("2023-02-15T00:00:00Z"),
			Location:   "SEOL",
		},
	}
	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   len(expectedResultes),
		TotalPages:  1,
	}
	suite.usecaseMock.On("FindAllBootcamp", paginationParam).Return(expectedResultes, expectedPaging, nil)

	req := httptest.NewRequest("GET", "/api/v1/bootcamps?page=1&limit=5", nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BootcampControllerTestSuite) TestListHandler_InternalServerError() {
	paginationParam := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("FindAllBootcamp", paginationParam).Return(nil, dto.Paging{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/bootcamps?page=1&limit=5", nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}


func (suite *BootcampControllerTestSuite) TestGetHandler_Success() {
	bootcampID := "1"
	expectedResult := model.Bootcamp{BootcampId: bootcampID, Name: "Bootcamp 1", StartDate: parseTime("2023-01-02T00:00:00Z"), EndDate: parseTime("2023-01-15T00:00:00Z")}
	suite.usecaseMock.On("FindByIdBootcamp", bootcampID).Return(expectedResult, nil)

	req := httptest.NewRequest("GET", "/api/v1/bootcamps/"+bootcampID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BootcampControllerTestSuite) TestGetHandler_InternalServerError() {
	bootcampID := "1"
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("FindByIdBootcamp", bootcampID).Return(model.Bootcamp{}, expectedError)

	req := httptest.NewRequest("GET", "/api/v1/bootcamps/"+bootcampID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *BootcampControllerTestSuite) TestUpdateHandler_Success() {
	expectedResult := model.Bootcamp{BootcampId: "1", Name: "Bootcamp 1", StartDate: parseTime("2023-01-02T00:00:00Z"), EndDate: parseTime("2023-01-15T00:00:00Z"), Location: "LA"}
	suite.usecaseMock.On("UpdateBootcamp", expectedResult).Return(nil)

	requestBody := []byte(`{"id": "1", "name": "Bootcamp 1", "start_date": "2023-01-02T00:00:00Z", "end_date": "2023-01-15T00:00:00Z", "location": "LA"}`)
	req := httptest.NewRequest("PUT", "/api/v1/bootcamps", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *BootcampControllerTestSuite) TestUpdateHandler_BadRequest() {
	req := httptest.NewRequest("PUT", "/api/v1/bootcamps", strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *BootcampControllerTestSuite) TestUpdateHandler_InternalServerError() {
	expectedResult := model.Bootcamp{BootcampId: "1", Name: "Bootcamp 1", StartDate: parseTime("2023-01-02T00:00:00Z"), EndDate: parseTime("2023-01-15T00:00:00Z"), Location: "LA"}
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("UpdateBootcamp", expectedResult).Return(expectedError)

	requestBody := []byte(`{"id": "1", "name": "Bootcamp 1", "start_date": "2023-01-02T00:00:00Z", "end_date": "2023-01-15T00:00:00Z", "location": "LA"}`)
	req := httptest.NewRequest("PUT", "/api/v1/bootcamps", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}

func (suite *BootcampControllerTestSuite) TestDeleteHandler_Success() {
	bootcampID := "1"
	suite.usecaseMock.On("DeleteBootcamp", bootcampID).Return(nil)

	req := httptest.NewRequest("DELETE", "/api/v1/bootcamps/"+bootcampID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
}

func (suite *BootcampControllerTestSuite) TestDeleteHandler_InternalServerError() {
	bootcampID := "1"
	expectedError := fmt.Errorf("internal server error")
	suite.usecaseMock.On("DeleteBootcamp", bootcampID).Return(expectedError)

	req := httptest.NewRequest("DELETE", "/api/v1/bootcamps/"+bootcampID, nil)

	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}