package manager

import (
	"interview_bootcamp/usecase"

	"github.com/cloudinary/cloudinary-go/v2"
)

type UseCaseManager interface {
	CandidateUseCase() usecase.CandidateUseCase
	ResumeUseCase() usecase.ResumeUseCase
	UserRolesUseCase() usecase.UserRolesUseCase         // user role
	UsersUseCase() usecase.UserUsecase                  //user
	HRRecruitmentUsecase() usecase.HRRecruitmentUsecase //hr
	InterviewerUseCase() usecase.InterviewerUseCase
	BootcampUseCase() usecase.BootcampUseCase
	SetCloudinaryInstance(cloudinary *cloudinary.Cloudinary)
	StatusUseCase() usecase.StatusUseCase
}

type useCaseManager struct {
	repoManager RepoManager
	cloudinary  *cloudinary.Cloudinary
}

func (u *useCaseManager) CandidateUseCase() usecase.CandidateUseCase {
	return usecase.NewCandidateUseCase(u.repoManager.CandidateRepo(), u.BootcampUseCase(), u.repoManager.CloudinaryInstance())
}
func (u *useCaseManager) BootcampUseCase() usecase.BootcampUseCase {
	return usecase.NewBootcampUseCase(u.repoManager.BootcampRepo())
}
func (u *useCaseManager) StatusUseCase() usecase.StatusUseCase {
	return usecase.NewStatusUseCase(u.repoManager.StatusRepo())
}

func (u *useCaseManager) ResumeUseCase() usecase.ResumeUseCase {
	return usecase.NewResumeUseCase(u.repoManager.ResumeRepo(), u.repoManager.CloudinaryInstance())
}

func (u *useCaseManager) InterviewerUseCase() usecase.InterviewerUseCase {
	return usecase.NewInterviewerUseCase(u.repoManager.InterviewerRepo())
}

// user role
func (u *useCaseManager) UserRolesUseCase() usecase.UserRolesUseCase {
	return usecase.NewUserRolesUseCase(u.repoManager.UserRolesRepo())
}

// user
func (u *useCaseManager) UsersUseCase() usecase.UserUsecase {
	return usecase.NewUserUsecase(u.repoManager.UsersRepo())
}

// hr
func (u *useCaseManager) HRRecruitmentUsecase() usecase.HRRecruitmentUsecase {
	return usecase.NewHRRecruitmentUsecase(u.repoManager.HRRecruitmentRepo(), u.repoManager.UsersRepo())
}

func (u *useCaseManager) SetCloudinaryInstance(cloudinary *cloudinary.Cloudinary) {
	u.cloudinary = cloudinary
}

func NewUseCaseManager(repoManager RepoManager) UseCaseManager {
	return &useCaseManager{repoManager: repoManager}
}
