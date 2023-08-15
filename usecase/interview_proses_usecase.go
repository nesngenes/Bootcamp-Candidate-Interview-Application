package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/repository"
	"interview_bootcamp/usecase/email" // Add this import
)

type InterviewProcessUseCase interface {
	RegisterNewInterviewProcess(payload model.InterviewProcess) error
	FindByIdInterviewProcess(id string) (dto.InterviewProcessResponseDto, error)
	FindAllInterviewProcess(requestPaging dto.PaginationParam) ([]dto.InterviewProcessResponseDto, dto.Paging, error)
}
type interviewProcessUseCase struct {
	repo          repository.InterviewProcessRepository
	canUseCase    CandidateUseCase
	intUseCase    InterviewerUseCase
	statusUseCase StatusUseCase
	formUseCase   FormUseCase
	emailService  email.EmailService // Add this field
}


func (i *interviewProcessUseCase) RegisterNewInterviewProcess(newInterviewProcess model.InterviewProcess) error {
	// get candidate
	candidate, err := i.canUseCase.FindByIdCandidate(newInterviewProcess.CandidateID)
	if err != nil {
		return fmt.Errorf("candidate with ID %s not found", newInterviewProcess.CandidateID)
	}
	interviewer, err := i.intUseCase.FindByIdInterviewer(newInterviewProcess.InterviewerID)
	if err != nil {
		return fmt.Errorf("interviewer with ID %s not found", newInterviewProcess.InterviewerID)
	}
	status, err := i.statusUseCase.FindByIdStatus(newInterviewProcess.StatusID)
	if err != nil {
		return fmt.Errorf("status with ID %s not found", newInterviewProcess.StatusID)
	}
	form, err := i.formUseCase.FindByIdForm(newInterviewProcess.FormID) // Retrieve the form using Form.ID
	if err != nil {
		return fmt.Errorf("form with ID %s not found", newInterviewProcess.FormID)
	}

	newInterviewProcess.CandidateID = candidate.CandidateID
	newInterviewProcess.InterviewerID = interviewer.InterviewerID
	newInterviewProcess.StatusID = status.StatusId
	newInterviewProcess.FormLink = form.FormLink // Set the FormLink based on the retrieved form

	err = i.repo.Create(newInterviewProcess)
	if err != nil {
		return fmt.Errorf("failed to register new interview process: %v", err)
	}

    // Send interview email to the candidate
	candidateEmail := candidate.Email
	interviewDetails := fmt.Sprintf("Interviewer: %s<br>Interview Date: %s<br>Meeting Link: %s",
	interviewer.FullName, newInterviewProcess.InterviewDatetime, newInterviewProcess.MeetingLink)

	err = i.emailService.SendInterviewEmail(candidateEmail, interviewDetails)
	if err != nil {
		// You can log the error or take other actions
		return fmt.Errorf("failed to send interview email: %v", err)
	}

	return nil
}

func (i *interviewProcessUseCase) FindAllInterviewProcess(requestPaging dto.PaginationParam) ([]dto.InterviewProcessResponseDto, dto.Paging, error) {
	return i.repo.List(requestPaging)
}
func (i *interviewProcessUseCase) FindByIdInterviewProcess(id string) (dto.InterviewProcessResponseDto, error) {
	var interviewPrResponseDto dto.InterviewProcessResponseDto
	interviewPrResponse, err := i.repo.Get(id)
	if err != nil {
		return dto.InterviewProcessResponseDto{}, fmt.Errorf("failed get by id interviewProcess: %v", err.Error())
	}

	interviewPrResponseDto = interviewPrResponse
	return interviewPrResponseDto, nil

}
func NewInterviewProcessUseCase(repo repository.InterviewProcessRepository, canUseCase CandidateUseCase, intUseCase InterviewerUseCase, statusUseCase StatusUseCase, formUseCase FormUseCase) InterviewProcessUseCase {
	return &interviewProcessUseCase{
		repo:          repo,
		canUseCase:    canUseCase,
		intUseCase:    intUseCase,
		statusUseCase: statusUseCase,
		formUseCase:   formUseCase,
	}
}
