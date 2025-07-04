package controller

import (
	"device-communication/src/dto"
	"device-communication/src/dtoError"
	"device-communication/src/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func communicationGroupRouter(g *gin.RouterGroup) {
	group := g.Group("/communication")
	group.Use(GetLoginFilter())
	group.GET("/main", communication.MainDeviceConnection)
	group.GET("/sub", communication.SubDeviceConnection)
}

var communication CommunicationController

type CommunicationController interface {
	MainDeviceConnection(c *gin.Context)
	SubDeviceConnection(c *gin.Context)
}

type communicationControllerImpl struct {
	errWarper     dtoError.ServiceErrorWarpper
	communication service.CommunicationSerivice
}

func init() {
	communication = &communicationControllerImpl{
		errWarper:     dtoError.GetServiceErrorWarpper(),
		communication: service.GetCommunicationSerivice(),
	}
}

func (ctl *communicationControllerImpl) MainDeviceConnection(c *gin.Context) {
	var req dto.MainDeviceConnectionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, id, _ := GetSessionValue(c)
	req.UserId = id
	err := ctl.communication.MainDeviceConnection(c, &req, c.Writer, c.Request)
	if err != nil {
		c.JSON(err.ToJsonResponse())
		return
	}
}

func (ctl *communicationControllerImpl) SubDeviceConnection(c *gin.Context) {
	var req dto.SubDeviceConnectionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, id, _ := GetSessionValue(c)
	req.UserId = id
	err := ctl.communication.SubDeviceConnection(c, &req, c.Writer, c.Request)
	if err != nil {
		c.JSON(err.ToJsonResponse())
		return
	}
}
