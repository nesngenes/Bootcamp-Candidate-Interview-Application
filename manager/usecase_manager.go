package manager

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"interview_bootcamp/usecase"
)

type UseCaseManager interface {
    CandidateUseCase() usecase.CandidateUseCase
    ResumeUseCase() usecase.ResumeUseCase
    InterviewerUseCase() usecase.InterviewerUseCase
	BootcampUseCase() usecase.BootcampUseCase
	SetCloudinaryInstance(cloudinary *cloudinary.Cloudinary)

}

type useCaseManager struct {
	repoManager RepoManager
	cloudinary  *cloudinary.Cloudinary
}

func (u *useCaseManager) CandidateUseCase() usecase.CandidateUseCase {
	return usecase.NewCandidateUseCase(u.repoManager.CandidateRepo())
}
func (u *useCaseManager) BootcampUseCase() usecase.BootcampUseCase {
	return usecase.NewBootcampUseCase(u.repoManager.BootcampRepo())
}

func (u *useCaseManager) ResumeUseCase() usecase.ResumeUseCase {
    return usecase.NewResumeUseCase(u.repoManager.ResumeRepo(), u.repoManager.CloudinaryInstance())
}

func (u *useCaseManager) InterviewerUseCase() usecase.InterviewerUseCase {
    return usecase.NewInterviewerUseCase(u.repoManager.InterviewerRepo())
}

func (u *useCaseManager) SetCloudinaryInstance(cloudinary *cloudinary.Cloudinary) {
    u.cloudinary = cloudinary
}


func NewUseCaseManager(repoManager RepoManager) UseCaseManager {
	return &useCaseManager{repoManager: repoManager}
}
