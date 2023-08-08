package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/repository"
)

type CandidateUseCase interface {
	RegisterNewCandidate(payload model.Candidate) error
	FindAllCandidate() ([]model.Candidate, error)
	FindByIdCandidate(id string) (model.Candidate, error)
	UpdateCandidate(payload model.Candidate) error
	DeleteCandidate(id string) error
}

type candidateUseCase struct {
	repo repository.CandidateRepository
}

// RegisterNewCandidate implements CandidateUseCase.
func (c *candidateUseCase) RegisterNewCandidate(payload model.Candidate) error {
	//pengecekan nama tidak boleh kosong
	if payload.FirstName == "" && payload.LastName == "" && payload.Email == "" && payload.Phone == "" && payload.Address == "" {
		return fmt.Errorf("first name, last name, email, phone, address, date of birth required fields")
	}

	//pengecekan email tidak boleh sama
	// isExistCandidate, _ := c.repo.GetByEmail(payload.Email)
	// if isExistCandidate.Email == payload.Email {
	// 	return fmt.Errorf("candidate with email %s exists", payload.Email)
	// }

	//pengecekan phone number tidak boleh sama
	isExistCandidate, _ := c.repo.GetByPhoneNumber(payload.Phone)
	if isExistCandidate.Phone == payload.Phone {
		return fmt.Errorf("candidate with phoone %s exists", payload.Phone)
	}

	err := c.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new candidate: %v", err)
	}
	return nil
}

// FindAllCandidate implements CandidateUseCase.
func (c *candidateUseCase) FindAllCandidate() ([]model.Candidate, error) {
	panic("")
}

// FindByIdCandidate implements CandidateUseCase.
func (c *candidateUseCase) FindByIdCandidate(id string) (model.Candidate, error) {
	panic("")
}

// DeleteCandidate implements CandidateUseCase.
func (c *candidateUseCase) DeleteCandidate(id string) error {
	panic("")
}

// UpdateCandidate implements CandidateUseCase.
func (c *candidateUseCase) UpdateCandidate(payload model.Candidate) error {
	//untuk mengecek apakah kolom nomor sudah diisi
	if payload.Phone == "" {
		return fmt.Errorf("kolom nomor harus di isi")
	}

	//untuk mengecek apakah data dengan nomor tersebut sudah ada
	isExistCandidate, _ := c.repo.GetByPhoneNumber(payload.Phone)
	if isExistCandidate.Phone == payload.Phone && isExistCandidate.CandidateID != payload.CandidateID {
		return fmt.Errorf("data dengan nomor: %s sudah ada", payload.Phone)
	}

	//untuk melakukan update pada data dengan nomor sesuai kolom
	err := c.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("gagal memperbarui nomor: %v", err)
	}

	return nil
}

func NewCandidateUseCase(repo repository.CandidateRepository) CandidateUseCase {
	return &candidateUseCase{repo: repo}
}
