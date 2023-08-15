package api

import (
	"interview_bootcamp/delivery/middleware"
	"interview_bootcamp/model"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InterviewerController struct {
	router  *gin.Engine
	usecase usecase.InterviewerUseCase
}

func (i *InterviewerController) createHandler(c *gin.Context) {
	var interviewer model.Interviewer
	if err := c.ShouldBindJSON(&interviewer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	interviewer.InterviewerID = common.GenerateID()
	if err := i.usecase.RegisterNewInterviewer(interviewer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, interviewer)
}

func (i *InterviewerController) listHandler(c *gin.Context) {
	interviewers, err := i.usecase.FindAllInterviewer()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	status := map[string]any{
		"code":        200,
		"description": "get all data succesfully",
	}
	c.JSON(200, gin.H{
		"status": status,
		"data":   interviewers,
	})
}

func (i *InterviewerController) getHandler(c *gin.Context) {
	id := c.Param("id")
	interviewer, err := i.usecase.FindByIdInterviewer(id)
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	status := map[string]any{
		"code":        200,
		"description": "Get By Id Data Successfully",
	}
	c.JSON(200, gin.H{
		"status": status,
		"data":   interviewer,
	})

}

func (i *InterviewerController) updateHandler(c *gin.Context) {
	var interviewer model.Interviewer
	if err := c.ShouldBindJSON(&interviewer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Ensure that InterviewerID is provided in the JSON payload
	if interviewer.InterviewerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "InterviewerID is required for updating"})
		return
	}

	if err := i.usecase.UpdateInterviewer(interviewer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, interviewer)
}

func (i *InterviewerController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := i.usecase.DeleteInterviewer(id); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	c.String(204, "")

}

func NewInterviewerController(r *gin.Engine, usecase usecase.InterviewerUseCase) *InterviewerController {
	controller := InterviewerController{
		router:  r,
		usecase: usecase,
	}
	rg := r.Group("/api/v1")
	rg.POST("/interviewers", middleware.AuthMiddleware("admin", "interviewer"), controller.createHandler)
	rg.GET("/interviewers", middleware.AuthMiddleware("admin", "interviewer", "hr_recruitmen"), controller.listHandler)
	rg.GET("/interviewers/:id", middleware.AuthMiddleware("admin", "hr_recruitment", "interviewer"), controller.getHandler)
	rg.PUT("/interviewers", middleware.AuthMiddleware("admin", "interviewer"), controller.updateHandler)
	rg.DELETE("/interviewers/:id", middleware.AuthMiddleware("admin", "interviewer"), controller.deleteHandler)
	return &controller
}
