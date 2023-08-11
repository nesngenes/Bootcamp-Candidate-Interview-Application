package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/repository"
)

type InterviewResultUseCase interface {
	RegisterNewInterviewResult(payload model.InterviewResult) error
	FindByIdInterviewResult(id string) (dto.InterviewResultResponseDto, error)
	FindAllInterviewResult(requestPaging dto.PaginationParam) ([]dto.InterviewResultResponseDto, dto.Paging, error)
}
type interviewResultUseCase struct {
	repo                    repository.InterviewResultRepository
	interviewProcessUseCase interviewProcessUseCase
	resultUseCase           ResultUseCase
}

func (i *interviewResultUseCase) RegisterNewInterviewResult(newInterviewResult model.InterviewResult) error {
	// get candidate
	InterviewP, err := i.interviewProcessUseCase.FindByIdInterviewProcess(newInterviewResult.InterviewId)
	if err != nil {
		return fmt.Errorf("interviewProcess with ID %s not found", newInterviewResult.InterviewId)
	}
	result, err := i.resultUseCase.FindByIdResult(newInterviewResult.ResultId)
	if err != nil {
		return fmt.Errorf("result with ID %s not found", newInterviewResult.ResultId)
	}

	newInterviewResult.InterviewId = InterviewP.ID
	newInterviewResult.ResultId = result.ResultId

	err = i.repo.Create(newInterviewResult)
	if err != nil {
		return fmt.Errorf("failed to register new interview Result %v", err)
	}

	return nil
}
func (i *interviewResultUseCase) FindAllInterviewResult(requestPaging dto.PaginationParam) ([]dto.InterviewResultResponseDto, dto.Paging, error) {
	return i.repo.List(requestPaging)
}
func (i *interviewResultUseCase) FindByIdInterviewResult(id string) (dto.InterviewResultResponseDto, error) {
	var interviewResultResponseDto dto.InterviewResultResponseDto
	interviewResultResponse, err := i.repo.Get(id)
	if err != nil {
		return dto.InterviewResultResponseDto{}, fmt.Errorf("failed get by id interviewResult: %v", err.Error())
	}

	interviewResultResponseDto = interviewResultResponse
	return interviewResultResponseDto, nil

}

func NewInterviewResultUseCase(repo repository.InterviewResultRepository, interviewPUseCase InterviewProcessUseCase, resultUseCase ResultUseCase) InterviewResultUseCase {
	return &interviewResultUseCase{
		repo:                    repo,
		interviewProcessUseCase: interviewProcessUseCase{},
		resultUseCase:           resultUseCase,
	}
}
