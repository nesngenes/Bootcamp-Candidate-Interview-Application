package manager


import (
    "github.com/cloudinary/cloudinary-go/v2"
	"interview_bootcamp/repository"
	"interview_bootcamp/config"
)

type RepoManager interface {
	// semua repo di daftarkan disini
	CandidateRepo() repository.CandidateRepository
	BootcampRepo() repository.BootcampRepository
	StatusRepo() repository.StatusRepository
	ResumeRepo() repository.ResumeRepository
	InterviewerRepo() repository.InterviewerRepository
	CloudinaryInstance() *cloudinary.Cloudinary
	UserRolesRepo() repository.UserRolesRepository //user role
	UsersRepo() repository.UserRepository          // user
}

type repoManager struct {
	infra InfraManager
	cloudinary   *cloudinary.Cloudinary
}

// UserRepo implements RepoManager.
func (r *repoManager) CandidateRepo() repository.CandidateRepository {
	return repository.NewCandidateRepository(r.infra.Conn())
}
func (r *repoManager) BootcampRepo() repository.BootcampRepository {
	return repository.NewBootcampRepository(r.infra.Conn())
}
func (r *repoManager) StatusRepo() repository.StatusRepository {
	return repository.NewStatusRepository(r.infra.Conn())
}

func (r *repoManager) ResumeRepo() repository.ResumeRepository {
	return repository.NewResumeRepository(r.infra.Conn())
}

func (r *repoManager) InterviewerRepo() repository.InterviewerRepository {
	return repository.NewInterviewerRepository(r.infra.Conn())
}

func (r *repoManager) CloudinaryInstance() *cloudinary.Cloudinary {
    return r.cloudinary
}
//user role
func (r *repoManager) UserRolesRepo() repository.UserRolesRepository {
	return repository.NewUserRolesRepository(r.infra.Conn())
}

//user
func (r *repoManager) UsersRepo() repository.UserRepository {
	return repository.NewUserRepository(r.infra.Conn())
}

func NewRepoManager(infra InfraManager) RepoManager {
	cfg, err := config.NewConfig()
	if err != nil {
		panic("Failed to read configuration")
	}

	cloudinaryInstance, err := cloudinary.NewFromParams(
        cfg.CloudinaryConfig.CloudinaryCloudName,
        cfg.CloudinaryConfig.CloudinaryAPIKey,
        cfg.CloudinaryConfig.CloudinaryAPISecret,
    )
    if err != nil {
        panic("Failed to initialize Cloudinary")
    }
	return &repoManager{
        infra: infra,
        cloudinary:   cloudinaryInstance,
    }
}
