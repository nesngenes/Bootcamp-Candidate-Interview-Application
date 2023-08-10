package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type repoMock struct {
	mock.Mock
}

func (r *repoMock) Create(payload model.Candidate) error {
	args := r.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *repoMock) List() ([]model.Candidate, error) {
	args := r.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Candidate), nil
}

func (r *repoMock) Get(id string) (model.Candidate, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.Candidate{}, args.Error(1)
	}
	return args.Get(0).(model.Candidate), nil
}
func (r *repoMock) GetByPhoneNumber(phoneNumber string) (model.Candidate, error) {
	args := r.Called(phoneNumber)
	if args.Get(1) != nil {
		return model.Candidate{}, args.Error(1)
	}
	return args.Get(0).(model.Candidate), nil
}
func (r *repoMock) GetByEmail(email string) (model.Candidate, error) {
	args := r.Called(email)
	if args.Get(1) != nil {
		return model.Candidate{}, args.Error(1)
	}
	return args.Get(0).(model.Candidate), nil
}

func (r *repoMock) Update(payload model.Candidate) error {
	args := r.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *repoMock) Delete(id string) error {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
func (r *repoMock) Paging(requestPaging dto.PaginationParam) ([]model.Candidate, dto.Paging, error) {
	args := r.Called(requestPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Candidate), args.Get(1).(dto.Paging), nil
}

type usecaseMock struct {
	mock.Mock
}

// FindByIdBootcamp implements BootcampUseCase.
func (*usecaseMock) FindByIdBootcamp(id string) (model.Bootcamp, error) {
	panic("unimplemented")
}

// GetBootcampByID implements BootcampUseCase.
func (u *usecaseMock) GetBootcampByID(id string) (model.Bootcamp, error) {
	args := u.Called(id)
	if args.Get(1) != nil {
		return model.Bootcamp{}, args.Error(1)
	}
	return args.Get(0).(model.Bootcamp), nil
}

func (*usecaseMock) DeleteBootcamp(id string) error {
	panic("unimplemented")
}

func (u *usecaseMock) FindAllBootcamp(requesPaging dto.PaginationParam) ([]model.Bootcamp, dto.Paging, error) {
	panic("unimplemented")
}

func (u *usecaseMock) RegisterNewBootcamp(payload model.Bootcamp) error {
	panic("unimplemented")
}

func (u *usecaseMock) UpdateBootcamp(payload model.Bootcamp) error {
	panic("unimplemented")
}

type CandidateUseCaseTestSuite struct {
	suite.Suite
	repoMock    *repoMock
	usecaseMock *usecaseMock
	usecase     CandidateUseCase
}

func (suite *CandidateUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(repoMock)
	suite.usecaseMock = new(usecaseMock)
	suite.usecase = NewCandidateUseCase(suite.repoMock, suite.usecaseMock)
}

func TestCandidateUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CandidateUseCaseTestSuite))
}

var candidateDummy = []model.Candidate{
	{
		CandidateID:        "1",
		FullName:           "John",
		Phone:              "1234567890",
		Email:              "john@example.com",
		DateOfBirth:        "2005-01-01",
		Address:            "Address",
		Bootcamp:           model.Bootcamp{BootcampId: "1"},
		InstansiPendidikan: "SMK",
		HackerRank:         90,
	},
	{
		CandidateID:        "2",
		FullName:           "John",
		Phone:              "1234567891",
		Email:              "johndoe@example.com",
		DateOfBirth:        "2005-01-02",
		Address:            "Address 2",
		Bootcamp:           model.Bootcamp{BootcampId: "1"},
		InstansiPendidikan: "SMK",
		HackerRank:         90,
	},
	{
		CandidateID:        "3",
		FullName:           "John",
		Phone:              "1234567892",
		Email:              "johntor@example.com",
		DateOfBirth:        "2005-01-03",
		Address:            "Address 3",
		Bootcamp:           model.Bootcamp{BootcampId: "1"},
		InstansiPendidikan: "SMK",
		HackerRank:         90,
	},
}

var bootcampDummy = model.Bootcamp{
	BootcampId: "1",
	Name:       "Golang",
	StartDate:  time.Date(2023, time.July, 17, 0, 0, 0, 0, time.UTC),
	EndDate:    time.Date(2023, time.September, 17, 0, 0, 0, 0, time.UTC),
	Location:   "Online",
}

func (suite *CandidateUseCaseTestSuite) TestRegisterNewProduct_Success() {
	dmProduct := candidateDummy[0]
	suite.repoMock.On("GetByEmail", dmProduct.Email).Return(model.Candidate{}, fmt.Errorf("error"))
	suite.repoMock.On("GetByPhoneNumber", dmProduct.Phone).Return(model.Candidate{}, fmt.Errorf("error"))
	suite.usecaseMock.On("FindByIdBootcamp", dmProduct.Bootcamp.BootcampId).Return(bootcampDummy, nil)
	suite.repoMock.On("Create", dmProduct).Return(nil)
	err := suite.usecase.RegisterNewCandidate(dmProduct)
	assert.Nil(suite.T(), err)
}
func (suite *CandidateUseCaseTestSuite) TestRegisterNewCandidate_EmptyFields() {
	emptyCandidate := model.Candidate{}
	err := suite.usecase.RegisterNewCandidate(emptyCandidate)

	expectedError := "fullname, email, phone, address, date of birth required fields"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *CandidateUseCaseTestSuite) TestRegisterNewCandidate_DuplicateEmail() {
	duplicateEmailCandidate := candidateDummy[0]

	suite.repoMock.On("GetByEmail", duplicateEmailCandidate.Email).Return(duplicateEmailCandidate, nil)

	err := suite.usecase.RegisterNewCandidate(duplicateEmailCandidate)

	expectedError := fmt.Sprintf("candidate with email %s exists", duplicateEmailCandidate.Email)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *CandidateUseCaseTestSuite) TestRegisterNewCandidate_DuplicatePhoneNumber() {
	duplicatePhoneCandidate := candidateDummy[0]

	suite.repoMock.On("GetByEmail", duplicatePhoneCandidate.Email).Return(model.Candidate{}, fmt.Errorf("error"))

	suite.repoMock.On("GetByPhoneNumber", duplicatePhoneCandidate.Phone).Return(duplicatePhoneCandidate, nil)

	err := suite.usecase.RegisterNewCandidate(duplicatePhoneCandidate)

	expectedError := fmt.Sprintf("candidate with phoone %s exists", duplicatePhoneCandidate.Phone)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *CandidateUseCaseTestSuite) TestRegisterNewCandidate_BootcampNotFound() {
	nonExistentBootcampCandidate := candidateDummy[0]

	suite.repoMock.On("GetByEmail", nonExistentBootcampCandidate.Email).Return(model.Candidate{}, fmt.Errorf("error"))

	suite.repoMock.On("GetByPhoneNumber", nonExistentBootcampCandidate.Phone).Return(model.Candidate{}, fmt.Errorf("error"))

	suite.usecaseMock.On("FindByIdBootcamp", nonExistentBootcampCandidate.Bootcamp.BootcampId).Return(model.Bootcamp{}, fmt.Errorf("error"))

	err := suite.usecase.RegisterNewCandidate(nonExistentBootcampCandidate)

	expectedError := fmt.Sprintf("bootcamp with ID %s not found", nonExistentBootcampCandidate.Bootcamp.BootcampId)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)

}

func (suite *CandidateUseCaseTestSuite) TestRegisterNewCandidate_Fail() {
	failedCreateCandidate := candidateDummy[0]

	suite.repoMock.On("GetByEmail", failedCreateCandidate.Email).Return(model.Candidate{}, fmt.Errorf("error"))

	suite.repoMock.On("GetByPhoneNumber", failedCreateCandidate.Phone).Return(model.Candidate{}, fmt.Errorf("error"))

	suite.usecaseMock.On("FindByIdBootcamp", failedCreateCandidate.Bootcamp.BootcampId).Return(bootcampDummy, nil)

	suite.repoMock.On("Create", failedCreateCandidate).Return(fmt.Errorf("error"))

	err := suite.usecase.RegisterNewCandidate(failedCreateCandidate)

	expectedError := fmt.Sprintf("failed to create new candidate: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)

}

func (suite *CandidateUseCaseTestSuite) TestFindAllCandidate_Success() {
	dummy := candidateDummy
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

	candidates, paging, err := suite.usecase.FindAllCandidate(requestPaging)

	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), dummy, candidates)
	assert.Equal(suite.T(), expectedPaging, paging)

}

func (suite *CandidateUseCaseTestSuite) TestFindAllCandidate_Fail() {
	requestPaging := dto.PaginationParam{
		Page: 1,
	}
	suite.repoMock.On("Paging", requestPaging).Return(nil, dto.Paging{}, fmt.Errorf("error"))

	candidates, paging, err := suite.usecase.FindAllCandidate(requestPaging)

	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "error")

	assert.Empty(suite.T(), candidates)
	assert.Empty(suite.T(), paging)
}

func (suite *CandidateUseCaseTestSuite) TestFindByIdCandidate_Success() {
	dummy := candidateDummy[0]

	suite.repoMock.On("Get", dummy.CandidateID).Return(dummy, nil)

	candidate, err := suite.usecase.FindByIdCandidate(dummy.CandidateID)

	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), dummy, candidate)
}

func (suite *CandidateUseCaseTestSuite) TestFindByIdCandidate_NotFound() {

	suite.repoMock.On("Get", "1234").Return(model.Candidate{}, fmt.Errorf("error"))

	candidate, err := suite.usecase.FindByIdCandidate("1234")

	expectError := fmt.Sprintf("candidate with id %v not found", "1234")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectError)

	assert.Empty(suite.T(), candidate)
}

func (suite *CandidateUseCaseTestSuite) TestDeleteCandidate_Success() {

	dummy := candidateDummy[0]
	suite.repoMock.On("Get", dummy.CandidateID).Return(dummy, nil)
	suite.repoMock.On("Delete", dummy.CandidateID).Return(nil)

	err := suite.usecase.DeleteCandidate(dummy.CandidateID)

	assert.Nil(suite.T(), err)
}

func (suite *CandidateUseCaseTestSuite) TestDeleteCandidate_CandidateNotFound() {
	nonExistentCandidateID := "1234"
	suite.repoMock.On("Get", nonExistentCandidateID).Return(model.Candidate{}, fmt.Errorf("error"))

	err := suite.usecase.DeleteCandidate(nonExistentCandidateID)

	expectedError := fmt.Sprintf("candidate with ID %s not found", nonExistentCandidateID)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)

}

func (suite *CandidateUseCaseTestSuite) TestDeleteCandidate_Failure() {

	candidateToDelete := candidateDummy[0]
	suite.repoMock.On("Get", candidateToDelete.CandidateID).Return(candidateToDelete, nil)
	suite.repoMock.On("Delete", candidateToDelete.CandidateID).Return(fmt.Errorf("error"))

	err := suite.usecase.DeleteCandidate(candidateToDelete.CandidateID)

	expectedError := fmt.Sprintf("failed to delete candidate: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *CandidateUseCaseTestSuite) TestUpdateCandidate_Success() {
	updatedCandidate := candidateDummy[0]
	suite.usecaseMock.On("FindByIdBootcamp", updatedCandidate.Bootcamp.BootcampId).Return(bootcampDummy, nil)
	suite.repoMock.On("Update", updatedCandidate).Return(nil)

	err := suite.usecase.UpdateCandidate(updatedCandidate)

	assert.Nil(suite.T(), err)
}

func (suite *CandidateUseCaseTestSuite) TestUpdateCandidate_BootcampNotFound() {
	updatedCandidate := candidateDummy[0]
	suite.usecaseMock.On("FindByIdBootcamp", updatedCandidate.Bootcamp.BootcampId).Return(model.Bootcamp{}, fmt.Errorf("error"))

	err := suite.usecase.UpdateCandidate(updatedCandidate)

	expectedError := fmt.Sprintf("bootcamp with ID %s not found", updatedCandidate.Bootcamp.BootcampId)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *CandidateUseCaseTestSuite) TestUpdateCandidate_EmptyPhone() {
	updatedCandidate := candidateDummy[0]
	updatedCandidate.Phone = ""

	err := suite.usecase.UpdateCandidate(updatedCandidate)

	expectedError := "kolom nomor harus di isi"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *CandidateUseCaseTestSuite) TestUpdateCandidate_Failure() {
	updatedCandidate := candidateDummy[0]
	suite.usecaseMock.On("FindByIdBootcamp", updatedCandidate.Bootcamp.BootcampId).Return(bootcampDummy, nil)
	suite.repoMock.On("Update", updatedCandidate).Return(fmt.Errorf("error"))

	err := suite.usecase.UpdateCandidate(updatedCandidate)

	expectedError := fmt.Sprintf("gagal memperbarui nomor: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}
