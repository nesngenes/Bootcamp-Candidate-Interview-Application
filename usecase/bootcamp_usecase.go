package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/repository"
)

type BootcampUseCase interface {
	RegisterNewBootcamp(payload model.Bootcamp) error
	FindAllBootcamp(requesPaging dto.PaginationParam) ([]model.Bootcamp, dto.Paging, error)
	FindByIdBootcamp(id string) (model.Bootcamp, error)
	UpdateBootcamp(payload model.Bootcamp) error
	DeleteBootcamp(id string) error
	GetBootcampByID(id string) (model.Bootcamp, error)
}
type bootcampUseCase struct {
	repo repository.BootcampRepository
}

func (b *bootcampUseCase) RegisterNewBootcamp(payload model.Bootcamp) error {
	if payload.Name == "" || payload.Location == "" || payload.StartDate.String() == "" || payload.EndDate.String() == "" {
		return fmt.Errorf("name ,location strat_date and end_date required fields")
	}
	isExistCandidate, _ := b.repo.GetByName(payload.Name)
	if isExistCandidate.Name == payload.Name {
		return fmt.Errorf("candidate with name %s exists", payload.Name)
	}
	err := b.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new botcamp: %v", err)
	}
	return nil
}

func (b *bootcampUseCase) GetBootcampByID(id string) (model.Bootcamp, error) {
	bootcamp, err := b.repo.GetByID(id)
	if err != nil {
		return model.Bootcamp{}, fmt.Errorf("bootcamp with id %s not found", id)
	}
	return bootcamp, nil
}


func (b *bootcampUseCase) FindAllBootcamp(requesPaging dto.PaginationParam) ([]model.Bootcamp, dto.Paging, error) {
	return b.repo.Paging(requesPaging)
}

func (b *bootcampUseCase) FindByIdBootcamp(id string) (model.Bootcamp, error) {
	bootcamp, err := b.repo.Get(id)
	if err != nil {
		return model.Bootcamp{}, fmt.Errorf("bootcamp with id %s not found", id)
	}
	return bootcamp, nil
}
func (b *bootcampUseCase) UpdateBootcamp(payload model.Bootcamp) error {
	if payload.Name == "" || payload.Location == "" || payload.StartDate.String() == "" || payload.EndDate.String() == "" {
		return fmt.Errorf("name ,location strat_date and end_date required fields")
	}

	err := b.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update bootcamp: %v", err)
	}

	return nil
}
func (b *bootcampUseCase) DeleteBootcamp(id string) error {
	bootcamp, err := b.FindByIdBootcamp(id)
	if err != nil {
		return fmt.Errorf("bootcamp with ID %s not found", id)
	}

	err = b.repo.Delete(bootcamp.BootcampId)
	if err != nil {
		return fmt.Errorf("failed to delete bootcamp: %v", err.Error())
	}
	return nil
}
func NewBootcampUseCase(repo repository.BootcampRepository) BootcampUseCase {
	return &bootcampUseCase{repo: repo}
}
