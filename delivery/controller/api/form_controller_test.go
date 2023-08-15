package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"mime/multipart" // Import the multipart package
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/utils/common"

	"github.com/cloudinary/cloudinary-go/v2"
)

type formUseCaseMock struct {
	mock.Mock
}

func (f *formUseCaseMock) RegisterNewForm(payload model.Form) error {
	args := f.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (f *formUseCaseMock) FindAllForm(requestPaging dto.PaginationParam) ([]model.Form, dto.Paging, error) {
	args := f.Called(requestPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Form), args.Get(1).(dto.Paging), nil
}

func (f *formUseCaseMock) UpdateForm(payload model.Form) error {
	args := f.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (f *formUseCaseMock) DeleteForm(id string) error {
	args := f.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}


func (f *formUseCaseMock) FindByIdForm(id string) (model.Form, error) {
	args := f.Called(id)
	if args.Get(1) != nil {
		return model.Form{}, args.Error(1)
	}
	return args.Get(0).(model.Form), nil
}

type FormControllerTestSuite struct {
	suite.Suite
	usecaseMock   *formUseCaseMock // Change the type here
	router        *gin.Engine
	cloudinaryAPI *cloudinary.Cloudinary
}

func (suite *FormControllerTestSuite) SetupTest() {
	suite.usecaseMock = new(formUseCaseMock)
	suite.cloudinaryAPI, _ = cloudinary.NewFromParams("cloud_name", "api_key", "api_secret")
	suite.router = gin.Default()
}

func TestFormControllerTestSuite(t *testing.T) {
	suite.Run(t, new(FormControllerTestSuite))
}

func (suite *FormControllerTestSuite) TestListHandlerForm_Success() {
	expectedPaginationParam := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	expectedForms := []model.Form{
		{
			FormID:   "1",
			FormLink: "https://example.com/form_link",
			// ... fill in other form fields ...
		},
		// ... add more dummy forms ...
	}

	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   len(expectedForms),
		TotalPages:  1,
	}

	suite.usecaseMock.On("FindAllForm", expectedPaginationParam).Return(expectedForms, expectedPaging, nil)

	NewFormController(suite.router, suite.usecaseMock, suite.cloudinaryAPI)

	request, _ := http.NewRequest(http.MethodGet, "/api/v1/forms?page=1&limit=5", nil)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = request
	suite.router.ServeHTTP(context.Writer, context.Request)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)

	// Unmarshal the response body into a structure and make assertions
	var responseBody struct {
		Status common.WebStatus `json:"status"`
		Data   []model.Form     `json:"data"`
		Paging dto.Paging       `json:"paging"`
	}
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedForms, responseBody.Data)
	assert.Equal(suite.T(), expectedPaging, responseBody.Paging)
}

func (suite *FormControllerTestSuite) TestUpdateHandlerForm_Success() {
	// Arrange: Set up the expected form data and mock behavior
	formID := "1ABC"
	expectedForm := model.Form{
		FormID:   formID,
		FormLink: "https://example.com/form_link_updated",
		// ... set other form fields ...
	}
	suite.usecaseMock.On("FindByIdForm", formID).Return(expectedForm, nil)
	suite.usecaseMock.On("UpdateForm", expectedForm).Return(nil)

	// Create a new multipart writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add form fields to the writer
	_ = writer.WriteField("id", formID)
	_ = writer.WriteField("form_link", expectedForm.FormLink)
	// ... add other form fields ...

	writer.Close()

	// Create a new HTTP request with the multipart form data
	request, _ := http.NewRequest(http.MethodPut, "/api/v1/forms", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Act: Send the request to the router
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)
	context.Request = request

	// Create a new FormController with the mock usecase and cloudinary
	controller := NewFormController(suite.router, suite.usecaseMock, suite.cloudinaryAPI)
	controller.updateHandler(context)

	// Assert: Verify the response and mock expectations
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response struct {
		Status common.WebStatus `json:"status"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, response.Status.Code)
	assert.Equal(suite.T(), "Update data successfully", response.Status.Description)

	suite.usecaseMock.AssertExpectations(suite.T())
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}
