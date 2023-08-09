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
	ResumeRepo() repository.ResumeRepository
	InterviewerRepo() repository.InterviewerRepository
	CloudinaryInstance() *cloudinary.Cloudinary
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

func (r *repoManager) ResumeRepo() repository.ResumeRepository {
	return repository.NewResumeRepository(r.infra.Conn())
}

func (r *repoManager) InterviewerRepo() repository.InterviewerRepository {
	return repository.NewInterviewerRepository(r.infra.Conn())
}

func (r *repoManager) CloudinaryInstance() *cloudinary.Cloudinary {
    return r.cloudinary
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
