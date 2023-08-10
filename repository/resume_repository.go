// repository/resume_repository.go
package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type ResumeRepository interface {
	Create(payload model.Resume) error
	Get(id string) (model.Resume, error)
	Delete(id string) error
	Update(payload model.Resume) error
}

type resumeRepository struct {
	db *sql.DB
}

func (r *resumeRepository) Create(payload model.Resume) error {
	_, err := r.db.Exec(
		"INSERT INTO resume (resume_id, candidate_id, cv_url) VALUES ($1, $2, $3)",
		payload.ResumeID, payload.CandidateID, payload.CvURL,
	)
	return err
}

func (r *resumeRepository) Get(id string) (model.Resume, error) {
	var resume model.Resume
	err := r.db.QueryRow(
		"SELECT resume_id, candidate_id, cv_url FROM resume WHERE resume_id = $1",
		id,
	).Scan(&resume.ResumeID, &resume.CandidateID, &resume.CvURL)
	return resume, err
}

func (r *resumeRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM resume WHERE resume_id = $1", id)
	return err
}

func (r *resumeRepository) Update(payload model.Resume) error {
	_, err := r.db.Exec(
		"UPDATE resume SET candidate_id = $2, cv_url = $3 WHERE resume_id = $1",
		payload.ResumeID, payload.CandidateID, payload.CvURL,
	)
	return err
}

func NewResumeRepository(db *sql.DB) ResumeRepository {
	return &resumeRepository{db: db}
}
