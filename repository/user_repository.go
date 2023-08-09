package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type UserRepository interface {
	BaseRepository[model.Users]
	GetByUserName(userName string) (model.Users, error)
	UpdatePassword(userID, newPassword string) error
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) Create(user model.Users) error {
	_, err := r.db.Exec("INSERT INTO users (id, email, username, password, role_id) VALUES ($1, $2, $3, $4, $5)",
		user.Id, user.Email, user.UserName, user.Password, user.UserRole.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) List() ([]model.Users, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.email, u.username, r.id, r.name as role_name
		FROM users u
		JOIN user_roles r ON u.role_id = r.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.Users
	for rows.Next() {
		var user model.Users
		err := rows.Scan(&user.Id, &user.Email, &user.UserName, &user.UserRole.Id, &user.UserRole.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) Get(id string) (model.Users, error) {
	var user model.Users
	err := r.db.QueryRow(`
		SELECT u.id, u.email, u.username, r.id, r.name
		FROM users u
		JOIN user_roles r ON u.role_id = r.id
		WHERE u.id = $1`, id).Scan(
		&user.Id,
		&user.Email,
		&user.UserName,
		&user.UserRole.Id,
		&user.UserRole.Name,
	)
	if err != nil {
		return model.Users{}, err
	}
	return user, nil
}

func (r *userRepository) GetByUserName(userName string) (model.Users, error) {
	var user model.Users
	err := r.db.QueryRow(`
		SELECT u.id, u.email, u.username, r.id, r.name
		FROM users u
		JOIN user_roles r ON u.role_id = r.id
		WHERE u.username = $1`, userName).Scan(
		&user.Id,
		&user.Email,
		&user.UserName,
		&user.UserRole.Id,
		&user.UserRole.Name,
	)
	if err != nil {
		return model.Users{}, err
	}
	return user, nil
}

func (r *userRepository) Update(payload model.Users) error {
	_, err := r.db.Exec(`
		UPDATE users SET email = $2, username = $3, role_id = $4
		WHERE id = $1`,
		payload.Id, payload.Email, payload.UserName, payload.UserRole.Id)
	if err != nil {
		return err
	}
	return nil
}

// using user_id to update the password only
func (r *userRepository) UpdatePassword(userID, newPassword string) error {
	_, err := r.db.Exec("UPDATE users SET password = $2 WHERE id = $1", userID, newPassword)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}
