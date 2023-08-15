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

type bootcampRepoMock struct {
	mock.Mock
}

func (b *bootcampRepoMock) Create(payload model.Bootcamp) error {
	args := b.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (b *bootcampRepoMock) Delete(id string) error {
	args := b.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (b *bootcampRepoMock) Get(id string) (model.Bootcamp, error) {
	args := b.Called(id)
	if args.Get(1) != nil {
		return model.Bootcamp{}, args.Error(1)
	}
	return args.Get(0).(model.Bootcamp), nil
}

func (b *bootcampRepoMock) GetByName(name string) (model.Bootcamp, error) {
	args := b.Called(name)
	if args.Get(1) != nil {
		return model.Bootcamp{}, args.Error(1)
	}
	return args.Get(0).(model.Bootcamp), nil

}

func (m *bootcampRepoMock) GetByID(id string) (model.Bootcamp, error) {
	args := m.Called(id)
	if args.Get(1) != nil {
		return model.Bootcamp{}, args.Error(1)
	}
	return args.Get(0).(model.Bootcamp), args.Error(1)
}

func (b *bootcampRepoMock) List() ([]model.Bootcamp, error) {
	panic("unimplemented")
}

func (b *bootcampRepoMock) Paging(requestPaging dto.PaginationParam) ([]model.Bootcamp, dto.Paging, error) {
	args := b.Called(requestPaging)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Bootcamp), args.Get(1).(dto.Paging), nil

}

func (b *bootcampRepoMock) Update(payload model.Bootcamp) error {
	args := b.Called(payload)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil

}

type BootcampUseCaseTestSuite struct {
	suite.Suite
	bootcampRepoMock *bootcampRepoMock
	usecase        BootcampUseCase
}

func (suite *BootcampUseCaseTestSuite) SetupTest() {
	suite.bootcampRepoMock = new(bootcampRepoMock)
	suite.usecase = NewBootcampUseCase(suite.bootcampRepoMock)
}

var bootcampsDummy = []model.Bootcamp{
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
    {
        BootcampId: "3",
        Name:       "Bootcamp 3",
        StartDate:  parseTime("2023-03-02T00:00:00Z"),
        EndDate:    parseTime("2023-03-15T00:00:00Z"),
        Location:   "JAPAN",
    },
}

func parseTime(timeStr string) time.Time {
    parsedTime, err := time.Parse(time.RFC3339, timeStr)
    if err != nil {
        panic(err) // Handle the error appropriately in your code
    }
    return parsedTime
}


func (suite *BootcampUseCaseTestSuite) TestRegisterNewBootcamp_Succes() {
	newBootcamp := bootcampsDummy[0]
	suite.bootcampRepoMock.On("GetByName", newBootcamp.Name).Return(model.Bootcamp{}, fmt.Errorf("error"))
	suite.bootcampRepoMock.On("Create", newBootcamp).Return(nil)

	err := suite.usecase.RegisterNewBootcamp(newBootcamp)

	assert.Nil(suite.T(), err)
}

func (suite *BootcampUseCaseTestSuite) TestRegisterNewBootcamp_EmptyName() {
	emptyBootcamp := model.Bootcamp{Name: ""}

	err := suite.usecase.RegisterNewBootcamp(emptyBootcamp)

	expectedError := "name ,location strat_date and end_date required fields"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *BootcampUseCaseTestSuite) TestRegisterNewBootcamp_Failure() {
	newBootcamp := model.Bootcamp{Name: "New Bootcamp"}
	suite.bootcampRepoMock.On("GetByName", newBootcamp.Name).Return(model.Bootcamp{}, fmt.Errorf("error"))
	suite.bootcampRepoMock.On("Create", newBootcamp).Return(fmt.Errorf("error"))

	err := suite.usecase.RegisterNewBootcamp(newBootcamp)

	expectedError := fmt.Sprintf("name ,location strat_date and end_date required fields")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)

}

func (suite *BootcampUseCaseTestSuite) TestFindAllBootcamp_Success() {
	dummy := bootcampsDummy
	expectedPaging := dto.Paging{
		Page:        1,
		RowsPerPage: 5,
		TotalRows:   3,
		TotalPages:  1,
	}
	requestPaging := dto.PaginationParam{
		Page: 1,
	}
	suite.bootcampRepoMock.On("Paging", requestPaging).Return(dummy, expectedPaging, nil)
	bootcamps, paging, err := suite.usecase.FindAllBootcamp(requestPaging)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummy, bootcamps)
	assert.Equal(suite.T(), expectedPaging, paging)
}

func (suite *BootcampUseCaseTestSuite) TestFindAllBootcamp_Failure() {
	requestPaging := dto.PaginationParam{
		Page: 1,
	}
	suite.bootcampRepoMock.On("Paging", requestPaging).Return(nil, dto.Paging{}, fmt.Errorf("error"))

	bootcamps, paging, err := suite.usecase.FindAllBootcamp(requestPaging)

	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "error")
	assert.Empty(suite.T(), bootcamps)
	assert.Empty(suite.T(), paging)

}

func (suite *BootcampUseCaseTestSuite) TestUpdateBootcamp_Success() {
	dummy := bootcampsDummy[0]
	suite.bootcampRepoMock.On("Update", dummy).Return(nil)

	err := suite.usecase.UpdateBootcamp(dummy)

	assert.Nil(suite.T(), err)

	suite.bootcampRepoMock.AssertCalled(suite.T(), "Update", dummy)
}

func (suite *BootcampUseCaseTestSuite) TestUpdateBootcamp_EmptyName() {
	emptyBootcamp := model.Bootcamp{Name: ""}

	err := suite.usecase.UpdateBootcamp(emptyBootcamp)

	expectedError := "name ,location strat_date and end_date required fields"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *BootcampUseCaseTestSuite) TestUpdateBootcamp_Failure() {
	dummy := bootcampsDummy[0]
	suite.bootcampRepoMock.On("Update", dummy).Return(fmt.Errorf("error"))

	err := suite.usecase.UpdateBootcamp(dummy)

	expectedError := fmt.Sprintf("failed to update bootcamp: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)

}

func (suite *BootcampUseCaseTestSuite) TestFindByIdBootcamp_Success() {
	dummyBootcamp := bootcampsDummy[0]
	suite.bootcampRepoMock.On("Get", dummyBootcamp.BootcampId).Return(dummyBootcamp, nil)

	bootcamp, err := suite.usecase.FindByIdBootcamp(dummyBootcamp.BootcampId)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyBootcamp, bootcamp)
}

func (suite *BootcampUseCaseTestSuite) TestFindByIdBootcamp_NotFound() {
	suite.bootcampRepoMock.On("Get", "1234").Return(model.Bootcamp{}, fmt.Errorf("error"))

	bootcamp, err := suite.usecase.FindByIdBootcamp("1234")

	expectedError := "bootcamp with id 1234 not found"
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
	assert.Equal(suite.T(), model.Bootcamp{}, bootcamp)
}
func (suite *BootcampUseCaseTestSuite) TestDeleteBootcamp_Success() {
	dummyBootcamp := bootcampsDummy[0]
	suite.bootcampRepoMock.On("Get", dummyBootcamp.BootcampId).Return(dummyBootcamp, nil)
	suite.bootcampRepoMock.On("Delete", dummyBootcamp.BootcampId).Return(nil)

	err := suite.usecase.DeleteBootcamp(dummyBootcamp.BootcampId)

	assert.Nil(suite.T(), err)
}

func (suite *BootcampUseCaseTestSuite) TestDeleteBootcamp_BootcampNotFound() {
	nonExistentBootcampID := "1234"
	suite.bootcampRepoMock.On("Get", nonExistentBootcampID).Return(model.Bootcamp{}, fmt.Errorf("error"))

	err := suite.usecase.DeleteBootcamp(nonExistentBootcampID)

	expectedError := fmt.Sprintf("bootcamp with ID %s not found", nonExistentBootcampID)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func (suite *BootcampUseCaseTestSuite) TestDeleteBootcamp_Failure() {
	dummyBootcamp := bootcampsDummy[0]
	suite.bootcampRepoMock.On("Get", dummyBootcamp.BootcampId).Return(dummyBootcamp, nil)
	suite.bootcampRepoMock.On("Delete", dummyBootcamp.BootcampId).Return(fmt.Errorf("error"))

	err := suite.usecase.DeleteBootcamp(dummyBootcamp.BootcampId)

	expectedError := fmt.Sprintf("failed to delete bootcamp: %v", "error")
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, expectedError)
}

func TestBootcampUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(BootcampUseCaseTestSuite))
}
