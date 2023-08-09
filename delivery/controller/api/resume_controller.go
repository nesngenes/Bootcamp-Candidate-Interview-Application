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
	router  *gin.Engine
	usecase usecase.ResumeUseCase
}

func (rc *ResumeController) createHandler(c *gin.Context) {
	cfg, err := config.NewConfig()
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

	fmt.Println("File Content Size:", len(fileContent))

	// Upload file to Cloudinary
	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryConfig.CloudinaryCloudName,
		cfg.CloudinaryConfig.CloudinaryAPIKey,
		cfg.CloudinaryConfig.CloudinaryAPISecret,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error initializing Cloudinary"})
		return
	}

	ctx := context.Background()

	uploadResult, err := cld.Upload.Upload(ctx, bytes.NewReader(fileContent), uploader.UploadParams{
		PublicID: "resumes/" + resume.ResumeID, // Use a meaningful folder structure
	})
	if err != nil {
		fmt.Println("Error uploading to Cloudinary:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error uploading to Cloudinary"})
		return
	}

	fmt.Println("Upload Result:", uploadResult)

	// Save the Cloudinary URL and CvFile info to the resume model
	resume.CvURL = uploadResult.SecureURL
	resume.CvFile = fileContent

	if err := rc.usecase.RegisterNewResume(resume); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resume)
}

func NewResumeController(r *gin.Engine, usecase usecase.ResumeUseCase) *ResumeController {
	controller := ResumeController{
		router:  r,
		usecase: usecase,
	}
	rg := r.Group("/api/v1")
	rg.POST("/resumes", controller.createHandler)
	// ... (other routes)
	return &controller
}
