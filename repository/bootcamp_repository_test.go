package repository

import (
	"database/sql"
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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

type BootcampRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    BootcampRepository
}

func (suite *BootcampRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewBootcampRepository(db)
}

func (suite *BootcampRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestBootcampRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BootcampRepositoryTestSuite))
}

func (suite *BootcampRepositoryTestSuite) TestCreateBootcamp_Success() {
	dummy := bootcampsDummy[0]
	query := "INSERT INTO bootcamp (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.BootcampId, dummy.Name, dummy.StartDate, dummy.EndDate, dummy.Location).WillReturnResult(sqlmock.NewResult(1, 1))

	actualError := suite.repo.Create(dummy)

	assert.Nil(suite.T(), actualError)
}

func (suite *BootcampRepositoryTestSuite) TestCreateBootcamp_Fai() {
	dummy := bootcampsDummy[0]
	query := "INSERT INTO bootcamp (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.BootcampId, dummy.Name, dummy.StartDate, dummy.EndDate, dummy.Location).WillReturnError(fmt.Errorf("error"))

	actualError := suite.repo.Create(dummy)

	expectedError := fmt.Errorf("error")
	assert.NotNil(suite.T(), actualError)
	assert.EqualError(suite.T(), actualError, expectedError.Error())
}

func (suite *BootcampRepositoryTestSuite) TestGetBootcampByName_Success() {
	dummy := bootcampsDummy[0]
	rows := sqlmock.NewRows([]string{"id", "name","start_date","end_date","location"})
	rows.AddRow(dummy.BootcampId, dummy.Name, dummy.StartDate, dummy.EndDate, dummy.Location)

	query := "SELECT (.+) FROM bootcamp WHERE name ILIKE (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs("%" + dummy.Name + "%").WillReturnRows(rows)

	actualBootcamp, actualError := suite.repo.GetByName(dummy.Name)

	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), actualBootcamp, dummy)
}

func (suite *BootcampRepositoryTestSuite) TestGetBootcampByName_NotFound() {
	dummyBootcampName := "Nonexistent Bootcamp"
	query := "SELECT (.+) FROM bootcamp WHERE name ILIKE (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs("%" + dummyBootcampName + "%").
		WillReturnError(sql.ErrNoRows)

	actualBootcamp, actualError := suite.repo.GetByName(dummyBootcampName)

	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), model.Bootcamp{}, actualBootcamp)
}

func (suite *BootcampRepositoryTestSuite) TestListBootcamp_Succes() {
	rows := sqlmock.NewRows([]string{"id", "name","start_date","end_date","location"})
	for _, dummy := range bootcampsDummy {
		rows.AddRow(dummy.BootcampId, dummy.Name, dummy.StartDate, dummy.EndDate, dummy.Location)
	}

	query := "SELECT (.+) FROM bootcamp"
	suite.mockSql.ExpectQuery(query).WillReturnRows(rows)
	statuses, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), statuses, 3)
	assert.Equal(suite.T(), bootcampsDummy[0], statuses[0])
	assert.Equal(suite.T(), bootcampsDummy[1], statuses[1])
	assert.Equal(suite.T(), bootcampsDummy[2], statuses[2])

}

func (suite *BootcampRepositoryTestSuite) TestListBootcamp_Fail() {
	query := "SELECT (.+) FROM bootcamp"
	suite.mockSql.ExpectQuery(query).WillReturnError(fmt.Errorf("error"))
	bootcamps, err := suite.repo.List()
	assert.Nil(suite.T(), bootcamps)
	assert.Error(suite.T(), err)
}

func (suite *BootcampRepositoryTestSuite) TestGetBootcampByID_Success() {
	dummy := bootcampsDummy[0]
	rows := sqlmock.NewRows([]string{"id", "name","start_date","end_date","location"})
	rows.AddRow(dummy.BootcampId, dummy.Name, dummy.StartDate, dummy.EndDate, dummy.Location)

	query := "SELECT (.+) FROM bootcamp WHERE id= (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs(dummy.BootcampId).WillReturnRows(rows)

	actualBootcamp, actualError := suite.repo.Get(dummy.BootcampId)

	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), actualBootcamp, dummy)
}

func (suite *BootcampRepositoryTestSuite) TestGetBootcampByID_NotFound() {
	dummyBootcampID := "3rujwke"
	query := "SELECT (.+) FROM bootcamp WHERE id= (.+)"
	suite.mockSql.ExpectQuery(query).WithArgs(dummyBootcampID).
		WillReturnError(sql.ErrNoRows)

	actualBootcamp, actualError := suite.repo.Get(dummyBootcampID)

	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), model.Bootcamp{}, actualBootcamp)
}

func (suite *BootcampRepositoryTestSuite) TestUpdateBootcamp_Success() {
	dummy := bootcampsDummy[0]
	query := "UPDATE bootcamp SET name = (.+) WHERE id = (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.BootcampId, dummy.Name, dummy.StartDate, dummy.EndDate, dummy.Location).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Update(dummy)

	assert.Nil(suite.T(), actualError)
	assert.NoError(suite.T(), actualError)
}

func (suite *BootcampRepositoryTestSuite) TestUpdateBootcamp_Fail() {
	dummy := bootcampsDummy[0]
	query := "UPDATE bootcamp SET name = (.+) WHERE id = (.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.BootcampId, dummy.Name, dummy.StartDate, dummy.EndDate, dummy.Location).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Update(dummy)

	assert.Error(suite.T(), actualError)
}

func (suite *BootcampRepositoryTestSuite) TestDeleteBootcamp_Success() {
	dummy := bootcampsDummy[0]
	query := "DELETE FROM bootcamp WHERE id=(.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.BootcampId).WillReturnResult(sqlmock.NewResult(0, 1))
	actualError := suite.repo.Delete(dummy.BootcampId)
	assert.Nil(suite.T(), actualError)
}

func (suite *BootcampRepositoryTestSuite) TestDeleteBootcamp_Fail() {
	dummy := bootcampsDummy[0]
	query := "DELETE FROM bootcamp WHERE id=(.+)"
	suite.mockSql.ExpectExec(query).WithArgs(dummy.BootcampId).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Delete(dummy.BootcampId)
	assert.Error(suite.T(), actualError)
}

func (suite *BootcampRepositoryTestSuite) TestPagingBootcamp_Succes() {
	requestPaging := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	paginationQuery := dto.PaginationQuery{
		Take: 5,
		Skip: 0,
	}
	rows := sqlmock.NewRows([]string{"id", "name","start_date","end_date","location"})
	for _, dummy := range bootcampsDummy {
		rows.AddRow(dummy.BootcampId, dummy.Name, dummy.StartDate, dummy.EndDate, dummy.Location)
	}
	selectQuery := "SELECT (.+) FROM bootcamp  LIMIT (.+) OFFSET (.+)"
	countQuery := "SELECT COUNT\\(.*\\) FROM bootcamp"
	suite.mockSql.ExpectQuery(selectQuery).WithArgs(paginationQuery.Take, paginationQuery.Skip).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(3)
	suite.mockSql.ExpectQuery(countQuery).WillReturnRows(rowCount)
	actualBoocamp, actualPaging, actualError := suite.repo.Paging(requestPaging)
	assert.Nil(suite.T(), actualError)
	assert.NotNil(suite.T(), actualBoocamp)
	assert.Equal(suite.T(), actualPaging.TotalRows, 3)
}

func (suite *BootcampRepositoryTestSuite) TestPagingBootcamp_QueryPagingError() {
	selectQuery := "SELECT (.+) FROM bootcamp  LIMIT (.+) OFFSET (.+)"
	suite.mockSql.ExpectQuery(selectQuery).WillReturnError(fmt.Errorf("error"))
	actualBootcamp, actualPaging, actualError := suite.repo.Paging(dto.PaginationParam{})
	assert.Error(suite.T(), actualError)
	assert.Nil(suite.T(), actualBootcamp)
	assert.Equal(suite.T(), actualPaging.TotalRows, 0)
}

func (suite *BootcampRepositoryTestSuite) TestPagingBootcamp_QueryCountError() {
	requestPaging := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	paginationQuery := dto.PaginationQuery{
		Take: 5,
		Skip: 0,
	}
	rows := sqlmock.NewRows([]string{"id", "name","start_date","end_date","location"})
	for _, dummy := range bootcampsDummy {
		rows.AddRow(dummy.BootcampId, dummy.Name, dummy.StartDate, dummy.EndDate, dummy.Location)
	}
	selectQuery := "SELECT (.+) FROM bootcamp  LIMIT (.+) OFFSET (.+)"
	countQuery := "SELECT COUNT\\(.*\\) FROM bootcamp"
	suite.mockSql.ExpectQuery(selectQuery).WithArgs(paginationQuery.Take, paginationQuery.Skip).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(3)
	suite.mockSql.ExpectQuery(countQuery).WillReturnError(fmt.Errorf("error"))
	_, actualPaging, actualError := suite.repo.Paging(requestPaging)
	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), actualPaging.TotalRows, 0)
}
