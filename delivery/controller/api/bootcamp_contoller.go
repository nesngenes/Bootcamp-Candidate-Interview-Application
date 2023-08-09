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

type BootcampController struct {
	router  *gin.Engine
	usecase usecase.BootcampUseCase
}

func (b *BootcampController) createHandler(c *gin.Context) {
	var bootcamp model.Bootcamp
	if err := c.ShouldBindJSON(&bootcamp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	bootcamp.BootcampId = common.GenerateID()
	if err := b.usecase.RegisterNewBootcamp(bootcamp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, bootcamp)
}
func (b *BootcampController) listHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	paginationParam := dto.PaginationParam{
		Page:  page,
		Limit: limit,
	}
	bootcamps, paging, err := b.usecase.FindAllBootcamp(paginationParam)
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
		"data":   bootcamps,
		"paging": paging,
	})

}
func (b *BootcampController) getHandler(c *gin.Context) {
	id := c.Param("id")
	bootcamps, err := b.usecase.FindByIdBootcamp(id)
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
		"data":   bootcamps,
	})

}
func (b *BootcampController) updateHandler(c *gin.Context) {
	var bootcamp model.Bootcamp
	if err := c.ShouldBindJSON(&bootcamp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := b.usecase.UpdateBootcamp(bootcamp); err != nil {
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
func (b *BootcampController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := b.usecase.DeleteBootcamp(id); err != nil {
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

func NewBootcampController(r *gin.Engine, usecase usecase.BootcampUseCase) *BootcampController {
	controller := BootcampController{
		router:  r,
		usecase: usecase,
	}
	rg := r.Group("/api/v1")
	rg.POST("/bootcamps", controller.createHandler)
	rg.GET("/bootcamps", controller.listHandler)
	rg.GET("/bootcamps/:id", controller.getHandler)
	rg.PUT("/bootcamps", controller.updateHandler)
	rg.DELETE("/bootcamps/:id", controller.deleteHandler)
	return &controller
}
