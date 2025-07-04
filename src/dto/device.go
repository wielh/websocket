package dto

type BindMainDeviceRequest struct {
	UserId   uint64
	Platform string `json:"platform" binding:"required"`
	Version  string `json:"version" binding:"required"`
	DeviceId string `json:"device_id" binding:"required"`
}

type BindMainDeviceResponse struct {
	Ok           bool   `json:"ok"`
	MainDeviceId uint64 `json:"main_device_id"`
}

type UnbindMainDeviceRequest struct {
	UserId       uint64
	MainDeviceId uint64 `json:"main_device_id" binding:"required"`
	SubDeviceId  uint64 `json:"sub_device_id" binding:"required"`
}

type UnbindMainDeviceResponse struct {
	Ok bool `json:"ok"`
}

type BindSubDeviceRequest struct {
	UserId       uint64
	MainDeviceId uint64 `json:"main_device_id" binding:"required"`
	Platform     string `json:"platform" binding:"required"`
	Version      string `json:"version" binding:"required"`
	DeviceId     string `json:"device_id" binding:"required"`
}

type BindSubDeviceResponse struct {
	Ok          bool   `json:"ok"`
	SubDeviceId uint64 `json:"sub_device_id"`
}

type UnbindSubDeviceRequest struct {
	UserId       uint64
	MainDeviceId uint64 `json:"main_device_id" binding:"required"`
	SubDeviceId  uint64 `json:"sub_device_id" binding:"required"`
}

type UnbindSubDeviceResponse struct {
}

type GetDevicesByUserIdRequest struct {
	UserId uint64
}

type GetDevicesByUserIdResponse struct {
	MainDevices []*MainDevice `json:"main_devices"`
}

type MainDevice struct {
	Id         uint64       `json:"id"`
	UserId     uint64       `json:"user_id"`
	Platform   string       `json:"platform"`
	Version    string       `json:"version"`
	DeviceId   string       `json:"device_id"`
	SubDevices []*SubDevice `json:"sub_devices,omitempty"`
}

type SubDevice struct {
	Id           uint64 `json:"id"`
	MainDeviceId uint64 `json:"main_device_id"`
	Platform     string `json:"platform"`
	Version      string `json:"version"`
	DeviceId     string `json:"device_id"`
}
