package manager

import "interview_bootcamp/usecase"

type UseCaseManager interface {
    CandidateUseCase() usecase.CandidateUseCase
    ResumeUseCase() usecase.ResumeUseCase
    InterviewerUseCase() usecase.InterviewerUseCase
	BootcampUseCase() usecase.BootcampUseCase
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

func (u *useCaseManager) InterviewerUseCase() usecase.InterviewerUseCase {
    return usecase.NewInterviewerUseCase(u.repoManager.InterviewerRepo())
}


func NewUseCaseManager(repoManager RepoManager) UseCaseManager {
	return &useCaseManager{repoManager: repoManager}
}
