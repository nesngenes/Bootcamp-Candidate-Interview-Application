package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/repository"
)

type ResultUseCase interface {
	RegisterNewResult(payload model.Result) error
	FindAllResult(requesPaging dto.PaginationParam) ([]model.Result, dto.Paging, error)
	FindByIdResult(id string) (model.Result, error)
	UpdateResult(payload model.Result) error
	DeleteResult(id string) error
}
type resultUseCase struct {
	repo repository.ResultRepository
}

func (r *resultUseCase) RegisterNewResult(payload model.Result) error {
	if payload.Name == "" {
		return fmt.Errorf("name  required fields")
	}
	isExistResult, _ := r.repo.GetByName(payload.Name)
	if isExistResult.Name == payload.Name {
		return fmt.Errorf("ERR result with name %s exits", payload.Name)
	}
	err := r.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new result: %v", err)
	}
	return nil
}
func (r *resultUseCase) FindAllResult(requesPaging dto.PaginationParam) ([]model.Result, dto.Paging, error) {
	return r.repo.Paging(requesPaging)
}
func (r *resultUseCase) FindByIdResult(id string) (model.Result, error) {
	result, err := r.repo.Get(id)
	if err != nil {
		return model.Result{}, fmt.Errorf("result with id %s not found", id)
	}
	return result, nil
}
func (r *resultUseCase) UpdateResult(payload model.Result) error {
	if payload.Name == "" {
		return fmt.Errorf("name  required fields")
	}

	err := r.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update result: %v", err)
	}

	return nil
}
func (r *resultUseCase) DeleteResult(id string) error {
	result, err := r.FindByIdResult(id)
	if err != nil {
		return fmt.Errorf("result with ID %s not found", id)
	}

	err = r.repo.Delete(result.ResultId)
	if err != nil {
		return fmt.Errorf("failed to delete result: %v", err.Error())
	}
	return nil
}
func NewResultUseCase(repo repository.ResultRepository) ResultUseCase {
	return &resultUseCase{repo: repo}
}
