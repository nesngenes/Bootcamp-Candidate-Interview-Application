package api

import (
	"interview_bootcamp/model"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HRRecruitmentController struct {
	router               *gin.Engine
	hrRecruitmentUsecase usecase.HRRecruitmentUsecase
}

// tambahan buat router
func (c *HRRecruitmentController) GetRouter() *gin.Engine {
	return c.router
}

func (c *HRRecruitmentController) createHandler(ctx *gin.Context) {
	var hrRecruitment model.HRRecruitment
	if err := ctx.ShouldBindJSON(&hrRecruitment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the HR recruitment record already exists
	_, err := c.hrRecruitmentUsecase.Get(hrRecruitment.ID)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "HR recruitment record already exists"})
		return
	}
	hrRecruitment.ID = common.GenerateID()
	err = c.hrRecruitmentUsecase.CreateHRRecruitment(hrRecruitment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, hrRecruitment)
}

func (c *HRRecruitmentController) updateHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var hrRecruitment model.HRRecruitment
	if err := ctx.ShouldBindJSON(&hrRecruitment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hrRecruitment.ID = id //perubahan
	existingHRRecruitment, err := c.hrRecruitmentUsecase.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "HR recruitment record not found"})
		return
	}

	if existingHRRecruitment.ID != hrRecruitment.ID {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Cannot update to a different HR recruitment record ID"})
		return
	}
	err = c.hrRecruitmentUsecase.UpdateHRRecruitment(hrRecruitment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, hrRecruitment)
}

func (c *HRRecruitmentController) getHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	hrRecruitment, err := c.hrRecruitmentUsecase.Get(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, hrRecruitment)
}

func (c *HRRecruitmentController) listHandler(ctx *gin.Context) {
	hrRecruitments, err := c.hrRecruitmentUsecase.ListHRRecruitments()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, hrRecruitments)
}

func (c *HRRecruitmentController) deleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.hrRecruitmentUsecase.DeleteHRRecruitment(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func NewHRRecruitmentController(r *gin.Engine, usecase usecase.HRRecruitmentUsecase) *HRRecruitmentController {
	controller := HRRecruitmentController{
		router:               r,
		hrRecruitmentUsecase: usecase,
	}
	//routernya kumpul sini
	rg := r.Group("/api/v1")
	rg.POST("/hr-recruitment", controller.createHandler)
	rg.GET("/hr-recruitment", controller.listHandler)
	rg.PUT("/hr-recruitment/:id", controller.updateHandler)
	rg.GET("/hr-recruitment/:id", controller.getHandler)
	rg.DELETE("/hr-recruitment/:id", controller.deleteHandler)
	return &controller
}
