package api

import (
	"bytes"
	"context"
	"fmt"

	// "interview_bootcamp/delivery/middleware"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
	"io"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

type CandidateController struct {
	router           *gin.Engine
	candidateUsecase usecase.CandidateUseCase
	bootcampUsecase  usecase.BootcampUseCase
	cloudinary       *cloudinary.Cloudinary
}

func (cc *CandidateController) createHandler(c *gin.Context) {
	// Parse the form data
	err := c.Request.ParseMultipartForm(10 << 20) // Max 10MB file size
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Error parsing form data"})
		return
	}

	// Retrieve the form fields
	fullName := c.PostForm("full_name")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	dateOfBirth := c.PostForm("date_of_birth")
	address := c.PostForm("address")
	instansiPendidikan := c.PostForm("instansi_pendidikan")
	hackerRank, _ := strconv.Atoi(c.PostForm("hackerrank_score"))
	bootcampID := c.PostForm("bootcamp")

	// Create a new candidate instance
	candidate := model.Candidate{
		CandidateID:        common.GenerateID(),
		FullName:           fullName,
		Phone:              phone,
		Email:              email,
		DateOfBirth:        dateOfBirth,
		Address:            address,
		InstansiPendidikan: instansiPendidikan,
		HackerRank:         hackerRank,
	}

	// Create a new file upload form field
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

	// Upload file to cloudinary
	uploadResult, err := cc.cloudinary.Upload.Upload(
		context.Background(), bytes.NewReader(fileContent),
		uploader.UploadParams{
			PublicID: "candidates/" + candidate.CandidateID,
		},
	)
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error uploading file to cloudinary:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error uploading file to cloudinary"})
		return
	}
	candidate.CvLink = uploadResult.SecureURL

	// Fetch the Bootcamp details by its ID
	bootcamp, err := cc.bootcampUsecase.GetBootcampByID(bootcampID)
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error fetching bootcamp details:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error fetching bootcamp details"})
		return
	}
	candidate.Bootcamp = bootcamp

	// Register the new candidate
	err = cc.candidateUsecase.RegisterNewCandidate(candidate)
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error registering new candidate:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error registering new candidate"})
		return
	}

	status := common.WebStatus{
		Code:        http.StatusCreated,
		Description: "Create Data Successfully",
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": status,
	})
}

func (cc *CandidateController) listHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	paginationParam := dto.PaginationParam{
		Page:  page,
		Limit: limit,
	}

	candidates, paging, err := cc.candidateUsecase.FindAllCandidate(paginationParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	status := common.WebStatus{
		Code:        http.StatusOK,
		Description: "Get All Data Successfully",
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"data":   candidates,
		"paging": paging,
	})
}
func (cc *CandidateController) getHandler(c *gin.Context) {
	id := c.Param("id")
	product, err := cc.candidateUsecase.FindByIdCandidate(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	status := map[string]any{
		"code":        http.StatusOK,
		"description": "Get By Id Data Successfully",
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"data":   product,
	})

}

func (cc *CandidateController) updateHandler(c *gin.Context) {
	// Parse the form data
	err := c.Request.ParseMultipartForm(10 << 20) // Max 10MB file size
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Error parsing form data"})
		return
	}

	// Retrieve the form fields
	candidateID := c.PostForm("candidate_id")
	fullName := c.PostForm("full_name")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	dateOfBirth := c.PostForm("date_of_birth")
	address := c.PostForm("address")
	instansiPendidikan := c.PostForm("instansi_pendidikan")
	hackerRank, _ := strconv.Atoi(c.PostForm("hackerrank_score"))
	bootcampID := c.PostForm("bootcamp")

	// Fetch the existing candidate by ID
	existingCandidate, err := cc.candidateUsecase.FindByIdCandidate(candidateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	// Update the candidate's fields
	existingCandidate.FullName = fullName
	existingCandidate.Phone = phone
	existingCandidate.Email = email
	existingCandidate.DateOfBirth = dateOfBirth
	existingCandidate.Address = address
	existingCandidate.InstansiPendidikan = instansiPendidikan
	existingCandidate.HackerRank = hackerRank

	// Create a new file upload form field
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		// If no new file is uploaded, proceed without updating the file
	} else {
		defer file.Close()

		// Read file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "Error reading file content"})
			return
		}

		// Upload file to cloudinary
		uploadResult, err := cc.cloudinary.Upload.Upload(
			context.Background(), bytes.NewReader(fileContent),
			uploader.UploadParams{
				PublicID: "candidates/" + candidateID,
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "Error uploading file to cloudinary"})
			return
		}
		existingCandidate.CvLink = uploadResult.SecureURL
	}

	// Fetch the Bootcamp details by its ID
	bootcamp, err := cc.bootcampUsecase.GetBootcampByID(bootcampID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error fetching bootcamp details"})
		return
	}
	existingCandidate.Bootcamp = bootcamp

	// Update the candidate in the database
	err = cc.candidateUsecase.UpdateCandidate(existingCandidate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error updating candidate"})
		return
	}

	status := map[string]interface{}{
		"code":        http.StatusOK,
		"description": "Update data successfully",
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}

func (cc *CandidateController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := cc.candidateUsecase.DeleteCandidate(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	status := map[string]any{
		"code":        http.StatusNoContent,
		"description": "delete data succesfully",
	}
	c.JSON(http.StatusNoContent, gin.H{
		"status": status,
	})
}

func NewCandidateController(r *gin.Engine, candidateUsecase usecase.CandidateUseCase, bootcampUsecase usecase.BootcampUseCase, cloudinary *cloudinary.Cloudinary) *CandidateController {
	controller := CandidateController{
		router:           r,
		candidateUsecase: candidateUsecase,
		bootcampUsecase:  bootcampUsecase,
		cloudinary:       cloudinary,
	}
	rg := r.Group("/api/v1")
	rg.POST("/candidates", controller.createHandler)
	// rg.GET("/candidates", middleware.AuthMiddleware("admin", "hr"), controller.listHandler)
	rg.GET("/candidates", controller.listHandler)
	rg.GET("/candidates/:id", controller.getHandler)
	rg.PUT("/candidates", controller.updateHandler)
	rg.DELETE("/candidates/:id", controller.deleteHandler)
	return &controller
}
