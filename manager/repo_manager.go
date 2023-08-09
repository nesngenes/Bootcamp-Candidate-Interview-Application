package manager

import "interview_bootcamp/repository"

type RepoManager interface {
	// semua repo di daftarkan disini
	CandidateRepo() repository.CandidateRepository
	ResumeRepo() repository.ResumeRepository
	InterviewerRepo() repository.InterviewerRepository
}

type repoManager struct {
	infra InfraManager
}

// UserRepo implements RepoManager.
func (r *repoManager) CandidateRepo() repository.CandidateRepository {
	return repository.NewCandidateRepository(r.infra.Conn())
}

func (r *repoManager) ResumeRepo() repository.ResumeRepository {
	return repository.NewResumeRepository(r.infra.Conn())
}

func (r *repoManager) InterviewerRepo() repository.InterviewerRepository {
	return repository.NewInterviewerRepository(r.infra.Conn())
}

func NewRepoManager(infra InfraManager) RepoManager {
	return &repoManager{infra: infra}
}
