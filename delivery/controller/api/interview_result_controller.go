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

type InterviewResultController struct {
	router            *gin.Engine
	interviewResultUC usecase.InterviewResultUseCase
}

func (i *InterviewResultController) createHandler(c *gin.Context) {
	var interviewResult model.InterviewResult
	if err := c.ShouldBindJSON(&interviewResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	interviewResult.Id = common.GenerateID()
	if err := i.interviewResultUC.RegisterNewInterviewResult(interviewResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, interviewResult)
}
func (i *InterviewResultController) listHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	paginationParam := dto.PaginationParam{
		Page:  page,
		Limit: limit,
	}
	InterviewsP, paging, err := i.interviewResultUC.FindAllInterviewResult(paginationParam)
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
		"data":   InterviewsP,
		"paging": paging,
	})
}
func (i *InterviewResultController) getHandler(c *gin.Context) {
	id := c.Param("id")
	interviewResult, err := i.interviewResultUC.FindByIdInterviewResult(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	status := map[string]any{
		"code":        200,
		"description": "Get By Id Data Successfully",
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"data":   interviewResult,
	})
}

func NewInterviewResultController(r *gin.Engine, usecase usecase.InterviewResultUseCase) *InterviewResultController {
	controller := InterviewResultController{
		router:            r,
		interviewResultUC: usecase,
	}

	rg := r.Group("/api/v1")
	rg.POST("/interviewresult", controller.createHandler)
	rg.GET("/interviewresult", controller.listHandler)
	rg.GET("/interviewresult/:id", controller.getHandler)
	return &controller
}
