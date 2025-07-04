package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"device-communication/src/dto"
	"device-communication/src/dtoError"
	"device-communication/src/service"
)

func userGroupRouter(g *gin.RouterGroup) {
	group := g.Group("/user")
	group.POST("/register", user.Register)
	group.POST("/login", user.Login)
	group.PUT("/reset_password", user.ResetPassword)
}

type UserController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	ResetPassword(c *gin.Context)
}

type userControllerImpl struct {
	errWarper   dtoError.ServiceErrorWarpper
	userService service.UserService
}

var user UserController

func init() {
	user = &userControllerImpl{
		errWarper:   dtoError.GetServiceErrorWarpper(),
		userService: service.GetUserService(),
	}
}

func (u *userControllerImpl) Register(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		serviceErr := u.errWarper.NewParseParametersFailedError(err)
		c.JSON(serviceErr.ToJsonResponse())
		return
	}

	res, serviceErr := u.userService.UserRegisterService(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.ToJsonResponse())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (u *userControllerImpl) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		serviceErr := u.errWarper.NewParseParametersFailedError(err)
		c.JSON(serviceErr.ToJsonResponse())
		return
	}

	res, serviceErr := service.GetUserService().UserLoginService(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.ToJsonResponse())
		return
	}

	_, err := SetSessionValue(c, res.ID, res.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (u *userControllerImpl) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		serviceErr := u.errWarper.NewParseParametersFailedError(err)
		c.JSON(serviceErr.ToJsonResponse())
		return
	}

	serviceErr := service.GetUserService().ResetPasswordService(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.ToJsonResponse())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
