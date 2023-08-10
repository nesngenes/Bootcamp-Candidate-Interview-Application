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

type StatusController struct {
	router  *gin.Engine
	usecase usecase.StatusUseCase
}

func (s *StatusController) createHandler(c *gin.Context) {
	var status model.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	status.StatusId = common.GenerateID()
	if err := s.usecase.RegisterNewStatus(status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, status)
}
func (s *StatusController) listHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	paginationParam := dto.PaginationParam{
		Page:  page,
		Limit: limit,
	}
	statuss, paging, err := s.usecase.FindAllStatus(paginationParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	status := map[string]any{
		"code":        200,
		"description": "Get All Data Successfully",
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"data":   statuss,
		"paging": paging,
	})

}
func (s *StatusController) getHandler(c *gin.Context) {
	id := c.Param("id")
	statuss, err := s.usecase.FindByIdStatus(id)
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
		"data":   statuss,
	})
}
func (s *StatusController) updateHandler(c *gin.Context) {
	var statuss model.Status
	if err := c.ShouldBindJSON(&statuss); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := s.usecase.UpdateStatus(statuss); err != nil {
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
func (s *StatusController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := s.usecase.DeleteStatus(id); err != nil {
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

func NewStatusController(r *gin.Engine, usecase usecase.StatusUseCase) *StatusController {
	controller := StatusController{
		router:  r,
		usecase: usecase,
	}
	rg := r.Group("/api/v1")
	rg.POST("/statuss", controller.createHandler)
	rg.GET("/statuss", controller.listHandler)
	rg.GET("/statuss/:id", controller.getHandler)
	rg.PUT("/statuss", controller.updateHandler)
	rg.DELETE("/statuss/:id", controller.deleteHandler)
	return &controller
}
