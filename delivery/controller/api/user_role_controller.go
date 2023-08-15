package api

import (
	"interview_bootcamp/delivery/middleware"
	"interview_bootcamp/model"
	"interview_bootcamp/usecase"
	"interview_bootcamp/utils/common"

	"github.com/gin-gonic/gin"
)

type UserRoleController struct {
	router     *gin.Engine
	userRoleUC usecase.UserRolesUseCase
}

func (u *UserRoleController) createHandler(c *gin.Context) {
	var userRole model.UserRoles
	if err := c.ShouldBindJSON(&userRole); err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	userRole.Id = common.GenerateID()
	if err := u.userRoleUC.RegisterNewUserRole(userRole); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	c.JSON(201, userRole)
}

func (u *UserRoleController) listHandler(c *gin.Context) {
	userRoles, err := u.userRoleUC.GetAllUserRoles()
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	c.JSON(200, userRoles)
}

func (u *UserRoleController) getHandler(c *gin.Context) {
	id := c.Param("id")
	userRole, err := u.userRoleUC.GetUserRoleByID(id)
	if err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	c.JSON(200, userRole)
}

func (u *UserRoleController) updateHandler(c *gin.Context) {
	var userRole model.UserRoles
	if err := c.ShouldBindJSON(&userRole); err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	if err := u.userRoleUC.UpdateUserRole(userRole); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}
	c.JSON(200, userRole)

}

func (u *UserRoleController) deleteHandler(c *gin.Context) {
	id := c.Param("id")
	if err := u.userRoleUC.DeleteUserRole(id); err != nil {
		c.JSON(500, gin.H{"err": err.Error()})
		return
	}

	c.String(204, "")
}

func NewUserRoleController(r *gin.Engine, usecase usecase.UserRolesUseCase) *UserRoleController {
	controller := UserRoleController{
		router:     r,
		userRoleUC: usecase,
	}

	rg := r.Group("/api/v1")
	rg.POST("/user-roles", controller.createHandler)
	rg.GET("/user-roles", controller.listHandler)
	rg.GET("/user-roles/:id", controller.getHandler)
	rg.PUT("/user-roles", middleware.AuthMiddleware("admin"), controller.updateHandler)
	rg.DELETE("/user-roles/:id", middleware.AuthMiddleware("admin"), controller.deleteHandler)

	return &controller
}
