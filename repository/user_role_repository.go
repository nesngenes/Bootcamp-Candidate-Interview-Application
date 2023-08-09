package repository

import (
	"database/sql"
	"interview_bootcamp/model"
)

type UserRolesRepository interface {
	BaseRepository[model.UserRoles]
	GetByName(name string) (model.UserRoles, error)
}

type userRolesRepository struct {
	db *sql.DB
}

func (r *userRolesRepository) Create(userRole model.UserRoles) error {
	_, err := r.db.Exec("INSERT INTO user_roles (id, name) VALUES ($1, $2)",
		userRole.Id, userRole.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRolesRepository) List() ([]model.UserRoles, error) {
	rows, err := r.db.Query("SELECT id, name FROM user_roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRoles []model.UserRoles
	for rows.Next() {
		var userRole model.UserRoles
		err := rows.Scan(&userRole.Id, &userRole.Name)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, userRole)
	}
	return userRoles, nil
}

// get by id
func (r *userRolesRepository) Get(id string) (model.UserRoles, error) {
	var userRole model.UserRoles
	err := r.db.QueryRow("SELECT id, name FROM user_roles WHERE id = $1", id).Scan(&userRole.Id, &userRole.Name)
	if err != nil {
		return model.UserRoles{}, err
	}
	return userRole, nil

}

// get role by name
func (r *userRolesRepository) GetByName(name string) (model.UserRoles, error) {
	var userRole model.UserRoles
	// Use ILIKE for case-insensitive search
	err := r.db.QueryRow("SELECT id, name FROM user_roles WHERE name ILIKE $1", "%"+name+"%").Scan(&userRole.Id, &userRole.Name)
	if err != nil {
		return model.UserRoles{}, err
	}
	return userRole, nil
}

func (r *userRolesRepository) Update(userRole model.UserRoles) error {
	_, err := r.db.Exec("UPDATE user_roles SET name = $2 WHERE id = $1",
		userRole.Id, userRole.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRolesRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM user_roles WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil

}

func NewUserRolesRepository(db *sql.DB) UserRolesRepository {
	return &userRolesRepository{db}
}
