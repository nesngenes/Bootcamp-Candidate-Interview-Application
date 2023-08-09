package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type CandidateRepository interface {
	BaseRepository[model.Candidate]
	GetByPhoneNumber(phoneNumber string) (model.Candidate, error)
	GetByEmail(email string) (model.Candidate, error)
}

type candidateRepository struct {
	db *sql.DB
}

func (c *candidateRepository) Create(payload model.Candidate) error {
	_, err := c.db.Exec("INSERT INTO candidate (id, full_name, phone,email, date_of_birth, address,cv_link,bootcamp_id,instansi_pendidikan,hackerrank_score) VALUES ($1, $2, $3, $4, $5, $6, $7,$8,$9,$10)", payload.CandidateID, payload.FullName, payload.Phone, payload.Email, payload.DateOfBirth, payload.Address, payload.CvLink, payload.BootcampId, payload.InstansiPendidikan, payload.HackerRank)
	if err != nil {
		return err
	}
	return nil

}

// GetPhoneNumber implements employeeRepository.
func (c *candidateRepository) GetByPhoneNumber(phoneNumber string) (model.Candidate, error) {
	var candidate model.Candidate
	err := c.db.QueryRow("SELECT * FROM candidate WHERE phone ILIKE $1", "%"+phoneNumber+"%").Scan(&candidate.CandidateID, &candidate.FullName, &candidate.Phone, &candidate.Email, &candidate.DateOfBirth, &candidate.Address, &candidate.CvLink, candidate.BootcampId, candidate.InstansiPendidikan, candidate.HackerRank)
	if err != nil {
		return model.Candidate{}, err
	}
	return candidate, nil
}

func (c *candidateRepository) List() ([]model.Candidate, error) {
	rows, err := c.db.Query("SELECT * FROM candidate")
	if err != nil {
		return nil, err
	}
	var candidates []model.Candidate
	for rows.Next() {
		var candidate model.Candidate
		err := rows.Scan(&candidate.CandidateID, &candidate.FullName, &candidate.Phone, &candidate.Email, &candidate.DateOfBirth, &candidate.Address, &candidate.CvLink, candidate.BootcampId, candidate.InstansiPendidikan, candidate.HackerRank)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}
	return candidates, nil
}

func (c *candidateRepository) Get(id string) (model.Candidate, error) {
	panic("")
}

func (c *candidateRepository) GetByEmail(email string) (model.Candidate, error) {
	var candidate model.Candidate
	err := c.db.QueryRow("SELECT * FROM candidate WHERE email ILIKE $1", "%"+email+"%").Scan(&candidate.CandidateID, &candidate.FullName, &candidate.Phone, &candidate.Email, &candidate.DateOfBirth, &candidate.Address, &candidate.CvLink, candidate.BootcampId, candidate.InstansiPendidikan, candidate.HackerRank)
	if err != nil {
		return model.Candidate{}, err
	}
	return candidate, nil

}

func (c *candidateRepository) Update(payload model.Candidate) error {
	_, err := c.db.Exec("UPDATE product SET full_name=$2, email=$3, date_of_birth=$4, address=$5, cv_link=$6, bootcamp_id=$7, instansi_pendidikan=$8, hackerrank_score=$9 WHERE id = $1", payload.FullName, payload.Email, payload.DateOfBirth, payload.Address, payload.CvLink, payload.BootcampId, payload.InstansiPendidikan, payload.HackerRank)
	if err != nil {
		return err
	}

	return nil
}

func (c *candidateRepository) Delete(id string) error {
	_, err := c.db.Exec("DELETE FROM candidate WHERE candidate_id=$1", id)
	if err != nil {
		return err
	}
	return nil

}

// Constructor
func NewCandidateRepository(db *sql.DB) CandidateRepository {
	return &candidateRepository{db: db}
}
