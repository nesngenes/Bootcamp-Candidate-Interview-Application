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

type InterviewProcessController struct {
	router             *gin.Engine
	interviewProcessUC usecase.InterviewProcessUseCase
}

func (i *InterviewProcessController) createHandler(c *gin.Context) {
	var interviewProcess model.InterviewProcess
	if err := c.ShouldBindJSON(&interviewProcess); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	interviewProcess.ID = common.GenerateID()
	if err := i.interviewProcessUC.RegisterNewInterviewProcess(interviewProcess); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, interviewProcess)
}

func (i *InterviewProcessController) listHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	paginationParam := dto.PaginationParam{
		Page:  page,
		Limit: limit,
	}
	InterviewsP, paging, err := i.interviewProcessUC.FindAllInterviewProcess(paginationParam)
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
func (i *InterviewProcessController) getHandler(c *gin.Context) {
	id := c.Param("id")
	interviewP, err := i.interviewProcessUC.FindByIdInterviewProcess(id)
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
		"data":   interviewP,
	})
}

func NewInterviewProcessController(r *gin.Engine, usecase usecase.InterviewProcessUseCase) *InterviewProcessController {
	controller := InterviewProcessController{
		router:             r,
		interviewProcessUC: usecase,
	}

	rg := r.Group("/api/v1")
	rg.POST("/interviewprocess", middleware.AuthMiddleware("admin", "hr_recruitment"), controller.createHandler)
	rg.GET("/interviewprocess", middleware.AuthMiddleware("admin", "hr_recruitment", "interviewer"), controller.listHandler)
	rg.GET("/interviewprocess/:id", middleware.AuthMiddleware("admin", "hr_recruitment", "interviewer"), controller.getHandler)
	return &controller
}
