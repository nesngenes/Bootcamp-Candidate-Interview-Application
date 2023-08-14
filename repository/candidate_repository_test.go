package repository

import (
	"database/sql"
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var candidateDummy = []model.Candidate{
	{
		CandidateID:        "1",
		FullName:           "John",
		Phone:              "1234567890",
		Email:              "john@example.com",
		DateOfBirth:        "2005-01-01",
		Address:            "Address",
		CvLink:             "https://example.com",
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
		CvLink:             "https://example.com",
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
		CvLink:             "https://example.com",
		Bootcamp:           model.Bootcamp{BootcampId: "1"},
		InstansiPendidikan: "SMK",
		HackerRank:         90,
	},
}

type CandidateRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    CandidateRepository
}

func (suite *CandidateRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewCandidateRepository(db)
}

func (suite *CandidateRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func TestCandidateRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CandidateRepositoryTestSuite))
}

func (suite *CandidateRepositoryTestSuite) TestCreateNewCandidate_Success() {
	dummy := candidateDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO candidate (.+)").WithArgs(dummy.CandidateID, dummy.FullName, dummy.Phone, dummy.Email, dummy.DateOfBirth, dummy.Address, dummy.CvLink, dummy.Bootcamp.BootcampId, dummy.InstansiPendidikan, dummy.HackerRank).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Create(dummy)
	assert.Nil(suite.T(), actualError)
	assert.NoError(suite.T(), actualError)
}

func (suite *CandidateRepositoryTestSuite) TestCreateNewCandidate_Fail() {
	dummy := candidateDummy[0]
	suite.mockSql.ExpectExec("INSERT INTO candidate (.+)").WithArgs(dummy.CandidateID, dummy.FullName, dummy.Phone, dummy.Email, dummy.DateOfBirth, dummy.Address, dummy.CvLink, dummy.Bootcamp.BootcampId, dummy.InstansiPendidikan, dummy.HackerRank).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Create(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *CandidateRepositoryTestSuite) TestListCandidate_Succes() {
	rows := sqlmock.NewRows([]string{"id", "full_name", "phone", "email", "date_of_birth", "address", "cv_link", "bootcamp_id", "bootcamp_name", "bootcamp_start", "bootcamp_end", "bootcamp_location", "instansi_pendidikan", "hackerrank_score"})
	for _, candidate := range candidateDummy {
		rows.AddRow(candidate.CandidateID, candidate.FullName, candidate.Phone, candidate.Email, candidate.DateOfBirth, candidate.Address, candidate.CvLink, candidate.Bootcamp.BootcampId, candidate.Bootcamp.Name, candidate.Bootcamp.StartDate, candidate.Bootcamp.EndDate, candidate.Bootcamp.Location, candidate.InstansiPendidikan, candidate.HackerRank)
	}
	suite.mockSql.ExpectQuery("SELECT (.+) from candidate (.+)").WillReturnRows(rows)
	candidates, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), candidates, 3)
	assert.Equal(suite.T(), candidateDummy[0], candidates[0])
	assert.Equal(suite.T(), candidateDummy[1], candidates[1])
	assert.Equal(suite.T(), candidateDummy[2], candidates[2])
}

func (suite *CandidateRepositoryTestSuite) TestListCandidate_Fail() {
	suite.mockSql.ExpectQuery("SELECT (.+) from candidate (.+)").WillReturnError(fmt.Errorf("error"))
	products, err := suite.repo.List()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), products)
}

func (suite *CandidateRepositoryTestSuite) TestGetCandidate_Success() {
	expectedCandidate := candidateDummy[0]
	rows := sqlmock.NewRows([]string{"id", "full_name", "phone", "email", "date_of_birth", "address", "cv_link", "bootcamp_id", "bootcamp_name", "bootcamp_start", "bootcamp_end", "bootcamp_location", "instansi_pendidikan", "hackerrank_score"})
	rows.AddRow(expectedCandidate.CandidateID, expectedCandidate.FullName, expectedCandidate.Phone, expectedCandidate.Email, expectedCandidate.DateOfBirth, expectedCandidate.Address, expectedCandidate.CvLink, expectedCandidate.Bootcamp.BootcampId, expectedCandidate.Bootcamp.Name, expectedCandidate.Bootcamp.StartDate, expectedCandidate.Bootcamp.EndDate, expectedCandidate.Bootcamp.Location, expectedCandidate.InstansiPendidikan, expectedCandidate.HackerRank)
	suite.mockSql.ExpectQuery("SELECT (.+) from candidate (.+) where c.id = ?").WithArgs(expectedCandidate.CandidateID).WillReturnRows(rows)
	actualCandidate, actualError := suite.repo.Get(expectedCandidate.CandidateID)
	assert.NoError(suite.T(), actualError)
	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), expectedCandidate, actualCandidate)
}

func (suite *CandidateRepositoryTestSuite) TestGetCandidate_Fail() {
	suite.mockSql.ExpectQuery("SELECT (.+) from candidate (.+) where c.id = ?").WithArgs("1xxx").WillReturnError(fmt.Errorf("error"))
	actualCandidate, err := suite.repo.Get("1xxx")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Candidate{}, actualCandidate)
}

func (suite *CandidateRepositoryTestSuite) TestGetCandidateByPhoneNumber_Success() {
	expectedCandidate := candidateDummy[0]
	rows := sqlmock.NewRows([]string{"id", "full_name", "phone", "email", "date_of_birth", "address", "cv_link", "bootcamp_id", "bootcamp_name", "bootcamp_start", "bootcamp_end", "bootcamp_location", "instansi_pendidikan", "hackerrank_score"})
	rows.AddRow(expectedCandidate.CandidateID, expectedCandidate.FullName, expectedCandidate.Phone, expectedCandidate.Email, expectedCandidate.DateOfBirth, expectedCandidate.Address, expectedCandidate.CvLink, expectedCandidate.Bootcamp.BootcampId, expectedCandidate.Bootcamp.Name, expectedCandidate.Bootcamp.StartDate, expectedCandidate.Bootcamp.EndDate, expectedCandidate.Bootcamp.Location, expectedCandidate.InstansiPendidikan, expectedCandidate.HackerRank)
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id WHERE phone ILIKE $1")).WithArgs("%" + expectedCandidate.Email + "%").WillReturnRows(rows)
	actualCandidate, actualError := suite.repo.GetByPhoneNumber(expectedCandidate.Email)
	assert.NoError(suite.T(), actualError)
	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), expectedCandidate, actualCandidate)
}
func (suite *CandidateRepositoryTestSuite) TestGetCandidateByPhoneNumber_Fail() {
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id WHERE phone ILIKE $1")).WithArgs("sc1291").WillReturnError(fmt.Errorf("error"))
	actualCandidate, err := suite.repo.GetByPhoneNumber("X1291")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Candidate{}, actualCandidate)
}

func (suite *CandidateRepositoryTestSuite) TestGetCandidateByEmail_Success() {
	expectedCandidate := candidateDummy[0]
	rows := sqlmock.NewRows([]string{"id", "full_name", "phone", "email", "date_of_birth", "address", "cv_link", "bootcamp_id", "bootcamp_name", "bootcamp_start", "bootcamp_end", "bootcamp_location", "instansi_pendidikan", "hackerrank_score"})
	rows.AddRow(expectedCandidate.CandidateID, expectedCandidate.FullName, expectedCandidate.Phone, expectedCandidate.Email, expectedCandidate.DateOfBirth, expectedCandidate.Address, expectedCandidate.CvLink, expectedCandidate.Bootcamp.BootcampId, expectedCandidate.Bootcamp.Name, expectedCandidate.Bootcamp.StartDate, expectedCandidate.Bootcamp.EndDate, expectedCandidate.Bootcamp.Location, expectedCandidate.InstansiPendidikan, expectedCandidate.HackerRank)
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id WHERE c.email ILIKE $1")).WithArgs("%" + expectedCandidate.Email + "%").WillReturnRows(rows)
	actualCandidate, actualError := suite.repo.GetByEmail(expectedCandidate.Email)
	assert.NoError(suite.T(), actualError)
	assert.Nil(suite.T(), actualError)
	assert.Equal(suite.T(), expectedCandidate, actualCandidate)
}
func (suite *CandidateRepositoryTestSuite) TestGetCandidateByEmail_Fail() {
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id where c.id = $1")).WithArgs("sc1291").WillReturnError(fmt.Errorf("error"))
	actualCandidate, err := suite.repo.GetByEmail("X1291")
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), model.Candidate{}, actualCandidate)
}

func (suite *CandidateRepositoryTestSuite) TestDeleteCandidate_Success() {
	suite.mockSql.ExpectExec("DELETE FROM candidate WHERE id=?").WithArgs(candidateDummy[0].CandidateID).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Delete(candidateDummy[0].CandidateID)
	assert.Nil(suite.T(), actualError)
}

func (suite *CandidateRepositoryTestSuite) TestDeleteCandidate_Fail() {
	suite.mockSql.ExpectExec("DELETE FROM candidate WHERE id=?").WithArgs("1ABC").WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Delete("1ABC")
	assert.Error(suite.T(), actualError)
}

func (suite *CandidateRepositoryTestSuite) TestUpdateCandidate_Success() {
	dummy := candidateDummy[0]
	suite.mockSql.ExpectExec(regexp.QuoteMeta("UPDATE candidate SET full_name = $2, phone = $3, email = $4, date_of_birth = $5, address = $6,cv_link = $7,bootcamp_id = $8,instansi_pendidikan= $9,hackerrank_score = $10  WHERE id = $1")).WithArgs(dummy.CandidateID, dummy.FullName, dummy.Phone, dummy.Email, dummy.DateOfBirth, dummy.Address, dummy.CvLink, dummy.Bootcamp.BootcampId, dummy.InstansiPendidikan, dummy.HackerRank).WillReturnResult(sqlmock.NewResult(1, 1))
	actualError := suite.repo.Update(dummy)
	assert.Nil(suite.T(), actualError)
}

func (suite *CandidateRepositoryTestSuite) TestUpdateCandidate_Fail() {
	dummy := candidateDummy[0]
	suite.mockSql.ExpectExec(regexp.QuoteMeta("UPDATE candidate SET full_name = $2, phone = $3, email = $4, date_of_birth = $5, address = $6,cv_link = $7,bootcamp_id = $8,instansi_pendidikan= $9,hackerrank_score = $10  WHERE id = $1")).WithArgs(dummy.CandidateID, dummy.FullName, dummy.Phone, dummy.Email, dummy.DateOfBirth, dummy.Address, dummy.CvLink, dummy.Bootcamp.BootcampId, dummy.InstansiPendidikan, dummy.HackerRank).WillReturnError(fmt.Errorf("error"))
	actualError := suite.repo.Update(dummy)
	assert.Error(suite.T(), actualError)
}

func (suite *CandidateRepositoryTestSuite) TestPagingCandidate_Success() {
	// err := common.LoadEnv()
	// execptions.CheckErr(err)
	requestPaging := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	paginationQuery := dto.PaginationQuery{
		Take: 5,
		Skip: 0,
	}

	rows := sqlmock.NewRows([]string{"id", "full_name", "phone", "email", "date_of_birth", "address", "cv_link", "bootcamp_id", "bootcamp_name", "bootcamp_start", "bootcamp_end", "bootcamp_location", "instansi_pendidikan", "hackerrank_score"})
	for _, candidate := range candidateDummy {
		rows.AddRow(candidate.CandidateID, candidate.FullName, candidate.Phone, candidate.Email, candidate.DateOfBirth, candidate.Address, candidate.CvLink, candidate.Bootcamp.BootcampId, candidate.Bootcamp.Name, candidate.Bootcamp.StartDate, candidate.Bootcamp.EndDate, candidate.Bootcamp.Location, candidate.InstansiPendidikan, candidate.HackerRank)
	}
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id LIMIT $1 OFFSET $2")).WithArgs(paginationQuery.Take, paginationQuery.Skip).WillReturnRows(rows)

	// COUNT
	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(3)
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM candidate")).WillReturnRows(rowCount)

	actualCandidates, actualPaging, actualError := suite.repo.Paging(requestPaging)
	assert.Nil(suite.T(), actualError)
	assert.NotNil(suite.T(), actualCandidates)
	assert.Equal(suite.T(), actualPaging.TotalRows, 3)
}

func (suite *CandidateRepositoryTestSuite) TestPagingCandidat_QueryPagingError() {
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id LIMIT $1 OFFSET $2")).WillReturnError(fmt.Errorf("error"))
	actualCandidates, actualPaging, actualError := suite.repo.Paging(dto.PaginationParam{})
	assert.Error(suite.T(), actualError)
	assert.Nil(suite.T(), actualCandidates)
	assert.Equal(suite.T(), actualPaging.TotalRows, 0)
}

func (suite *CandidateRepositoryTestSuite) TestPagingCandidat_QueryCountError() {
	requestPaging := dto.PaginationParam{
		Page:  1,
		Limit: 5,
	}

	paginationQuery := dto.PaginationQuery{
		Take: 5,
		Skip: 0,
	}

	rows := sqlmock.NewRows([]string{"id", "full_name", "phone", "email", "date_of_birth", "address", "cv_link", "bootcamp_id", "bootcamp_name", "bootcamp_start", "bootcamp_end", "bootcamp_location", "instansi_pendidikan", "hackerrank_score"})
	for _, candidate := range candidateDummy {
		rows.AddRow(candidate.CandidateID, candidate.FullName, candidate.Phone, candidate.Email, candidate.DateOfBirth, candidate.Address, candidate.CvLink, candidate.Bootcamp.BootcampId, candidate.Bootcamp.Name, candidate.Bootcamp.StartDate, candidate.Bootcamp.EndDate, candidate.Bootcamp.Location, candidate.InstansiPendidikan, candidate.HackerRank)
	}
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id LIMIT $1 OFFSET $2")).WithArgs(paginationQuery.Take, paginationQuery.Skip).WillReturnRows(rows)
	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(3)
	suite.mockSql.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM candidate")).WillReturnError(fmt.Errorf("error"))
	_, actualPaging, actualError := suite.repo.Paging(requestPaging)
	assert.Error(suite.T(), actualError)
	assert.Equal(suite.T(), dto.Paging{}, actualPaging)
}
