package api

import (
	"interview_bootcamp/model"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	router *gin.Engine
	userUC usecase.UserUsecase
}

func (c *UserController) createHandler(ctx *gin.Context) {
	var user model.Users
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Id = common.GenerateID()
	err := c.userUC.RegisterNewUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) listHandler(ctx *gin.Context) {
	users, err := c.userUC.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) getHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.userUC.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) getByUsernameHandler(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := c.userUC.GetUserByUserName(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) updateHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var user model.Users
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Id = id
	err := c.userUC.UpdateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) deleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.userUC.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func NewUserController(r *gin.Engine, usecase usecase.UserUsecase) *UserController {
	controller := UserController{
		router: r,
		userUC: usecase,
	}
	// Register routes
	rg := r.Group("/api/v1")
	rg.POST("/users", controller.createHandler)
	rg.GET("/users", controller.listHandler)
	rg.GET("/users/:id", controller.getHandler)
	rg.GET("/users/by-username/:username", controller.getByUsernameHandler)
	rg.PUT("/users/:id", controller.updateHandler) //bisa update asal role id and role name tetep sama

	rg.DELETE("/users/:id", controller.deleteHandler)

	return &controller
}
