package dto

type MainDeviceConnectionRequest struct {
	UserId       uint64 `binding:"-"`
	MainDeviceId uint64 `form:"main_device_id" binding:"required"`
}

type SubDeviceConnectionRequest struct {
	UserId       uint64 `binding:"-"`
	MainDeviceId uint64 `form:"main_device_id" binding:"required"`
	SubDeviceId  uint64 `form:"sub_device_id" binding:"required"`
}
