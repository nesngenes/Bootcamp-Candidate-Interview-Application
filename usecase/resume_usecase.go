package usecase

import (
    "fmt"
    "interview_bootcamp/model"
    "interview_bootcamp/repository"
    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
    "context"
    "bytes"
)

type ResumeUseCase interface {
    RegisterNewResume(payload model.Resume) error
    FindAllResume() ([]model.Resume, error)
    FindByIdResume(id string) (model.Resume, error)
    UpdateResume(payload model.Resume) error
    DeleteResume(id string) error
}

type resumeUseCase struct {
    repo       repository.ResumeRepository
    cloudinary *cloudinary.Cloudinary
}

func (r *resumeUseCase) RegisterNewResume(payload model.Resume) error {
	if payload.ResumeID == "" && payload.CandidateID == "" {
		return fmt.Errorf("resume id and candidate id are required fields")
	}

	err := r.repo.Create(payload)
	if err != nil {
		return fmt.Errorf("failed to create new resume: %v", err)
	}

	return nil
}

// FindAllResume implements ResumeUseCase.
func (r *resumeUseCase) FindAllResume() ([]model.Resume, error) {
    // resumes, err := r.repo.GetAll() // Implement a method to retrieve all resumes from the repository
    // if err != nil {
    //     return nil, fmt.Errorf("failed to retrieve resumes: %v", err)
    // }
    // return resumes, nil
    panic("")
}

// FindByIdResume implements ResumeUseCase.
func (r *resumeUseCase) FindByIdResume(id string) (model.Resume, error) {
    resume, err := r.repo.Get(id) // Implement a method to retrieve a resume by its ID from the repository
    if err != nil {
        return model.Resume{}, fmt.Errorf("failed to retrieve resume: %v", err)
    }
    return resume, nil
}


func (r *resumeUseCase) DeleteResume(id string) error {
    // Retrieve the resume details before deleting
    resume, err := r.repo.Get(id)
    if err != nil {
        return fmt.Errorf("failed to retrieve resume: %v", err)
    }

    // Delete the resume record from the database
    err = r.repo.Delete(id)
    if err != nil {
        return fmt.Errorf("failed to delete resume: %v", err)
    }

    // Delete the file from Cloudinary
    publicID := "resumes/" + resume.ResumeID
    _, err = r.cloudinary.Upload.Destroy(context.Background(), uploader.DestroyParams{
        PublicID: publicID,
    })
    if err != nil {
        return fmt.Errorf("failed to delete file from Cloudinary: %v", err)
    }

    return nil
}

func (r *resumeUseCase) UpdateResume(payload model.Resume) error {
    if payload.ResumeID == "" {
        return fmt.Errorf("resume_id is required for update")
    }

    // Retrieve the existing resume details
    existingResume, err := r.repo.Get(payload.ResumeID)
    if err != nil {
        return fmt.Errorf("failed to retrieve existing resume: %v", err)
    }

    // Update the resume details in the database
    err = r.repo.Update(payload)
    if err != nil {
        return fmt.Errorf("failed to update resume: %v", err)
    }

    // If CvFile is provided, update the file on Cloudinary
    if len(payload.CvFile) > 0 {
        ctx := context.Background()

        updatedUploadResult, err := r.cloudinary.Upload.Upload(ctx, bytes.NewReader(payload.CvFile), uploader.UploadParams{
            PublicID: "resumes/" + payload.ResumeID,
        })
        if err != nil {
            fmt.Println("Error uploading updated file to Cloudinary:", err)
            return fmt.Errorf("error uploading updated file to Cloudinary: %v", err)
        }

        payload.CvURL = updatedUploadResult.SecureURL
    } else {
        // Keep the existing file details if no new file is uploaded
        payload.CvURL = existingResume.CvURL
        payload.CvFile = existingResume.CvFile
    }

    return nil
}

func NewResumeUseCase(repo repository.ResumeRepository, cloudinary *cloudinary.Cloudinary) ResumeUseCase {
    return &resumeUseCase{
        repo:       repo,
        cloudinary: cloudinary,
    }
}