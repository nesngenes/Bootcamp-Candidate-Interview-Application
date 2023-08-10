package usecase

import (
	"fmt"
	"interview_bootcamp/model"
)

type AuthUseCAse interface {
	Login(username string) (string, error)
}

type authUseCase struct {
	user UserUsecase
}

func (a *authUseCase) Login(username string) (string, error) {
	user, err := a.user.GetUserByUserName(username)
	if err != nil {
		return "un authorization", fmt.Errorf("user un authorization")
	}

	//bila user ditemukan maka token akan muncul
	token, err := GenerateToken(user)
	if err != nil {
		return "gagal", fmt.Errorf("token gagal dibuat")
	}

	return token, nil
}

// untuk membuat token sebagai autentikasi
func GenerateToken(user model.Users) (string, error) {
	return "iniToken", nil
}

func NewAuthUseCase(user AuthUseCAse) AuthUseCAse {
	return &authUseCase{user: user}
}
