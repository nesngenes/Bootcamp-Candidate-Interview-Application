package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"interview_bootcamp/model"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
)

type CandidateController struct {
	router  *gin.Engine
	usecase usecase.CandidateUseCase
}

func (cc *CandidateController) createHandler(c *gin.Context) {
	var candidate model.Candidate
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	candidate.CandidateID = common.GenerateID()
	if err := cc.usecase.RegisterNewCandidate(candidate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, candidate)
}

func (cc *CandidateController) listHandler(c *gin.Context) {
	panic("")
}
func (cc *CandidateController) getHandler(c *gin.Context) {
	panic("")
}
func (cc *CandidateController) updateHandler(c *gin.Context) {
	panic("")
}
func (cc *CandidateController) deleteHandler(c *gin.Context) {
	panic("")
}

func NewCandidateController(r *gin.Engine, usecase usecase.CandidateUseCase) *CandidateController {
	controller := CandidateController{
		router:  r,
		usecase: usecase,
	}
	rg := r.Group("/api/v1")
	rg.POST("/candidates", controller.createHandler)
	rg.GET("/candidates", controller.listHandler)
	rg.GET("/candidates/:id", controller.getHandler)
	rg.PUT("/candidates", controller.updateHandler)
	rg.DELETE("/candidates/:id", controller.deleteHandler)
	return &controller
}
