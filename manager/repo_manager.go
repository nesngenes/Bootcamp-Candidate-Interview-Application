package manager

import "interview_bootcamp/repository"

type RepoManager interface {
	// semua repo di daftarkan disini
	CandidateRepo() repository.CandidateRepository
}

type repoManager struct {
	infra InfraManager
}

// UserRepo implements RepoManager.
func (r *repoManager) CandidateRepo() repository.CandidateRepository {
	return repository.NewCandidateRepository(r.infra.Conn())
}

func NewRepoManager(infra InfraManager) RepoManager {
	return &repoManager{infra: infra}
}
