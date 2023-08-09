package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type ResumeRepository interface {
	BaseRepository[model.Resume]
	Create(payload model.Resume) error
	Update(payload model.Resume) error
	GetByName(name string) (model.Resume, error)
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
	_, err := r.db.Exec("UPDATE product SET full_name=$2, email=$3, date_of_birth=$4, address=$5, cv_link=$6, bootcamp_id=$7, instansi_pendidikan=$8, hackerrank_score=$9 WHERE id = $1", payload.FullName, payload.Email, payload.DateOfBirth, payload.Address, payload.CvLink, payload.BootcampId, payload.InstansiPendidikan, payload.HackerRank)
	if err != nil {
		return err
	}

	return nil

}

func (r *resumeRepository) Delete(id string) error {
	panic("")

}

// Constructor
func NewResumeRepository(db *sql.DB) ResumeRepository {
	return &resumeRepository{db: db}
}
