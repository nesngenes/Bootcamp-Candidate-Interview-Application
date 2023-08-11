package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/repository"
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
}

func (i *interviewProcessUseCase) RegisterNewInterviewProcess(newInterviewProses model.InterviewProcess) error {
	// get candidate
	candidate, err := i.canUseCase.FindByIdCandidate(newInterviewProses.CandidateID)
	if err != nil {
		return fmt.Errorf("candidate with ID %s not found", newInterviewProses.CandidateID)
	}
	interviewer, err := i.intUseCase.FindByIdInterviewer(newInterviewProses.InterviewerID)
	if err != nil {
		return fmt.Errorf("interviewer with ID %s not found", newInterviewProses.InterviewerID)
	}
	status, err := i.statusUseCase.FindByIdStatus(newInterviewProses.StatusID)
	if err != nil {
		return fmt.Errorf("status with ID %s not found", newInterviewProses.StatusID)
	}
	newInterviewProses.CandidateID = candidate.CandidateID
	newInterviewProses.InterviewerID = interviewer.InterviewerID
	newInterviewProses.StatusID = status.StatusId

	err = i.repo.Create(newInterviewProses)
	if err != nil {
		return fmt.Errorf("failed to register new interview proses %v", err)
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
func NewInterviewProcessUseCase(repo repository.InterviewProcessRepository, canUseCase CandidateUseCase, intUseCase InterviewerUseCase, statusUseCase StatusUseCase) InterviewProcessUseCase {
	return &interviewProcessUseCase{
		repo:          repo,
		canUseCase:    canUseCase,
		intUseCase:    intUseCase,
		statusUseCase: statusUseCase,
	}
}
