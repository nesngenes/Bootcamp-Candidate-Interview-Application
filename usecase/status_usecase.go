package usecase

import (
	"fmt"

	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/repository"
)

type StatusUseCase interface {
	RegisterNewStatus(payload model.Status) error
	FindAllStatus(requesPaging dto.PaginationParam) ([]model.Status, dto.Paging, error)
	FindByIdStatus(id string) (model.Status, error)
	UpdateStatus(payload model.Status) error
	DeleteStatus(id string) error
}
type statusUseCase struct {
	repo repository.StatusRepository
}

func (s *statusUseCase) RegisterNewStatus(payload model.Status) error {
	if payload.Name == "" {
		return fmt.Errorf("name  required fields")
	}
	isExistStatus, _ := s.repo.GetByName(payload.Name)
	if isExistStatus.Name == payload.Name {
		return fmt.Errorf("ERR status with name %s exits", payload.Name)
	}
	err := s.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new status: %v", err)
	}
	return nil
}
func (s *statusUseCase) FindAllStatus(requesPaging dto.PaginationParam) ([]model.Status, dto.Paging, error) {
	return s.repo.Paging(requesPaging)
}
func (s *statusUseCase) FindByIdStatus(id string) (model.Status, error) {
	status, err := s.repo.Get(id)
	if err != nil {
		return model.Status{}, fmt.Errorf("status with id %s not found", id)
	}
	return status, nil
}
func (s *statusUseCase) UpdateStatus(payload model.Status) error {
	if payload.Name == "" {
		return fmt.Errorf("name  required fields")
	}

	err := s.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update status: %v", err)
	}

	return nil
}
func (s *statusUseCase) DeleteStatus(id string) error {
	status, err := s.FindByIdStatus(id)
	if err != nil {
		return fmt.Errorf("status with ID %s not found", id)
	}

	err = s.repo.Delete(status.StatusId)
	if err != nil {
		return fmt.Errorf("failed to delete status: %v", err.Error())
	}
	return nil
}
func NewStatusUseCase(repo repository.StatusRepository) StatusUseCase {
	return &statusUseCase{repo: repo}
}
