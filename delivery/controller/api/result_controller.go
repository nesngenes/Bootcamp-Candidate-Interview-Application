package api

import (
	"interview_bootcamp/delivery/middleware"
	"interview_bootcamp/model"
	"interview_bootcamp/model/dto"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResultController struct {
	router  *gin.Engine
	usecase usecase.ResultUseCase
}

func (r *ResultController) createHandler(c *gin.Context) {
	var result model.Result
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	result.ResultId = common.GenerateID()
	if err := r.usecase.RegisterNewResult(result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}
func (r *ResultController) listHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	paginationParam := dto.PaginationParam{
		Page:  page,
		Limit: limit,
	}
	results, paging, err := r.usecase.FindAllResult(paginationParam)
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
		"data":   results,
		"paging": paging,
	})

}
func (r *ResultController) getHandler(c *gin.Context) {
	id := c.Param("id")
	results, err := r.usecase.FindByIdResult(id)
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
		"data":   results,
	})
}
func (r *ResultController) updateHandler(c *gin.Context) {
	var results model.Result
	if err := c.ShouldBindJSON(&results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := r.usecase.UpdateResult(results); err != nil {
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
func (r *ResultController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := r.usecase.DeleteResult(id); err != nil {
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

func NewResultController(r *gin.Engine, usecase usecase.ResultUseCase) *ResultController {
	controller := ResultController{
		router:  r,
		usecase: usecase,
	}
	rg := r.Group("/api/v1")
	rg.POST("/results", middleware.AuthMiddleware("admin", "hr_recruitment", "interviewer"), controller.createHandler)
	rg.GET("/results", middleware.AuthMiddleware("admin", "hr_recruitment", "interviewer"), controller.listHandler)
	rg.GET("/results/:id", middleware.AuthMiddleware("admin", "hr_recruitment", "interviewer"), controller.getHandler)
	rg.PUT("/results", middleware.AuthMiddleware("admin", "hr_recruitment", "interviewer"), controller.updateHandler)
	rg.DELETE("/results/:id", middleware.AuthMiddleware("admin", "hr_recruitment", "interviewer"), controller.deleteHandler)
	return &controller
}
