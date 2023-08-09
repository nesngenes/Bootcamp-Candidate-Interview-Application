package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/repository"
)

type UserUsecase interface {
	RegisterNewUser(payload model.Users) error
	List() ([]model.Users, error)
	GetUserByEmail(email string) (model.Users, error)
	GetUserByUserName(userName string) (model.Users, error)
	GetUserByID(id string) (model.Users, error)
	UpdateUser(payload model.Users) error
	DeleteUser(id string) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func (u *userUsecase) RegisterNewUser(payload model.Users) error {
	// cek klo kosong
	if payload.Email == "" || payload.UserName == "" || payload.Password == "" {
		return fmt.Errorf("email, username, and password are required fields")
	}

	// cek kali username ada yang sama
	existingUserByUsername, err := u.userRepo.GetByUserName(payload.UserName)
	if err == nil && existingUserByUsername.UserName == payload.UserName {
		return fmt.Errorf("user with username %s already exists", payload.UserName)
	}

	// Check if a user with the given email already exists
	existingUserByEmail, err := u.userRepo.GetByEmail(payload.Email)
	if err == nil && existingUserByEmail.Email == payload.Email {
		return fmt.Errorf("user with email %s already exists", payload.Email)
	}

	// Create the new user
	err = u.userRepo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to register new user: %v", err)
	}

	return nil
}

func (u *userUsecase) List() ([]model.Users, error) {
	return u.userRepo.List()
}

func (u *userUsecase) GetUserByEmail(email string) (model.Users, error) {
	return u.userRepo.GetByEmail(email)
}

func (u *userUsecase) GetUserByUserName(userName string) (model.Users, error) {
	return u.userRepo.GetByUserName(userName)
}

func (u *userUsecase) GetUserByID(id string) (model.Users, error) {
	return u.userRepo.Get(id)
}

func (u *userUsecase) UpdateUser(payload model.Users) error {
	// Check if the user exists
	existingUser, err := u.userRepo.Get(payload.Id)
	if err != nil {
		return fmt.Errorf("error checking for existing user: %v", err)
	}

	//Check if the updated username or email conflicts with existing users
	if payload.UserName != existingUser.UserName {
		existingUserByUsername, err := u.userRepo.GetByUserName(payload.UserName)
		if err == nil && existingUserByUsername.Id != payload.Id {
			return fmt.Errorf("user with username %s already exists", payload.UserName)
		}
	}

	if payload.Email != existingUser.Email {
		existingUserByEmail, err := u.userRepo.GetByEmail(payload.Email)
		if err == nil && existingUserByEmail.Id != payload.Id {
			return fmt.Errorf("user with email %s already exists", payload.Email)
		}
	}

	// Update the user
	err = u.userRepo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

func (u *userUsecase) DeleteUser(id string) error {
	return u.userRepo.Delete(id)
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}
