package manager

import "interview_bootcamp/usecase"

type UseCaseManager interface {
	CandidateUseCase() usecase.CandidateUseCase
	BootcampUseCase() usecase.BootcampUseCase
	ResumeUseCase() usecase.ResumeUseCase
}

type useCaseManager struct {
	repoManager RepoManager
}

func (u *useCaseManager) CandidateUseCase() usecase.CandidateUseCase {
	return usecase.NewCandidateUseCase(u.repoManager.CandidateRepo())
}
func (u *useCaseManager) BootcampUseCase() usecase.BootcampUseCase {
	return usecase.NewBootcampUseCase(u.repoManager.BootcampRepo())
}

func (u *useCaseManager) ResumeUseCase() usecase.ResumeUseCase {
	return usecase.NewResumeUseCase(u.repoManager.ResumeRepo())
}

func NewUseCaseManager(repoManager RepoManager) UseCaseManager {
	return &useCaseManager{repoManager: repoManager}
}
