package manager

import (
	"interview_bootcamp/usecase"

	"github.com/cloudinary/cloudinary-go/v2"
)

type UseCaseManager interface {
	CandidateUseCase() usecase.CandidateUseCase
	ResultUseCase() usecase.ResultUseCase
	SetCloudinaryInstance(cloudinary *cloudinary.Cloudinary)
	InterviewProcessUseCase() usecase.InterviewProcessUseCase
	UsersUseCase() usecase.UserUsecase          //user
	UserRolesUseCase() usecase.UserRolesUseCase         // user role
	HRRecruitmentUsecase() usecase.HRRecruitmentUsecase //hr
	InterviewerUseCase() usecase.InterviewerUseCase
	BootcampUseCase() usecase.BootcampUseCase
	StatusUseCase() usecase.StatusUseCase
	AuthUseCase() usecase.AuthUseCase
	FormUseCase() usecase.FormUseCase
	InterviewResultUseCase() usecase.InterviewResultUseCase
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

func (u *useCaseManager) FormUseCase() usecase.FormUseCase {
	return usecase.NewFormUseCase(u.repoManager.FormRepo(), u.repoManager.CloudinaryInstance())
}

func (u *useCaseManager) InterviewProcessUseCase() usecase.InterviewProcessUseCase {
	return usecase.NewInterviewProcessUseCase(u.repoManager.InterviewProcessRepo(), u.CandidateUseCase(), u.InterviewerUseCase(), u.StatusUseCase(), u.FormUseCase())
}

func (u *useCaseManager) InterviewResultUseCase() usecase.InterviewResultUseCase {
	return usecase.NewInterviewResultUseCase(u.repoManager.InterviewResultRepo(), u.InterviewProcessUseCase(), u.ResultUseCase())
}


func (u *useCaseManager) InterviewerUseCase() usecase.InterviewerUseCase {
	return usecase.NewInterviewerUseCase(u.repoManager.InterviewerRepo())
}
func (u *useCaseManager) ResultUseCase() usecase.ResultUseCase {
	return usecase.NewResultUseCase(u.repoManager.ResultRepo())
}

// user role
func (u *useCaseManager) UserRolesUseCase() usecase.UserRolesUseCase {
	return usecase.NewUserRolesUseCase(u.repoManager.UserRolesRepo())
}

// user
func (u *useCaseManager) UsersUseCase() usecase.UserUsecase {
	return usecase.NewUserUsecase(u.repoManager.UsersRepo())
}

func (u *useCaseManager) SetCloudinaryInstance(cloudinary *cloudinary.Cloudinary) {
	u.cloudinary = cloudinary
}

func (u *useCaseManager) AuthUseCase() usecase.AuthUseCase {
	return usecase.NewAuthUseCase(u.UsersUseCase(), u.UserRolesUseCase())
}

func (u *useCaseManager) HRRecruitmentUsecase() usecase.HRRecruitmentUsecase {
	return usecase.NewHRRecruitmentUsecase(u.repoManager.HRRecruitmentRepo(), u.repoManager.UsersRepo())
}

func NewUseCaseManager(repoManager RepoManager) UseCaseManager {
	return &useCaseManager{repoManager: repoManager}
}
