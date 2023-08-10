package manager

import "interview_bootcamp/repository"

type RepoManager interface {
	// semua repo di daftarkan disini
	CandidateRepo() repository.CandidateRepository
	ResumeRepo() repository.ResumeRepository
	UserRolesRepo() repository.UserRolesRepository         //user role
	UsersRepo() repository.UserRepository                  // user
	HRRecruitmentRepo() repository.HRRecruitmentRepository //hr
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

//user role
func (r *repoManager) UserRolesRepo() repository.UserRolesRepository {
	return repository.NewUserRolesRepository(r.infra.Conn())
}

//user
func (r *repoManager) UsersRepo() repository.UserRepository {
	return repository.NewUserRepository(r.infra.Conn())
}

//hr
func (r *repoManager) HRRecruitmentRepo() repository.HRRecruitmentRepository {
	return repository.NewHRRecruitmentRepository(r.infra.Conn())
}

func NewRepoManager(infra InfraManager) RepoManager {
	return &repoManager{infra: infra}
}
