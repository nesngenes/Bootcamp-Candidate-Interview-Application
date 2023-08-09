package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/repository"
)

type UserRolesUseCase interface {
	RegisterNewUserRole(payload model.UserRoles) error
	GetAllUserRoles() ([]model.UserRoles, error)
	GetUserRoleByID(id string) (model.UserRoles, error)
	UpdateUserRole(payload model.UserRoles) error
	DeleteUserRole(id string) error
}

type userRolesUseCase struct {
	repo repository.UserRolesRepository
}

func (u *userRolesUseCase) RegisterNewUserRole(payload model.UserRoles) error {
	// check for empty field
	if payload.Name == "" {
		return fmt.Errorf("name is a required field")
	}

	// Check if user role with the same name already exists
	isExistUserRole, _ := u.repo.GetByName(payload.Name)
	if isExistUserRole.Name == payload.Name {
		return fmt.Errorf("user role with name %s already exists", payload.Name)
	}

	err := u.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to register new user role: %v", err)
	}
	return nil
}

func (u *userRolesUseCase) GetAllUserRoles() ([]model.UserRoles, error) {
	return u.repo.List()
}

func (u *userRolesUseCase) GetUserRoleByID(id string) (model.UserRoles, error) {
	return u.repo.Get(id)
}

func (u *userRolesUseCase) UpdateUserRole(payload model.UserRoles) error {
	if payload.Name == "" {
		return fmt.Errorf("name is required field")
	}

	isExistRole, _ := u.repo.GetByName(payload.Name)
	if isExistRole.Name == payload.Name && isExistRole.Id != payload.Id {
		return fmt.Errorf("user role with name %s exists", payload.Name)
	}

	err := u.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update user role: %v", err)
	}

	return nil

}

func (u *userRolesUseCase) DeleteUserRole(id string) error {
	user, err := u.GetUserRoleByID(id)
	if err != nil {
		return fmt.Errorf("user role with ID %s not found", id)
	}

	err = u.repo.Delete(user.Id)
	if err != nil {
		return fmt.Errorf("failed to delete user role: %v", err)
	}
	return nil

}

func NewUserRolesUseCase(repo repository.UserRolesRepository) UserRolesUseCase {
	return &userRolesUseCase{repo: repo}
}
