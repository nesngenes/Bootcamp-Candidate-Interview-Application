package api

import (
	"interview_bootcamp/model"
	"interview_bootcamp/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	router  *gin.Engine
	usecase usecase.AuthUseCAse
}

func (a *AuthController) loginHandler(c *gin.Context) {
	var payload model.Users
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	token, err := a.usecase.Login(payload.UserName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, token)
}

func NewAuthUseCase(e *gin.Engine, usecase usecase.AuthUseCAse) *AuthController {
	controller := AuthController{
		router:  e,
		usecase: usecase,
	}

	rGroup := e.Group("api/v1")
	rGroup.POST("login", controller.loginHandler)
	return &controller
}
