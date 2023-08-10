package usecase

import (
	"fmt"
	"interview_bootcamp/utils/security"
)

type AuthUseCase interface {
	Login(username string, password string) (string, error)
}

type authUseCase struct {
	usecase UserUsecase
}

func (a *authUseCase) Login(username string, password string) (string, error) {
	user, err := a.usecase.FindByUsernamePassword(username, password)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	token, err := security.CreateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}
	return token, nil
}

func NewAuthUseCase(usecase UserUsecase) AuthUseCase {
	return &authUseCase{usecase: usecase}
}
