package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type ResumeRepository interface {
	BaseRepository[model.Resume]
}

type resumeRepository struct {
	db *sql.DB
}

func (r *resumeRepository) Create(payload model.Resume) error {
	_, err := r.db.Exec("INSERT INTO resume (resume_id, candidate_id, cv_file) VALUES ($1, $2, $3)", payload.ResumeID, payload.CandidateID, payload.CvURL)
	if err != nil {
		return err
	}
	return nil

}


func (r *resumeRepository) List() ([]model.Resume, error) {
	panic("")
}

func (r *resumeRepository) Get(id string) (model.Resume, error) {
	panic("")

}

func (r *resumeRepository) GetByName(name string) (model.Resume, error) {
	panic("")

}

func (r *resumeRepository) Update(payload model.Resume) error {
	panic("")
	
}

func (r *resumeRepository) Delete(id string) error {
	panic("")

}

// Constructor
func NewResumeRepository(db *sql.DB) ResumeRepository {
	return &resumeRepository{db: db}
}
