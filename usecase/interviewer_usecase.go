package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/repository"
)

type InterviewerUseCase interface {
	RegisterNewInterviewer(payload model.Interviewer) error
	FindAllInterviewer() ([]model.Interviewer, error)
	FindByIdInterviewer(id string) (model.Interviewer, error)
	UpdateInterviewer(payload model.Interviewer) error
	DeleteInterviewer(id string) error
}

type interviewerUseCase struct {
	repo repository.InterviewerRepository
}

// RegisterNewInterviewer implements InterviewerUseCase.
func (i *interviewerUseCase) RegisterNewInterviewer(payload model.Interviewer) error {
	//pengecekan field tidak boleh kosong
	if payload.InterviewerID == "" && payload.FullName == "" && payload.UserID == "" {
		return fmt.Errorf("id, full name, user id required fields")
	}

	err := i.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new interviewer: %v", err)
	}
	return nil
}

// FindAllInterviewer implements InterviewerUseCase.
func (i *interviewerUseCase) FindAllInterviewer() ([]model.Interviewer, error) {
	return i.repo.List()
}

// FindByIdInterviewer implements InterviewerUseCase.
func (i *interviewerUseCase) FindByIdInterviewer(id string) (model.Interviewer, error) {
	return i.repo.Get(id)
}

// DeleteInterviewer implements InterviewerUseCase.
func (i *interviewerUseCase) DeleteInterviewer(id string) error {
	interviewer, err := i.FindByIdInterviewer(id)
	if err != nil {
		return fmt.Errorf("interviewer with ID %s not found", id)
	}

	err = i.repo.Delete(interviewer.InterviewerID)
	if err != nil {
		return fmt.Errorf("failed to delete interviewer: %v", err.Error())
	}
	return nil
}

// UpdateInterviewer implements InterviewerUseCase.
func (i *interviewerUseCase) UpdateInterviewer(payload model.Interviewer) error {
	if payload.InterviewerID == "" || payload.FullName == "" || payload.UserID == "" {
		return fmt.Errorf("id, full name , user id are required fields")
	}

	err := i.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update interviewer: %v", err.Error())
	}
	return nil

}

func NewInterviewerUseCase(repo repository.InterviewerRepository) InterviewerUseCase {
	return &interviewerUseCase{repo: repo}
}
