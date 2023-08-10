package api

import (
	"bytes"
	"context"
	"fmt"
	"interview_bootcamp/config"
	"interview_bootcamp/model"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
	"io"
	"net/http"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type ResumeController struct {
    router      *gin.Engine
    usecase     usecase.ResumeUseCase
    cloudinary  *cloudinary.Cloudinary
}

func (rc *ResumeController) createHandler(c *gin.Context) {
	_, err := config.NewConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error reading config"})
		return
	}

	// Parse JSON data from the form
	var resume model.Resume
	err = c.ShouldBind(&resume)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Generate ResumeID
	resume.ResumeID = common.GenerateID()

	// Extract candidate ID from JSON data
	candidateID := c.PostForm("candidate_id")
	resume.CandidateID = candidateID

	// Retrieve the uploaded file
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Error retrieving the file"})
		return
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error reading file content"})
		return
	}

	// Create a channel to communicate results
    resultCh := make(chan error)

    // Use a goroutine to upload the file to Cloudinary
    go func() {
        ctx := context.Background()

        uploadResult, err := rc.cloudinary.Upload.Upload(ctx, bytes.NewReader(fileContent), uploader.UploadParams{
            PublicID: "resumes/" + resume.ResumeID, // Use a meaningful folder structure
        })
        if err != nil {
            resultCh <- err
            return
        }

        // Save the Cloudinary URL and CvFile info to the resume model
        resume.CvURL = uploadResult.SecureURL
        resume.CvFile = fileContent

        resultCh <- nil
    }()

    // Create the resume record in the database concurrently
    go func() {
        if err := rc.usecase.RegisterNewResume(resume); err != nil {
            resultCh <- err
            return
        }
        resultCh <- nil
    }()

    // Wait for both goroutines to finish
    for i := 0; i < 2; i++ {
        if err := <-resultCh; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
            return
        }
    }

    c.JSON(http.StatusCreated, resume)
}

func (rc *ResumeController) deleteHandler(c *gin.Context) {
    resumeID := c.Param("resume_id")
    
    // Create a channel to communicate result
    resultCh := make(chan error)
    
    // Use a goroutine to initiate the deletion process
    go func() {
        if err := rc.usecase.DeleteResume(resumeID); err != nil {
            resultCh <- err
            return
        }
        resultCh <- nil
    }()

    // Wait for the goroutine to finish
    if err := <-resultCh; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Resume deleted successfully"})
}

func (rc *ResumeController) updateHandler(c *gin.Context) {
    // Retrieve the resume_id and cv_file_name from the form data
    resumeID := c.PostForm("resume_id")

    // Retrieve the uploaded file, if any
    file, _, err := c.Request.FormFile("file")
    if err != nil {
        file = nil
    }
    defer func() {
        if file != nil {
            file.Close()
        }
    }()

    fmt.Printf("Received Resume ID: %s\n", resumeID)

    // If the resume_id is provided in the form data, use it for the update
    if resumeID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"err": "resume_id is required for update"})
        return
    }

    // Retrieve the existing resume details
    existingResume, err := rc.usecase.FindByIdResume(resumeID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
        return
    }

    // Create a channel to communicate results
    resultCh := make(chan error)
    
    // If a new file is uploaded, read its content
    if file != nil {
        // Use a goroutine to upload the file to Cloudinary
        go func() {
            ctx := context.Background()

            updatedFileContent, err := io.ReadAll(file)
            if err != nil {
                resultCh <- err
                return
            }

            updatedUploadResult, err := rc.cloudinary.Upload.Upload(ctx, bytes.NewReader(updatedFileContent), uploader.UploadParams{
                PublicID: "resumes/" + resumeID, // Use the provided resume_id for update
            })
            if err != nil {
                resultCh <- err
                return
            }

            existingResume.CvURL = updatedUploadResult.SecureURL
            existingResume.CvFile = updatedFileContent

            resultCh <- nil
        }()
    } else {
        // If no file update, send a nil error to the channel
        resultCh <- nil
    }

    // Wait for the goroutine to finish if applicable
    err = <-resultCh

    // Update the resume details in the database
    if err == nil {
        err = rc.usecase.UpdateResume(existingResume)
    }

    // Respond based on the final error
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
        return
    }

    c.JSON(http.StatusOK, existingResume)
}

func NewResumeController(r *gin.Engine, usecase usecase.ResumeUseCase, cloudinary *cloudinary.Cloudinary) *ResumeController {
	controller := &ResumeController{
		router:      r,
		usecase:     usecase,
		cloudinary:  cloudinary,
	}
	rg := r.Group("/api/v1")
	rg.POST("/resumes", controller.createHandler)
	rg.DELETE("/resumes/:resume_id", controller.deleteHandler)
	rg.PUT("/resumes", controller.updateHandler)
	// ... (other routes)
	return controller
}