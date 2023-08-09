package api

import (
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
	"net/http"
	"strconv"

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

	candidates, paging, err := cc.usecase.FindAllCandidate(paginationParam)
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
	product, err := cc.usecase.FindByIdCandidate(id)
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

	var candidate model.Candidate
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := cc.usecase.UpdateCandidate(candidate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	status := map[string]any{
		"code":        http.StatusOK,
		"description": "update data succesfully",
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
func (cc *CandidateController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := cc.usecase.DeleteCandidate(id); err != nil {
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
