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
	_, err := r.db.Exec("INSERT INTO users (user_id, user_name, password, role_id) VALUES ($1, $2, $3, $4)",
		user.UserID, user.UserName, user.Password, user.UserRole.RoleID)
	if err != nil {
		return err
	}
	return nil
}

// from here on the password isnt being retrieve
func (r *userRepository) List() ([]model.Users, error) {
	rows, err := r.db.Query(`
		SELECT u.user_id, u.user_name, r.role_id, r.role_name
		FROM users u
		JOIN user_role r ON u.role_id = r.role_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.Users
	for rows.Next() {
		var user model.Users
		err := rows.Scan(&user.UserID, &user.UserName, &user.UserRole.RoleID, &user.UserRole.RoleName)
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
		SELECT u.user_id, u.user_name, r.role_id, r.role_name
		FROM users u
		JOIN user_role r ON u.role_id = r.role_id
		WHERE u.user_id = $1`, id).Scan(
		&user.UserID,
		&user.UserName,
		&user.UserRole.RoleID,
		&user.UserRole.RoleName,
	)
	if err != nil {
		return model.Users{}, err
	}
	return user, nil
}

// get the data based on their user name
func (r *userRepository) GetByUserName(userName string) (model.Users, error) {
	var user model.Users
	err := r.db.QueryRow(`
		SELECT u.user_id, u.user_name, r.role_id, r.role_name
		FROM users u
		JOIN user_role r ON u.role_id = r.role_id
		WHERE u.user_name = $1`, userName).Scan(
		&user.UserID,
		&user.UserName,
		&user.UserRole.RoleID,
		&user.UserRole.RoleName,
	)
	if err != nil {
		return model.Users{}, err
	}
	return user, nil
}

// update all
func (r *userRepository) Update(payload model.Users) error {
	_, err := r.db.Exec(`
		UPDATE users SET user_name = $2, password = $3, role_id = $4
		WHERE user_id = $1`,
		payload.UserID, payload.UserName, payload.Password, payload.UserRole.RoleID)
	if err != nil {
		return err
	}
	return nil
}

// buat yang mau update password aja.
func (r *userRepository) UpdatePassword(userID, newPassword string) error {
	_, err := r.db.Exec("UPDATE users SET password = $2 WHERE user_id = $1", userID, newPassword)
	if err != nil {
		return err
	}
	return nil
}

// delete
func (r *userRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE user_id=$1", id)
	if err != nil {
		return err
	}
	return nil

}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}
