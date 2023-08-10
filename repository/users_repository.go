package repository

import (
	"database/sql"
	"fmt"
	"interview_bootcamp/model"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	BaseRepository[model.Users]
	GetByEmail(email string) (model.Users, error)
	GetByUserName(username string) (model.Users, error)
	GetUsernamePassword(username string, password string) (model.Users, error)
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

func (r *userRepository) GetByEmail(email string) (model.Users, error) {
	var user model.Users
	err := r.db.QueryRow(`
		SELECT u.id, u.email, u.username, u.password, r.id, r.name
		FROM users u
		JOIN user_roles r ON u.role_id = r.id
		WHERE u.email ILIKE $1`, "%"+email+"%").Scan(
		&user.Id,
		&user.Email,
		&user.UserName,
		&user.Password,
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
		SELECT u.id, u.email, u.username,u.password, r.id, r.name
		FROM users u
		JOIN user_roles r ON u.role_id = r.id
		WHERE u.username ILIKE $1`, "%"+userName+"%").Scan(
		&user.Id,
		&user.Email,
		&user.UserName,
		&user.Password,
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

func (r *userRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) GetUsernamePassword(username string, password string) (model.Users, error) {

	user, err := u.GetByUserName(username)
	if err != nil {
		return model.Users{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.Users{}, fmt.Errorf("failed to verivy password hash : %v", err)
	}

	return user, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}
