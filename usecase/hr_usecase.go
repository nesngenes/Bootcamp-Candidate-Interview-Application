package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/repository"
)

type HRRecruitmentUsecase interface {
	CreateHRRecruitment(payload model.HRRecruitment) error
	Get(id string) (model.HRRecruitment, error)
	ListHRRecruitments() ([]model.HRRecruitment, error)
	UpdateHRRecruitment(payload model.HRRecruitment) error
	DeleteHRRecruitment(id string) error
}

type hrRecruitmentUsecase struct {
	hrRecruitmentRepo repository.HRRecruitmentRepository
	userRepo          repository.UserRepository
}

func (u *hrRecruitmentUsecase) CreateHRRecruitment(payload model.HRRecruitment) error {
	// Check if the new user ID is already in use by another HR recruitment record
	_, err := u.userRepo.Get(payload.UserID)
	if err != nil {
		return fmt.Errorf("user with ID %s not found", payload.UserID)
	}

	existingHRRecruitment, err := u.hrRecruitmentRepo.GetByUserID(payload.UserID)
	if err == nil && existingHRRecruitment.ID != "" {
		return fmt.Errorf("user ID %s is already in use by another HR recruitment record", payload.UserID)
	}

	return u.hrRecruitmentRepo.Create(payload)
}

func (u *hrRecruitmentUsecase) Get(id string) (model.HRRecruitment, error) {
	hrRecruitment, err := u.hrRecruitmentRepo.Get(id)
	if err != nil {
		return model.HRRecruitment{}, err
	}

	return hrRecruitment, nil
}

func (u *hrRecruitmentUsecase) ListHRRecruitments() ([]model.HRRecruitment, error) {
	return u.hrRecruitmentRepo.List()
}

func (u *hrRecruitmentUsecase) UpdateHRRecruitment(payload model.HRRecruitment) error {
	// Check if the HR recruitment record exists
	_, err := u.hrRecruitmentRepo.Get(payload.ID)
	if err != nil {
		return fmt.Errorf("HR recruitment record with ID %s not found", payload.ID) //brarti user_id nya emg ga ada.
	}

	// Check if the new user ID is already in use by another HR recruitment record
	otherHRRecruitment, err := u.hrRecruitmentRepo.GetByUserID(payload.UserID)
	if err == nil && otherHRRecruitment.ID != payload.ID {
		return fmt.Errorf("user ID %s is already in use by another HR recruitment record", payload.UserID) //ada hr lain yang udah make user_id nya
	}

	// Update the HR recruitment record
	err = u.hrRecruitmentRepo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update HR recruitment record: %v", err)
	}

	return nil
}

func (u *hrRecruitmentUsecase) DeleteHRRecruitment(id string) error {
	return u.hrRecruitmentRepo.Delete(id)
}

func NewHRRecruitmentUsecase(hrRecruitmentRepo repository.HRRecruitmentRepository, userRepo repository.UserRepository) HRRecruitmentUsecase {
	return &hrRecruitmentUsecase{
		hrRecruitmentRepo: hrRecruitmentRepo,
		userRepo:          userRepo,
	}
}
