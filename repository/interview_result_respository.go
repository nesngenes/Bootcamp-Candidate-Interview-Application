package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type InterviewResultRepository interface {
	BaseRepository[model.InterviewResult]
	GetByIdInterviewResult(result string) (model.InterviewResult, error)
	CreateInterviewResult(payload model.InterviewResult) error
	ListInterviewResult() ([]model.InterviewResult, error)
	UpdateInterviewResult(payload model.InterviewResult) error
	DeleteInterviewResult(id string) error
}

type interviewResultRepository struct {
	db *sql.DB
}

// membuat data hasil interview pertama kali
func (cr *interviewResultRepository) CreateInterviewResult(payload model.InterviewResult) error {
	_, err := cr.db.Exec("INSERT INTO interview_result (id, interview_id, result_id, note) VALUES ($1, $2, $3, $4)", payload.Id, payload.InterviewId, payload.ResultId, payload.Note)
	if err != nil {
		return err
	}
	return nil
}

// mendapatkan data hasil interview berdasarkan id
func (cr *interviewResultRepository) GetByIdInterviewResult(id string) (model.InterviewResult, error) {
	var interviewResult model.InterviewResult
	err := cr.db.QueryRow("SELECT * FROM interview_result WHERE id=$1", interviewResult.Id)
	if err != nil {
		return model.InterviewResult{}, nil
	}
	return interviewResult, nil
}

// mendapatkan seluruh hasil interview
func (cr *interviewResultRepository) ListInterviewResult() ([]model.InterviewResult, error) {
	rows, err := cr.db.Query("SELECT * FROM interview_result")
	if err != nil {
		return nil, err
	}

	var interviewResults []model.InterviewResult
	for rows.Next() {
		var interviewResult model.InterviewResult
		err := rows.Scan(&interviewResult.Id, &interviewResult.InterviewId, &interviewResult.ResultId, &interviewResult.Note)
		if err != nil {
			return nil, err
		}
		interviewResults = append(interviewResults, interviewResult)
	}

	return interviewResults, nil
}

// melakukan perubahan pada data hasil interview
func (cr *interviewResultRepository) UpdateInterviewResult(payload model.InterviewResult) error {
	_, err := cr.db.Exec("UPDATE product SET note=$2 WHERE id = $1", payload.Id, payload.Note)
	if err != nil {
		return err
	}

	return nil
}

// melakukan hapus data pada hasil review
func (cr *interviewResultRepository) DeleteInterviewResult(id string) error {
	_, err := cr.db.Exec("DELETE FROM interview_result WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil

}

// Constructor
func NewInterviewResultRepository(db *sql.DB) InterviewerRepository {
	return &interviewerRepository{db: db}
}
