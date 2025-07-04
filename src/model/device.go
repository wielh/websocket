package model

type MainDevice struct {
	Id         uint64      `gorm:"primaryKey;column:id"`
	UserId     uint64      `gorm:"not null;column:user_id"`
	Platform   string      `gorm:"not null;column:platform"`
	Version    string      `gorm:"not null;column:version"`
	DeviceId   string      `gorm:"not null;column:device_id"`
	SubDevices []SubDevice `gorm:"foreignKey:MainDeviceId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Base
}

type SubDevice struct {
	Id           uint64 `gorm:"primaryKey;column:id"`
	MainDeviceId uint64 `gorm:"not null;column:main_device_id"`
	Platform     string `gorm:"not null;column:platform"`
	Version      string `gorm:"not null;column:version"`
	DeviceId     string `gorm:"not null;column:device_id"`
	Base
}
