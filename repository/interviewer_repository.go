package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type InterviewerRepository interface {
	BaseRepository[model.Interviewer]
	GetByPhoneNumber(name string) (model.Interviewer, error)
	GetByEmail(email string) (model.Interviewer, error)
}

type interviewerRepository struct {
	db *sql.DB
}

func (c *interviewerRepository) Create(payload model.Interviewer) error {
	_, err := c.db.Exec("INSERT INTO interviewer (interviewer_id, user_id, first_name, last_name, email, phone, specialization) VALUES ($1, $2, $3, $4, $5, $6, $7)", payload.InterviewerID, payload.UserID, payload.FirstName, payload.LastName, payload.Email, payload.Phone, payload.Specialization)
	if err != nil {
		return err
	}
	return nil

}

// GetPhoneNumber implements employeeRepository.
func (c *interviewerRepository) GetByPhoneNumber(phoneNumber string) (model.Interviewer, error) {
	var interviewer model.Interviewer
	err := c.db.QueryRow(`
		SELECT i.interviewer_id, i.user_id, u.username, i.first_name, i.last_name, i.email, i.phone, i.specialization, r.role_name
		FROM interviewer i
		JOIN users u ON i.user_id = u.user_id
		JOIN user_role r ON u.role_id = r.role_id
		WHERE i.phone ILIKE $1`, "%"+phoneNumber+"%").Scan(
		&interviewer.InterviewerID,
		&interviewer.UserID,
		&interviewer.UserName,
		&interviewer.FirstName,
		&interviewer.LastName,
		&interviewer.Email,
		&interviewer.Phone,
		&interviewer.Specialization,
		&interviewer.RoleName,
	)
	if err != nil {
		return model.Interviewer{}, err
	}
	return interviewer, nil
}

// get by email
func (c *interviewerRepository) GetByEmail(email string) (model.Interviewer, error) {
	var interviewer model.Interviewer
	err := c.db.QueryRow(`
	SELECT i.interviewer_id, i.user_id, u.username, i.first_name, i.last_name, i.email, i.phone, i.specialization, r.role_name
	FROM interviewer i
	JOIN users u ON i.user_id = u.user_id
	JOIN user_role r ON u.role_id = r.role_id
	WHERE i.email ILIKE $1`, "%"+email+"%").Scan(
		&interviewer.InterviewerID,
		&interviewer.UserID,
		&interviewer.UserName,
		&interviewer.FirstName,
		&interviewer.LastName,
		&interviewer.Email,
		&interviewer.Phone,
		&interviewer.Specialization,
		&interviewer.RoleName,
	)
	if err != nil {
		return model.Interviewer{}, err
	}
	return interviewer, nil

}

func (c *interviewerRepository) List() ([]model.Interviewer, error) {
	rows, err := c.db.Query(`
		SELECT i.interviewer_id, u.user_id, u.username, i.first_name, i.last_name, i.email, i.phone, i.specialization, r.role_name 
		FROM interviewer i
		JOIN users u ON i.user_id = u.user_id
		JOIN user_role r ON u.role_id = r.role_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interviewers []model.Interviewer
	for rows.Next() {
		var interviewer model.Interviewer
		err := rows.Scan(
			&interviewer.InterviewerID,
			&interviewer.UserID,
			&interviewer.UserName,
			&interviewer.FirstName,
			&interviewer.LastName,
			&interviewer.Email,
			&interviewer.Phone,
			&interviewer.Specialization,
			&interviewer.RoleName,
		)
		if err != nil {
			return nil, err
		}
		interviewers = append(interviewers, interviewer)
	}
	return interviewers, nil
}

func (c *interviewerRepository) Get(id string) (model.Interviewer, error) {
	var interviewer model.Interviewer
	err := c.db.QueryRow(`
	SELECT i.interviewer_id, i.user_id, u.username, i.first_name, i.last_name, i.email, i.phone, i.specialization, r.role_name
	FROM interviewer i
	JOIN users u ON i.user_id = u.user_id
	JOIN user_role r ON u.role_id = r.role_id
	WHERE i.interviewer_id = $1`, id).Scan(
		&interviewer.InterviewerID,
		&interviewer.UserID,
		&interviewer.UserName,
		&interviewer.FirstName,
		&interviewer.LastName,
		&interviewer.Email,
		&interviewer.Phone,
		&interviewer.Specialization,
		&interviewer.RoleName,
	)
	if err != nil {
		return model.Interviewer{}, err
	}
	return interviewer, nil

}

func (c *interviewerRepository) Update(payload model.Interviewer) error {
	_, err := c.db.Exec("UPDATE interviewer SET user_id = $2, first_name = $3, last_name = $4, email = $5, phone = $6, specialization = $7  WHERE id = $1", payload.InterviewerID, payload.UserID, payload.FirstName, payload.LastName, payload.Email, payload.Phone, payload.Specialization)
	if err != nil {
		return err
	}
	return nil

}

func (c *interviewerRepository) Delete(id string) error {
	_, err := c.db.Exec("DELETE FROM interviewer WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil

}

// Constructor
func NewInterviewerRepository(db *sql.DB) InterviewerRepository {
	return &interviewerRepository{db: db}
}
