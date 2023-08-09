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


func (r *resumeUseCase) DeleteResume(id string) error {
    panic("")
}

func (r *resumeUseCase) UpdateResume(payload model.Resume) error {
    panic("")
}

func NewResumeUseCase(repo repository.ResumeRepository) ResumeUseCase {
    return &resumeUseCase{repo: repo}
}

func NewResumeUseCase(repo repository.ResumeRepository, cloudinary *cloudinary.Cloudinary) ResumeUseCase {
    return &resumeUseCase{
        repo:       repo,
        cloudinary: cloudinary,
    }
}