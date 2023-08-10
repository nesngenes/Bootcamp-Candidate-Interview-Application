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
	roleUC  UserRolesUseCase
}

func (a *authUseCase) Login(username string, password string) (string, error) {
	user, err := a.usecase.FindByUsernamePassword(username, password)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	role, _ := a.roleUC.GetUserRoleByID(user.UserRole.Id)
	user.UserRole = role

	token, err := security.CreateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}
	return token, nil
}

func NewAuthUseCase(usecase UserUsecase, roleUC UserRolesUseCase) AuthUseCase {
	return &authUseCase{usecase: usecase, roleUC: roleUC}
}
