package api

import (
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

func (cc *InterviewerController) createHandler(c *gin.Context) {
	//reminder: cek lagi,kayaknya masih harus di rombak createnya
	var interviewer model.Interviewer
	if err := c.ShouldBindJSON(&interviewer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	interviewer.InterviewerID = common.GenerateID()

	if err := cc.usecase.RegisterNewInterviewer(interviewer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, interviewer)
}

func (cc *InterviewerController) listHandler(c *gin.Context) {
	panic("")
}
func (cc *InterviewerController) getHandler(c *gin.Context) {
	id := c.Param("interviewer_id")
	interviewer, err := cc.usecase.FindByIdInterviewer(id)
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
func (cc *InterviewerController) updateHandler(c *gin.Context) {
	panic("")
}
func (cc *InterviewerController) deleteHandler(c *gin.Context) {
	id := c.Param("candidate_id")
	if err := cc.usecase.DeleteInterviewer(id); err != nil {
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
	rg.POST("/interviewers", controller.createHandler)
	rg.GET("/interviewers", controller.listHandler)
	rg.GET("/interviewers/:id", controller.getHandler)
	rg.PUT("/interviewers", controller.updateHandler)
	rg.DELETE("/interviewers/:id", controller.deleteHandler)
	return &controller
}
