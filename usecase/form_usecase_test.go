package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	// "interview_bootcamp/repository"
	"interview_bootcamp/model/dto"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

)

type formRepoMock struct {
	mock.Mock
}

func (r *formRepoMock) Create(payload model.Form) error {
	args := r.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *formRepoMock) List() ([]model.Form, error) {
	args := r.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Form), nil
}

func (r *formRepoMock) Get(id string) (model.Form, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.Form{}, args.Error(1)
	}
	return args.Get(0).(model.Form), nil
}
func (r *formRepoMock) GetByPhoneNumber(phoneNumber string) (model.Form, error) {
	args := r.Called(phoneNumber)
	if args.Get(1) != nil {
		return model.Form{}, args.Error(1)
	}
	return args.Get(0).(model.Form), nil
}
func (r *formRepoMock) GetByEmail(email string) (model.Form, error) {
	args := r.Called(email)
	if args.Get(1) != nil {
		return model.Form{}, args.Error(1)
	}
	return args.Get(0).(model.Form), nil
}

func (r *formRepoMock) Update(payload model.Form) error {
	args := r.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *formRepoMock) Delete(id string) error {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *formRepoMock) Paging(requestPaging dto.PaginationParam) ([]model.Form, dto.Paging, error) {
	args := r.Called(requestPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Form), args.Get(1).(dto.Paging), nil
}

type FormUseCaseTestSuite struct {
	suite.Suite
	repoMock *formRepoMock
	cloudinaryMock *CloudinaryMock
	usecase FormUseCase
}

func (suite *FormUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(formRepoMock)
	suite.cloudinaryMock = NewCloudinaryMock()
	suite.usecase = NewFormUseCase(suite.repoMock, &suite.cloudinaryMock.Cloudinary)
}

func TestFormUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(FormUseCaseTestSuite))
}

var formDummy = []model.Form{
	{
		FormID: "1",
		FormLink: "test_link_1.com",
	},
	{
		FormID: "2",
		FormLink: "test_link_2.com",
	},
}

func (suite *FormUseCaseTestSuite) TestRegisterNewForm_Success() {
	dmForm := formDummy[0]
	suite.repoMock.On("Create", dmForm).Return(nil)
	err := suite.usecase.RegisterNewForm(dmForm)
	assert.Nil(suite.T(), err)
}

func (suite *FormUseCaseTestSuite) TestRegisterNewForm_EmptyFields() {
	emptyForm := model.Form{} // Create an empty form

	// Set up an expectation for the Create method with the expected argument
	suite.repoMock.On("Create", emptyForm).Return(nil)

	// Make the usecase call
	err := suite.usecase.RegisterNewForm(emptyForm)

	// Assert the error is nil (no error expected)
	assert.Nil(suite.T(), err)

	// Verify that the expected method call was made
	suite.repoMock.AssertExpectations(suite.T())
}

func (suite *FormUseCaseTestSuite) TestRegisterNewForm_Fail() {
	failedCreateForm := formDummy[0]

	suite.repoMock.On("Create", failedCreateForm).Return(fmt.Errorf("error"))

	err := suite.usecase.RegisterNewForm(failedCreateForm)

	expectedError := fmt.Sprintf("failed to create new form: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)

}

func (suite *FormUseCaseTestSuite) TestFindAllForm_Success() {
	dummy := formDummy
	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   5,
		TotalPages:  1,
	}
	requestPaging := dto.PaginationParam{
		Page: 1,
	}

	suite.repoMock.On("Paging", requestPaging).Return(dummy, expectedPaging, nil)

	forms, paging, err := suite.usecase.FindAllForm(requestPaging)

	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), dummy, forms)
	assert.Equal(suite.T(), expectedPaging, paging)

}

func (suite *FormUseCaseTestSuite) TestFindAllForm_Fail() {
	requestPaging := dto.PaginationParam{
		Page: 1,
	}
	suite.repoMock.On("Paging", requestPaging).Return(nil, dto.Paging{}, fmt.Errorf("error"))

	forms, paging, err := suite.usecase.FindAllForm(requestPaging)

	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "failed to retrieve forms: error")

	assert.Empty(suite.T(), forms)
	assert.Empty(suite.T(), paging)
}

func (suite *FormUseCaseTestSuite) TestFindByIdForm_Success() {
	dummy := formDummy[0]

	suite.repoMock.On("Get", dummy.FormID).Return(dummy, nil)

	form, err := suite.usecase.FindByIdForm(dummy.FormID)

	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), dummy, form)
}

func (suite *FormUseCaseTestSuite) TestFindByIdForm_NotFound() {

	suite.repoMock.On("Get", "1234").Return(model.Form{}, fmt.Errorf("error"))

	form, err := suite.usecase.FindByIdForm("1234")

	expectError := fmt.Sprintf("failed to retrieve form: error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectError)

	assert.Empty(suite.T(), form)
}


func (suite *FormUseCaseTestSuite) TestDeleteForm_FormNotFound() {
	nonExistentFormID := "1234"
	suite.repoMock.On("Get", nonExistentFormID).Return(model.Form{}, fmt.Errorf("error"))

	// Call method yang ingin diuji
	err := suite.usecase.DeleteForm(nonExistentFormID)

	// Periksa bahwa error yang diharapkan muncul
	expectedError := fmt.Sprintf("form with ID %s not found", nonExistentFormID)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *FormUseCaseTestSuite) TestDeleteForm_Failure() {

	formToDelete := formDummy[0]
	suite.repoMock.On("Get", formToDelete.FormID).Return(formToDelete, nil)
	suite.repoMock.On("Delete", formToDelete.FormID).Return(fmt.Errorf("error"))

	err := suite.usecase.DeleteForm(formToDelete.FormID)

	expectedError := fmt.Sprintf("failed to delete form: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *FormUseCaseTestSuite) TestUpdateForm_Success() {
	updatedForm := formDummy[0]
	suite.repoMock.On("Update", updatedForm).Return(nil)

	err := suite.usecase.UpdateForm(updatedForm)

	assert.Nil(suite.T(), err)
}

func (suite *FormUseCaseTestSuite) TestUpdateForm_Failure() {
	updatedForm := formDummy[0]

	// Set up an expectation for the Update method with the expected argument
	expectedError := fmt.Errorf("repository error")
	suite.repoMock.On("Update", updatedForm).Return(expectedError)

	// Make the usecase call
	err := suite.usecase.UpdateForm(updatedForm)

	// Assert that the error matches the expected error
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, fmt.Sprintf("failed to update form: %v", expectedError))

	// Verify that the expected method call was made
	suite.repoMock.AssertExpectations(suite.T())
}

