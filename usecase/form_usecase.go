package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/repository"
	"github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"interview_bootcamp/model/dto"
	"context"
)

type FormUseCase interface {
	RegisterNewForm(payload model.Form) error
	DeleteForm(id string) error
	FindByIdForm(id string) (model.Form, error)
    UpdateForm(payload model.Form) error
	FindAllForm(dto.PaginationParam) ([]model.Form, dto.Paging, error)
}

type formUseCase struct {
	repo repository.FormRepository
	cloudinary *cloudinary.Cloudinary
}

func (f *formUseCase) RegisterNewForm(payload model.Form) error {
	err := f.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new form: %v", err)
	}
	return nil
}

func (f *formUseCase) FindByIdForm(id string) (model.Form, error) {
    form, err := f.repo.Get(id) // Implement a method to retrieve a form by its ID from the repository
    if err != nil {
        return model.Form{}, fmt.Errorf("failed to retrieve form: %v", err)
    }
    return form, nil
}

func (f *formUseCase) FindAllForm(requestPaging dto.PaginationParam) ([]model.Form, dto.Paging, error) {
	forms, paging, err := f.repo.Paging(requestPaging)
	if err != nil {
		return nil, dto.Paging{}, fmt.Errorf("failed to retrieve forms: %v", err)
	}
	return forms, paging, nil
}

func (f *formUseCase) DeleteForm(id string) error {
	form, err := f.FindByIdForm(id)
	if err != nil {
		return fmt.Errorf("form with ID %s not found", id)
	}

	err = f.repo.Delete(form.FormID)
	if err != nil {
		return fmt.Errorf("failed to delete form: %v", err.Error())
	}

	// Hapus file dari cloudinary
	publicID := "forms/" + form.FormID
	_, err = f.cloudinary.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from Cloudinary: %v", err)
	}	
	return nil
}

// UpdateForm implements FormUseCase.
func (f *formUseCase) UpdateForm(payload model.Form) error {
	err := f.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update form: %v", err)
	}
	return nil
}

func NewFormUseCase(repo repository.FormRepository, cloudinary *cloudinary.Cloudinary) FormUseCase {
	return &formUseCase{
		repo: repo,
		cloudinary: cloudinary,
	}
}