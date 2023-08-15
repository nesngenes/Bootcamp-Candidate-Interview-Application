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

type FormController struct {
	router     *gin.Engine
	usecase    usecase.FormUseCase
	cloudinary *cloudinary.Cloudinary
}

func (f *FormController) createHandler(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Error parsing form data"})
		return
	}

	var form model.Form
	err = c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	form.FormID = common.GenerateID()

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
	uploadResult, err := f.cloudinary.Upload.Upload(
		context.Background(), bytes.NewReader(fileContent),
		uploader.UploadParams{
			PublicID: "forms/" + form.FormID,
		},
	)
	if err != nil {
		// Log the error for debugging
		fmt.Println("Error uploading file to cloudinary:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error uploading file to cloudinary"})
		return
	}
	form.FormLink = uploadResult.SecureURL

	err = f.usecase.RegisterNewForm(form)
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

func (f *FormController) listHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	paginationParam := dto.PaginationParam{
		Page:  page,
		Limit: limit,
	}

	forms, paging, err := f.usecase.FindAllForm(paginationParam)
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
		"data":   forms,
		"paging": paging,
	})
}

func (f *FormController) updateHandler(c *gin.Context) {
	// Parse the form data
	err := c.Request.ParseMultipartForm(10 << 20) // Max 10MB file size
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Error parsing form data"})
		return
	}

	// Retrieve the form fields
	formID := c.PostForm("id")

	// Fetch the existing form by ID
	existingForm, err := f.usecase.FindByIdForm(formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

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
		uploadResult, err := f.cloudinary.Upload.Upload(
			context.Background(), bytes.NewReader(fileContent),
			uploader.UploadParams{
				PublicID: "forms/" + formID,
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "Error uploading file to cloudinary"})
			return
		}
		existingForm.FormLink = uploadResult.SecureURL
	}

	// Update the form in the database
	err = f.usecase.UpdateForm(existingForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error updating form"})
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

func (f *FormController) deleteHandler(c *gin.Context) {
	id := c.Param("id")

	if err := f.usecase.DeleteForm(id); err != nil {
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

func NewFormController(r *gin.Engine, usecase usecase.FormUseCase, cloudinary *cloudinary.Cloudinary) *FormController {
	controller := FormController{
		router:     r,
		usecase:    usecase, // Change this line
		cloudinary: cloudinary,
	}
	rg := r.Group("/api/v1")
	// rg.POST("/forms", middleware.AuthMiddleware("admin", "interviewer"), controller.createHandler)
	// rg.GET("/forms", middleware.AuthMiddleware("admin", "interviewer", "hr_recruitment"), controller.listHandler)
	// rg.DELETE("/forms/:id", middleware.AuthMiddleware("admin", "interviewer"), controller.deleteHandler)
	// rg.PUT("/forms", middleware.AuthMiddleware("admin", "interviewer"), controller.updateHandler)

	rg.POST("/forms", controller.createHandler)
	rg.GET("/forms", controller.listHandler)
	rg.DELETE("/forms/:id", controller.deleteHandler)
	rg.PUT("/forms", controller.updateHandler)
	return &controller
}
