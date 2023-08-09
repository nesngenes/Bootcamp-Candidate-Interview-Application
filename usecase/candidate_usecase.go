package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/repository"
	"github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"context"
)

type CandidateUseCase interface {
	RegisterNewCandidate(payload model.Candidate) error
	FindAllCandidate(dto.PaginationParam) ([]model.Candidate, dto.Paging, error)
	FindByIdCandidate(id string) (model.Candidate, error)
	UpdateCandidate(payload model.Candidate) error
	DeleteCandidate(id string) error
}

type candidateUseCase struct {
	repo       repository.CandidateRepository
	bootcampUC BootcampUseCase
	cloudinary *cloudinary.Cloudinary
}

// RegisterNewCandidate implements CandidateUseCase.
func (c *candidateUseCase) RegisterNewCandidate(payload model.Candidate) error {
	//pengecekan nama tidak boleh kosong
	if payload.FullName == "" && payload.Phone == "" && payload.Email == "" && payload.Address == "" {
		return fmt.Errorf("fullname, email, phone, address, date of birth required fields")
	}

	// pengecekan email tidak boleh sama
	isExistCandidateEmail, _ := c.repo.GetByEmail(payload.Email)
	if isExistCandidateEmail.Email == payload.Email {
		return fmt.Errorf("candidate with email %s exists", payload.Email)
	}
	


	isExistCandidatePhone, _ := c.repo.GetByPhoneNumber(payload.Phone)
	if isExistCandidatePhone.Phone == payload.Phone {
		return fmt.Errorf("candidate with phoone %s exists", payload.Phone)
	}

	_, err := c.bootcampUC.FindByIdBootcamp(payload.Bootcamp.BootcampId)
	if err != nil {
		return fmt.Errorf("bootcamp with ID %s not found", payload.Bootcamp.BootcampId)
	}

	err = c.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new candidate: %v", err)
	}
	return nil
}

// FindAllCandidate implements CandidateUseCase.
func (c *candidateUseCase) FindAllCandidate(requesPaging dto.PaginationParam) ([]model.Candidate, dto.Paging, error) {
	return c.repo.Paging(requesPaging)
}

// FindByIdCandidate implements CandidateUseCase.
func (c *candidateUseCase) FindByIdCandidate(id string) (model.Candidate, error) {
	candidate, err := c.repo.Get(id)
	if err != nil {
		return model.Candidate{}, fmt.Errorf("candidate with id %s not found", id)
	}
	return candidate, nil
}

// DeleteCandidate implements CandidateUseCase.
func (c *candidateUseCase) DeleteCandidate(id string) error {
	candidate, err := c.FindByIdCandidate(id)
	if err != nil {
		return fmt.Errorf("candidate with ID %s not found", id)
	}

	err = c.repo.Delete(candidate.CandidateID)
	if err != nil {
		return fmt.Errorf("failed to delete candidate: %v", err.Error())
	}

	// Hapus file dari cloudinary
	 publicID := "candidates/" + candidate.CandidateID
	 _, err = c.cloudinary.Upload.Destroy(context.Background(), uploader.DestroyParams{
		 PublicID: publicID,
	 })
	 if err != nil {
		 return fmt.Errorf("failed to delete file from Cloudinary: %v", err)
	 }
	return nil
}

// UpdateCandidate implements CandidateUseCase.
func (c *candidateUseCase) UpdateCandidate(payload model.Candidate) error {

	if payload.Phone == "" {
		return fmt.Errorf("number phone form must be fill")
	}

	_, err := c.bootcampUC.FindByIdBootcamp(payload.Bootcamp.BootcampId)
	if err != nil {
		return fmt.Errorf("bootcamp with ID %s not found", payload.Bootcamp.BootcampId)
	// pengecekan email tidak boleh sama
	isExistCandidateS, _ := c.repo.GetByEmail(payload.Email)
	if isExistCandidateS.Email == payload.Email {
		return fmt.Errorf("candidate with email %s exists", payload.Email)
	}

	//untuk mengecek apakah data dengan nomor tersebut sudah ada
	isExistCandidate, _ := c.repo.GetByPhoneNumber(payload.Phone)
	if isExistCandidate.Phone == payload.Phone && isExistCandidate.CandidateID != payload.CandidateID {
		return fmt.Errorf("data dengan nomor: %s sudah ada", payload.Phone)
	}

	err = c.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("gagal memperbarui nomor: %v", err)
	}

	return nil
}

func NewCandidateUseCase(repo repository.CandidateRepository, bootcampUC BootcampUseCase, cloudinary *cloudinary.Cloudinary) CandidateUseCase {
	return &candidateUseCase{
		repo: repo, 
		bootcampUC: bootcampUC,
		cloudinary: cloudinary,
	}
}
