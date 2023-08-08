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
func (c *interviewerUseCase) RegisterNewInterviewer(payload model.Interviewer) error {
	//pengecekan field tidak boleh kosong
	if payload.FirstName == "" && payload.LastName == "" && payload.Email == "" && payload.Phone == "" && payload.Specialization == "" {
		return fmt.Errorf("first name, last name, email, phone,specialization required fields")
	}

	// pengecekan email tidak boleh sama
	isExistInterviewer, _ := c.repo.GetByEmail(payload.Email)
	if isExistInterviewer.Email == payload.Email {
		return fmt.Errorf("interviewer with email %s exists", payload.Email)
	}

	//pengecekan phone number tidak boleh sama
	isExistInterviewers, _ := c.repo.GetByPhoneNumber(payload.Phone)
	if isExistInterviewers.Phone == payload.Phone {
		return fmt.Errorf("interviewer with phone %s exists", payload.Phone)
	}

	err := c.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new interviewer: %v", err)
	}
	return nil
}

// FindAllInterviewer implements InterviewerUseCase.
func (c *interviewerUseCase) FindAllInterviewer() ([]model.Interviewer, error) {
	return c.repo.List()
}

// FindByIdInterviewer implements InterviewerUseCase.
func (c *interviewerUseCase) FindByIdInterviewer(id string) (model.Interviewer, error) {
	return c.repo.Get(id)
}

// DeleteInterviewer implements InterviewerUseCase.
func (c *interviewerUseCase) DeleteInterviewer(id string) error {
	interviewer, err := c.FindByIdInterviewer(id)
	if err != nil {
		return fmt.Errorf("interviewer with ID %s not found", id)
	}

	err = c.repo.Delete(interviewer.InterviewerID)
	if err != nil {
		return fmt.Errorf("failed to delete interviewer: %v", err.Error())
	}
	return nil
}

// UpdateInterviewer implements InterviewerUseCase.
func (c *interviewerUseCase) UpdateInterviewer(payload model.Interviewer) error {
	if payload.FirstName == "" || payload.Phone == "" || payload.LastName == "" || payload.Specialization == "" {
		return fmt.Errorf("first name, last name , phone number, specialization are required fields")
	}

	// pengecekan email tidak boleh sama
	isExistInterviewer, _ := c.repo.GetByEmail(payload.Email)
	if isExistInterviewer.Email == payload.Email {
		return fmt.Errorf("interviewer with email %s exists", payload.Email)
	}

	interviewer, _ := c.repo.GetByPhoneNumber(payload.Phone)
	if interviewer.Phone == payload.Phone && interviewer.InterviewerID != payload.InterviewerID {
		return fmt.Errorf("interviewer with phone number %s already exists", payload.Phone)
	}
	err := c.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update interviewer: %v", err.Error())
	}
	return nil

}

func NewInterviewerUseCase(repo repository.InterviewerRepository) InterviewerUseCase {
	return &interviewerUseCase{repo: repo}
}
