package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type InterviewerRepository interface {
	BaseRepository[model.Interviewer]
}

type interviewerRepository struct {
	db *sql.DB
}

func (i *interviewerRepository) Create(payload model.Interviewer) error {
	_, err := i.db.Exec("INSERT INTO interviewer (id, full_name, user_id) VALUES ($1, $2, $3)", payload.InterviewerID, payload.FullName, payload.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (i *interviewerRepository) List() ([]model.Interviewer, error) {
	rows, err := i.db.Query(`
		SELECT i.id, i.full_name, u.id
		FROM interviewer i
		INNER JOIN users u ON u.id = i.user_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interviewers []model.Interviewer
	for rows.Next() {
		var interviewer model.Interviewer
		err := rows.Scan(
			&interviewer.InterviewerID,
			&interviewer.FullName,
			&interviewer.UserID,
		)
		if err != nil {
			return nil, err
		}
		interviewers = append(interviewers, interviewer)
	}
	return interviewers, nil
}

func (i *interviewerRepository) Get(id string) (model.Interviewer, error) {
	var interviewer model.Interviewer
	err := i.db.QueryRow("SELECT i.id, i.full_name, u.id FROM interviewer i INNER JOIN users u ON u.id = i.user_id WHERE i.id = $1", id).Scan(
		&interviewer.InterviewerID,
		&interviewer.FullName,
		&interviewer.UserID,
	)
	if err != nil {
		return model.Interviewer{}, err
	}
	return interviewer, nil

}

func (i *interviewerRepository) Update(payload model.Interviewer) error {
	_, err := i.db.Exec("UPDATE interviewer SET full_name = $2, user_id = $3 WHERE id = $1", payload.InterviewerID, payload.FullName, payload.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (i *interviewerRepository) Delete(id string) error {
	_, err := i.db.Exec("DELETE FROM interviewer WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil

}

// Constructor
func NewInterviewerRepository(db *sql.DB) InterviewerRepository {
	return &interviewerRepository{db: db}
}
