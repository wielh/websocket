package controller

import (
	"device-communication/src/dto"
	"device-communication/src/dtoError"
	"device-communication/src/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var device DeviceController

func deviceGroupRouter(g *gin.RouterGroup) {
	group := g.Group("/device")
	group.Use(GetLoginFilter())
	group.PUT("/main", device.BindMainDevice)
	group.DELETE("/main", device.UnBindMainDevice)
	group.PUT("/sub", device.BindSubDevice)
	group.DELETE("/sub", device.UnBindSubDevice)
	group.GET("/", device.GetDevicesByUserId)
}

type DeviceController interface {
	BindMainDevice(c *gin.Context)
	UnBindMainDevice(c *gin.Context)
	BindSubDevice(c *gin.Context)
	UnBindSubDevice(c *gin.Context)
	GetDevicesByUserId(c *gin.Context)
}

type deviceControllerImpl struct {
	errWarper     dtoError.ServiceErrorWarpper
	deviceService service.DeviceService
}

func (d *deviceControllerImpl) BindMainDevice(c *gin.Context) {
	var req dto.BindMainDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		serviceErr := d.errWarper.NewParseParametersFailedError(err)
		c.JSON(serviceErr.ToJsonResponse())
		return
	}

	_, id, _ := GetSessionValue(c)
	req.UserId = id
	res, serviceErr := d.deviceService.BindMainDevice(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.ToJsonResponse())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (d *deviceControllerImpl) BindSubDevice(c *gin.Context) {
	var req dto.BindSubDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		serviceErr := d.errWarper.NewParseParametersFailedError(err)
		c.JSON(serviceErr.ToJsonResponse())
		return
	}

	_, id, _ := GetSessionValue(c)
	req.UserId = id
	res, serviceErr := d.deviceService.BindSubDevice(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.ToJsonResponse())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (d *deviceControllerImpl) GetDevicesByUserId(c *gin.Context) {
	_, id, _ := GetSessionValue(c)
	req := dto.GetDevicesByUserIdRequest{UserId: id}
	res, serviceErr := d.deviceService.GetDevicesByUserId(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.ToJsonResponse())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (d *deviceControllerImpl) UnBindMainDevice(c *gin.Context) {
	var req dto.UnbindMainDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		serviceErr := d.errWarper.NewParseParametersFailedError(err)
		c.JSON(serviceErr.ToJsonResponse())
		return
	}

	_, id, _ := GetSessionValue(c)
	req.UserId = id
	res, serviceErr := d.deviceService.UnBindMainDevice(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.ToJsonResponse())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (d *deviceControllerImpl) UnBindSubDevice(c *gin.Context) {
	var req dto.UnbindSubDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		serviceErr := d.errWarper.NewParseParametersFailedError(err)
		c.JSON(serviceErr.ToJsonResponse())
		return
	}

	_, id, _ := GetSessionValue(c)
	req.UserId = id
	res, serviceErr := d.deviceService.UnBindSubDevice(c, &req)
	if serviceErr != nil {
		c.JSON(serviceErr.ToJsonResponse())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": res})
}

func init() {
	device = &deviceControllerImpl{
		errWarper:     dtoError.GetServiceErrorWarpper(),
		deviceService: service.GetDeviceService(),
	}
}
