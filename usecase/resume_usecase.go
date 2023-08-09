package usecase

import (
	"fmt"
	"interview_bootcamp/model"
	"interview_bootcamp/repository"
)

type ResumeUseCase interface {
	RegisterNewResume(payload model.Resume) error
	FindAllResume() ([]model.Resume, error)
	FindByIdResume(id string) (model.Resume, error)
	UpdateResume(payload model.Resume) error
	DeleteResume(id string) error
}

type resumeUseCase struct {
	repo repository.ResumeRepository
}

func (r *resumeUseCase) RegisterNewResume(payload model.Resume) error {
	if payload.ResumeID == "" && payload.CandidateID == "" {
		return fmt.Errorf("resume id and candidate id are required fields")
	}
	fmt.Println("CANDIDATE ID WOUYYYYYYY")
	fmt.Println(payload.CandidateID)
	err := r.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new resume: %v", err)
	}

	return nil
}

// FindAllResume implements ResumeUseCase.
func (r *resumeUseCase) FindAllResume() ([]model.Resume, error) {
	panic("")
}

// FindByIdResume implements ResumeUseCase.
func (r *resumeUseCase) FindByIdResume(id string) (model.Resume, error) {
	panic("")
}

// DeleteResume implements ResumeUseCase.
func (r *resumeUseCase) DeleteResume(id string) error {
	panic("")
}

// UpdateResume implements ResumeUseCase.
func (r *resumeUseCase) UpdateResume(payload model.Resume) error {
	//untuk mengecek apakah kolom nomor sudah diisi
	if payload.ResumeID == "" {
		return fmt.Errorf("kolom resume id harus di isi")
	}

	// pengecekan id resume tidak boleh sama
	isExistResume, _ := r.repo.FindByIdResume(payload.ResumeID)
	if isExistResume.ResumeID == payload.ResumeID {
		return fmt.Errorf("candidate with email %s exists", payload.ResumeID)
	}

	//untuk melakukan update pada data dengan id resume sesuai kolom
	err := r.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("gagal memperbarui resume id: %v", err)
	}

	return nil
}

func NewResumeUseCase(repo repository.ResumeRepository) ResumeUseCase {
	return &resumeUseCase{repo: repo}
}
