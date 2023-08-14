package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockBootcampUseCase struct {
	mock.Mock
}

// DeleteBootcamp implements usecase.BootcampUseCase.
func (*mockBootcampUseCase) DeleteBootcamp(id string) error {
	panic("unimplemented")
}

// FindAllBootcamp implements usecase.BootcampUseCase.
func (*mockBootcampUseCase) FindAllBootcamp(requesPaging dto.PaginationParam) ([]model.Bootcamp, dto.Paging, error) {
	panic("unimplemented")
}

// FindByIdBootcamp implements usecase.BootcampUseCase.
func (*mockBootcampUseCase) FindByIdBootcamp(id string) (model.Bootcamp, error) {
	panic("unimplemented")
}

// RegisterNewBootcamp implements usecase.BootcampUseCase.
func (*mockBootcampUseCase) RegisterNewBootcamp(payload model.Bootcamp) error {
	panic("unimplemented")
}

// UpdateBootcamp implements usecase.BootcampUseCase.
func (*mockBootcampUseCase) UpdateBootcamp(payload model.Bootcamp) error {
	panic("unimplemented")
}

func (m *mockBootcampUseCase) GetBootcampByID(id string) (model.Bootcamp, error) {
	args := m.Called(id)
	return args.Get(0).(model.Bootcamp), args.Error(1)
}

type mockCandidateUseCase struct {
	mock.Mock
}

// DeleteCandidate implements usecase.CandidateUseCase.
func (c *mockCandidateUseCase) DeleteCandidate(id string) error {
	args := c.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

// FindAllCandidate implements usecase.CandidateUseCase.
func (c *mockCandidateUseCase) FindAllCandidate(requestPaging dto.PaginationParam) ([]model.Candidate, dto.Paging, error) {
	args := c.Called(requestPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Candidate), args.Get(1).(dto.Paging), nil

}

// FindByIdCandidate implements usecase.CandidateUseCase.
func (c *mockCandidateUseCase) FindByIdCandidate(id string) (model.Candidate, error) {
	args := c.Called(id)
	if args.Get(1) != nil {
		return model.Candidate{}, args.Error(1)
	}
	return args.Get(0).(model.Candidate), nil

}

// UpdateCandidate implements usecase.CandidateUseCase.
func (*mockCandidateUseCase) UpdateCandidate(payload model.Candidate) error {
	panic("unimplemented")
}

func (m *mockCandidateUseCase) RegisterNewCandidate(candidate model.Candidate) error {
	args := m.Called(candidate)
	return args.Error(0)
}

type CandidateControllerSuite struct {
	suite.Suite
	mockCandidateUsecase *mockCandidateUseCase
	mockBootcampUsecase  *mockBootcampUseCase
	mockCloudinary       *cloudinary.Cloudinary
	controller           *CandidateController
	router               *gin.Engine
}

func (suite *CandidateControllerSuite) SetupTest() {
	suite.mockCandidateUsecase = new(mockCandidateUseCase)
	suite.mockBootcampUsecase = new(mockBootcampUseCase)
	suite.mockCloudinary = &cloudinary.Cloudinary{} // Mock the Cloudinary object

	suite.router = gin.Default()
	suite.controller = NewCandidateController(suite.router, suite.mockCandidateUsecase, suite.mockBootcampUsecase, suite.mockCloudinary)
}

func TestCandidateControllerSuite(t *testing.T) {
	suite.Run(t, new(CandidateControllerSuite))
}

func (suite *CandidateControllerSuite) TestCreateHandler_Failure_UploadFile() {
	// Prepare mock input data
	candidate := model.Candidate{
		CandidateID: "1",
		FullName:    "John Doe",
		Phone:       "1111",
		Email:       "john@example.com",
		DateOfBirth: "2005-05-01",
		Address:     "st12",
		CvLink:      "http://example.com",
		Bootcamp: model.Bootcamp{
			BootcampId: "1",
		},
		InstansiPendidikan: "smk",
		HackerRank:         90,
	}
	suite.mockCandidateUsecase.On("RegisterNewCandidate", candidate).Return(nil)

	// Create a test request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/candidates", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	assert.Equal(suite.T(), 400, rec.Code)
}

func (suite *CandidateControllerSuite) TestListHandler_Success() {
	expectedPaginationParam := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	candidates := []model.Candidate{
		{
			CandidateID:        "1",
			FullName:           "ali",
			Phone:              "1111",
			Email:              "a@example.com",
			DateOfBirth:        "2020-01-01",
			Address:            "st12",
			CvLink:             "https://www.example.com",
			Bootcamp:           model.Bootcamp{},
			InstansiPendidikan: "smk",
			HackerRank:         90,
		},
	}
	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   1,
		TotalPages:  1,
	}

	suite.mockCandidateUsecase.On("FindAllCandidate", expectedPaginationParam).Return(candidates, expectedPaging, nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/v1/candidates?page=1&limit=5", nil)
	suite.router.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *CandidateControllerSuite) TestListHandler_Fail() {
	suite.mockCandidateUsecase.On("FindAllCandidate", dto.PaginationParam{Page: 1, Limit: 5}).Return(nil, dto.Paging{}, errors.New("error"))
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/v1/candidates?page=1&limit=5", nil)
	suite.router.ServeHTTP(recorder, request)
	response := recorder.Body.Bytes()
	var actualError struct {
		Err string
	}
	json.Unmarshal(response, &actualError)
	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
	assert.Equal(suite.T(), "error", actualError.Err)
}

func (suite *CandidateControllerSuite) TestGetHandler_Success() {
	// Prepare mock input data
	candidateID := "1"
	candidate := model.Candidate{
		CandidateID:        "1",
		FullName:           "ali",
		Phone:              "1111",
		Email:              "a@example.com",
		DateOfBirth:        "2020-01-01",
		Address:            "st12",
		CvLink:             "https://www.example.com",
		Bootcamp:           model.Bootcamp{},
		InstansiPendidikan: "smk",
		HackerRank:         90,
	}
	suite.mockCandidateUsecase.On("FindByIdCandidate", candidateID).Return(candidate, nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/candidates/"+candidateID, nil)
	suite.router.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

}

func (suite *CandidateControllerSuite) TestGetHandler_Failure() {
	// Prepare mock input data
	candidateID := "1"
	errorMsg := "Candidate not found"
	suite.mockCandidateUsecase.On("FindByIdCandidate", candidateID).Return(model.Candidate{}, errors.New(errorMsg))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/v1/candidates/"+candidateID, nil)
	suite.router.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
}

func (suite *CandidateControllerSuite) TestDeleteHandler_Success() {
	candidateID := "1"
	suite.mockCandidateUsecase.On("DeleteCandidate", candidateID).Return(nil)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/v1/candidates/"+candidateID, nil)
	suite.router.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusNoContent, recorder.Code)
}

func (suite *CandidateControllerSuite) TestDeleteHandler_Failure() {
	candidateID := "1"
	errorMsg := "Failed to delete candidate"
	suite.mockCandidateUsecase.On("DeleteCandidate", candidateID).Return(errors.New(errorMsg))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("DELETE", "/api/v1/candidates/"+candidateID, nil)
	suite.router.ServeHTTP(recorder, request)
	assert.Equal(suite.T(), http.StatusInternalServerError, recorder.Code)
}

func (suite *CandidateControllerSuite) TestUpdateHandler_Success() {
	// Prepare mock input data
	candidateID := "1"
	existingCandidate := model.Candidate{
		CandidateID: candidateID,
		FullName:    "dffgdfg",
		Phone:       "234254",
		Email:       "sdfsdf",
		DateOfBirth: "sdfds",
		Address:     "sdf",
		CvLink:      "sdf",
		Bootcamp: model.Bootcamp{
			BootcampId: "1",
		},
		InstansiPendidikan: "sdffd",
		HackerRank:         90,
	}
	bootcampID := "1"
	bootcamp := model.Bootcamp{
		BootcampId: bootcampID,
		Name:       "sad",
		StartDate:  time.Time{},
		EndDate:    time.Time{},
		Location:   "asds",
	}

	suite.mockCandidateUsecase.On("FindByIdCandidate", candidateID).Return(existingCandidate, nil)
	suite.mockCandidateUsecase.On("UpdateCandidate", existingCandidate).Return(nil)
	suite.mockBootcampUsecase.On("GetBootcampByID", bootcampID).Return(bootcamp, nil)

	// Create a test request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// Add form fields and file here
	writer.Close()

	req := httptest.NewRequest("PUT", "/api/v1/candidates/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	suite.router.ServeHTTP(rec, req)

	assert.Equal(suite.T(), 307, rec.Code)

}
