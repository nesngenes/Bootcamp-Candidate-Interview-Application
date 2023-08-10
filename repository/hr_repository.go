package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type HRRecruitmentRepository interface {
	BaseRepository[model.HRRecruitment]
	GetByUserID(userID string) (model.HRRecruitment, error)
}

type hrRecruitmentRepository struct {
	db *sql.DB
}

func (r *hrRecruitmentRepository) Create(hr model.HRRecruitment) error {
	_, err := r.db.Exec("INSERT INTO hr_recruitment (id, full_name, user_id) VALUES ($1, $2, $3)",
		hr.ID, hr.FullName, hr.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *hrRecruitmentRepository) Get(id string) (model.HRRecruitment, error) {
	var hrRecruitment model.HRRecruitment
	err := r.db.QueryRow(`
		SELECT hr.id, hr.full_name, hr.user_id, u.id, u.email, u.username, ur.id, ur.name
		FROM hr_recruitment hr
		JOIN users u ON hr.user_id = u.id
		JOIN user_roles ur ON u.role_id = ur.id
		WHERE hr.id = $1`, id).Scan(
		&hrRecruitment.ID,
		&hrRecruitment.FullName,
		&hrRecruitment.UserID,
		&hrRecruitment.User.Id,
		&hrRecruitment.User.Email,
		&hrRecruitment.User.UserName,
		&hrRecruitment.User.UserRole.Id,
		&hrRecruitment.User.UserRole.Name,
	)
	if err != nil {
		return model.HRRecruitment{}, err
	}
	return hrRecruitment, nil
}

func (r *hrRecruitmentRepository) GetByUserID(userID string) (model.HRRecruitment, error) {
	var hrRecruitment model.HRRecruitment
	err := r.db.QueryRow(`
		SELECT hr.id, hr.full_name, hr.user_id, u.id, u.email, u.username, ur.id, ur.name
		FROM hr_recruitment hr
		JOIN users u ON hr.user_id = u.id
		JOIN user_roles ur ON u.role_id = ur.id
		WHERE hr.user_id = $1`, userID).Scan(
		&hrRecruitment.ID,
		&hrRecruitment.FullName,
		&hrRecruitment.UserID,
		&hrRecruitment.User.Id,
		&hrRecruitment.User.Email,
		&hrRecruitment.User.UserName,
		&hrRecruitment.User.UserRole.Id,
		&hrRecruitment.User.UserRole.Name,
	)
	if err != nil {
		return model.HRRecruitment{}, err
	}
	return hrRecruitment, nil
}

func (r *hrRecruitmentRepository) List() ([]model.HRRecruitment, error) {
	rows, err := r.db.Query(`
		SELECT hr.id, hr.full_name, hr.user_id, u.id, u.email, u.username, ur.id, ur.name
		FROM hr_recruitment hr
		JOIN users u ON hr.user_id = u.id
		JOIN user_roles ur ON u.role_id = ur.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hrRecruitments []model.HRRecruitment
	for rows.Next() {
		var hrRecruitment model.HRRecruitment
		var user model.Users
		var userRole model.UserRoles
		err := rows.Scan(
			&hrRecruitment.ID,
			&hrRecruitment.FullName,
			&hrRecruitment.UserID,
			&user.Id,
			&user.Email,
			&user.UserName,
			&userRole.Id,
			&userRole.Name,
		)
		if err != nil {
			return nil, err
		}
		user.UserRole = userRole
		hrRecruitment.User = user
		hrRecruitments = append(hrRecruitments, hrRecruitment)
	}
	return hrRecruitments, nil
}

func (r *hrRecruitmentRepository) Update(hr model.HRRecruitment) error {
	_, err := r.db.Exec(`
		UPDATE hr_recruitment SET full_name = $2, user_id = $3
		WHERE id = $1`,
		hr.ID, hr.FullName, hr.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *hrRecruitmentRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM hr_recruitment WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func NewHRRecruitmentRepository(db *sql.DB) HRRecruitmentRepository {
	return &hrRecruitmentRepository{db}
}
