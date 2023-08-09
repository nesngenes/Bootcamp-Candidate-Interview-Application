package manager

import "interview_bootcamp/usecase"

type UseCaseManager interface {
	CandidateUseCase() usecase.CandidateUseCase
	ResumeUseCase() usecase.ResumeUseCase
	UserRolesUseCase() usecase.UserRolesUseCase // user role
	UsersUseCase() usecase.UserUsecase          //user
}

type useCaseManager struct {
	repoManager RepoManager
}

func (u *useCaseManager) CandidateUseCase() usecase.CandidateUseCase {
	return usecase.NewCandidateUseCase(u.repoManager.CandidateRepo())
}

func (u *useCaseManager) ResumeUseCase() usecase.ResumeUseCase {
	return usecase.NewResumeUseCase(u.repoManager.ResumeRepo())
}

//user role
func (u *useCaseManager) UserRolesUseCase() usecase.UserRolesUseCase {
	return usecase.NewUserRolesUseCase(u.repoManager.UserRolesRepo())
}

//user
func (u *useCaseManager) UsersUseCase() usecase.UserUsecase {
	return usecase.NewUserUsecase(u.repoManager.UsersRepo())
}

func NewUseCaseManager(repoManager RepoManager) UseCaseManager {
	return &useCaseManager{repoManager: repoManager}
}
