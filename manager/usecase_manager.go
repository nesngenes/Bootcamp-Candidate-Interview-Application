package manager

import "interview_bootcamp/usecase"

type UseCaseManager interface {
	CandidateUseCase() usecase.CandidateUseCase
}

type useCaseManager struct {
	repoManager RepoManager
}

// CandidateUseCase implements UseCaseManager.
func (u *useCaseManager) CandidateUseCase() usecase.CandidateUseCase {
	return usecase.NewCandidateUseCase(u.repoManager.CandidateRepo())
}

func NewUseCaseManager(repoManager RepoManager) UseCaseManager {
	return &useCaseManager{repoManager: repoManager}
}
