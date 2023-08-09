package api

import (
	"interview_bootcamp/model"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
	"net/http"

	"github.com/gin-gonic/gin"
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
	id := c.Param("id")
	uom, err := cc.usecase.FindByIdCandidate(id)
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
		"data":   uom,
	})
}
func (cc *CandidateController) updateHandler(c *gin.Context) {
	panic("")
}
func (cc *CandidateController) deleteHandler(c *gin.Context) {
	id := c.Param("candidate_id")
	if err := cc.usecase.DeleteCandidate(id); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	c.String(204, "")

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
package api

import (
	"interview_bootcamp/model"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
	"net/http"

	"github.com/gin-gonic/gin"
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
	candidates, err := cc.usecase.FindAllCandidate()
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
		"data":   candidates,
	})
}
func (cc *CandidateController) getHandler(c *gin.Context) {
	panic("")
}
func (cc *CandidateController) updateHandler(c *gin.Context) {

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
