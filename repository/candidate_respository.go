package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type CandidateRepository interface {
	BaseRepository[model.Candidate]
	GetByPhoneNumber(name string) (model.Candidate, error)
}

type candidateRepository struct {
	db *sql.DB
}

func (c *candidateRepository) Create(payload model.Candidate) error {
	_, err := c.db.Exec("INSERT INTO candidate (candidate_id, first_name, last_name, email, phone, address, date_of_birth) VALUES ($1, $2, $3, $4, $5, $6, $7)", payload.CandidateID, payload.FirstName, payload.LastName, payload.Email, payload.Phone, payload.Address, payload.DateOfBirth)
	if err != nil {
		return err
	}
	return nil

}

// GetPhoneNumber implements employeeRepository.
func (c *candidateRepository) GetByPhoneNumber(phoneNumber string) (model.Candidate, error) {
	var candidate model.Candidate
	err := c.db.QueryRow("SELECT candidate_id, first_name, last_name, email, phone, address, date_of_birth FROM candidate WHERE phone=$1", phoneNumber).Scan(&candidate.CandidateID, &candidate.FirstName, &candidate.LastName, &candidate.Phone, &candidate.Address, &candidate.DateOfBirth)
	if err != nil {
		return model.Candidate{}, err
	}
	return candidate, nil
}

func (c *candidateRepository) List() ([]model.Candidate, error) {
	panic("")
}

func (c *candidateRepository) Get(id string) (model.Candidate, error) {
	panic("")

}

func (c *candidateRepository) GetByName(name string) (model.Candidate, error) {
	panic("")

}

func (c *candidateRepository) Update(payload model.Candidate) error {
	panic("")

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
