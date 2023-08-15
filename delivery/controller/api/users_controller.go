package api

import (
	"interview_bootcamp/delivery/middleware"
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

	// Create the user response object
	userResponse := map[string]interface{}{
		"id":        user.Id,
		"email":     user.Email,
		"user_name": user.UserName,
		"user_role": map[string]interface{}{
			"id":   user.UserRole.Id,
			"name": user.UserRole.Name,
		},
	}

	ctx.JSON(http.StatusCreated, userResponse)
}

func (u *UserController) listHandler(c *gin.Context) {
	users, err := u.userUC.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	status := map[string]interface{}{
		"code":        200,
		"description": "Get All Data Successfully",
	}

	userResponses := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		userResponse := map[string]interface{}{
			"id":        user.Id,
			"email":     user.Email,
			"user_name": user.UserName,
			"user_role": map[string]interface{}{
				"id":   user.UserRole.Id,
				"name": user.UserRole.Name,
			},
		}
		userResponses = append(userResponses, userResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"data":   userResponses,
	})
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
	username := ctx.Param("user_name")
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
	rg.GET("/users", middleware.AuthMiddleware("admin"), controller.listHandler)
	rg.GET("/users/:id", middleware.AuthMiddleware("admin", "hr_recruitment", "interviewer"), controller.getHandler)
	rg.GET("/users/by-username/:username", middleware.AuthMiddleware("admin"), controller.getByUsernameHandler)
	rg.PUT("/users/:id", middleware.AuthMiddleware("admin"), controller.updateHandler) //bisa update asal role id and role name tetep sama
	rg.DELETE("/users/:id", middleware.AuthMiddleware("admin"), controller.deleteHandler)

	return &controller
}
