package repository

import (
	"database/sql"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/utils/common"
)

type CandidateRepository interface {
	BaseRepository[model.Candidate]
	GetByPhoneNumber(phoneNumber string) (model.Candidate, error)
	GetByEmail(email string) (model.Candidate, error)
	BaseRepositoryPaging[model.Candidate]
}

type candidateRepository struct {
	db *sql.DB
}

func (c *candidateRepository) Create(payload model.Candidate) error {
	_, err := c.db.Exec("INSERT INTO candidate (id, full_name, phone,email, date_of_birth, address,cv_link,bootcamp_id,instansi_pendidikan,hackerrank_score) VALUES ($1, $2, $3, $4, $5, $6, $7,$8,$9,$10)", payload.CandidateID, payload.FullName, payload.Phone, payload.Email, payload.DateOfBirth, payload.Address, payload.CvLink, payload.Bootcamp.BootcampId, payload.InstansiPendidikan, payload.HackerRank)
	if err != nil {
		return err
	}
	return nil

}

// GetPhoneNumber implements employeeRepository.
func (c *candidateRepository) GetByPhoneNumber(phoneNumber string) (model.Candidate, error) {
	var candidate model.Candidate
	err := c.db.QueryRow("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id WHERE phone ILIKE $1", "%"+phoneNumber+"%").Scan(&candidate.CandidateID, &candidate.FullName, &candidate.Phone, &candidate.Email, &candidate.DateOfBirth, &candidate.Address, &candidate.CvLink, &candidate.Bootcamp.BootcampId, &candidate.Bootcamp.Name, &candidate.Bootcamp.StartDate, &candidate.Bootcamp.EndDate, &candidate.Bootcamp.Location, &candidate.InstansiPendidikan, &candidate.HackerRank)
	if err != nil {
		return model.Candidate{}, err
	}
	return candidate, nil
}

func (c *candidateRepository) List() ([]model.Candidate, error) {
	rows, err := c.db.Query("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id")
	if err != nil {
		return nil, err
	}
	var candidates []model.Candidate
	for rows.Next() {
		var candidate model.Candidate
		err := rows.Scan(&candidate.CandidateID, &candidate.FullName, &candidate.Phone, &candidate.Email, &candidate.DateOfBirth, &candidate.Address, &candidate.CvLink, &candidate.Bootcamp.BootcampId, &candidate.Bootcamp.Name, &candidate.Bootcamp.StartDate, &candidate.Bootcamp.EndDate, &candidate.Bootcamp.Location, &candidate.InstansiPendidikan, &candidate.HackerRank)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}
	return candidates, nil
}

func (c *candidateRepository) Get(id string) (model.Candidate, error) {
	var candidate model.Candidate
	row := c.db.QueryRow("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id where c.id = $1", id)
	err := row.Scan(&candidate.CandidateID, &candidate.FullName, &candidate.Phone, &candidate.Email, &candidate.DateOfBirth, &candidate.Address, &candidate.CvLink, &candidate.Bootcamp.BootcampId, &candidate.Bootcamp.Name, &candidate.Bootcamp.StartDate, &candidate.Bootcamp.EndDate, &candidate.Bootcamp.Location, &candidate.InstansiPendidikan, &candidate.HackerRank)
	if err != nil {
		return model.Candidate{}, err
	}
	return candidate, nil

}

func (c *candidateRepository) GetByEmail(email string) (model.Candidate, error) {
	var candidate model.Candidate
	err := c.db.QueryRow("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id WHERE c.email ILIKE $1", "%"+email+"%").Scan(&candidate.CandidateID, &candidate.FullName, &candidate.Phone, &candidate.Email, &candidate.DateOfBirth, &candidate.Address, &candidate.CvLink, &candidate.Bootcamp.BootcampId, &candidate.Bootcamp.Name, &candidate.Bootcamp.StartDate, &candidate.Bootcamp.EndDate, &candidate.Bootcamp.Location, &candidate.InstansiPendidikan, &candidate.HackerRank)
	if err != nil {
		return model.Candidate{}, err
	}
	return candidate, nil

}

func (c *candidateRepository) Update(payload model.Candidate) error {
	_, err := c.db.Exec("UPDATE candidate SET full_name = $2, phone = $3, email = $4, date_of_birth = $5, address = $6,cv_link = $7,bootcamp_id = $8,instansi_pendidikan= $9,hackerrank_score = $10  WHERE id = $1", payload.CandidateID, payload.FullName, payload.Phone, payload.Email, payload.DateOfBirth, payload.Address, payload.CvLink, payload.Bootcamp.BootcampId, payload.InstansiPendidikan, payload.HackerRank)
	if err != nil {
		return err
	}
	return nil

}

func (c *candidateRepository) Delete(id string) error {
	_, err := c.db.Exec("DELETE FROM candidate WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (c *candidateRepository) Paging(requestPaging dto.PaginationParam) ([]model.Candidate, dto.Paging, error) {
	paginationQuery := common.GetPaginationParams(requestPaging)
	rows, err := c.db.Query("SELECT c.id,c.full_name,c.phone,c.email,c.date_of_birth,c.address,c.cv_link,b.id,b.name,b.start_date,b.end_date,b.location,c.instansi_pendidikan,c.hackerrank_score from candidate c INNER JOIN bootcamp b on b.id = c.bootcamp_id  LIMIT $1 OFFSET $2", paginationQuery.Take, paginationQuery.Skip)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var candidates []model.Candidate
	for rows.Next() {
		var candidate model.Candidate
		err := rows.Scan(&candidate.CandidateID, &candidate.FullName, &candidate.Phone, &candidate.Email, &candidate.DateOfBirth, &candidate.Address, &candidate.CvLink, &candidate.Bootcamp.BootcampId, &candidate.Bootcamp.Name, &candidate.Bootcamp.StartDate, &candidate.Bootcamp.EndDate, &candidate.Bootcamp.Location, &candidate.InstansiPendidikan, &candidate.HackerRank)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		candidates = append(candidates, candidate)
	}

	var totalRows int
	row := c.db.QueryRow("SELECT COUNT(*) FROM candidate")
	err = row.Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	return candidates, common.Paginate(paginationQuery.Page, paginationQuery.Take, totalRows), nil

}

// Constructor
func NewCandidateRepository(db *sql.DB) CandidateRepository {
	return &candidateRepository{db: db}
}
