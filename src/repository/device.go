package repository

import (
	"context"
	"device-communication/src/config"
	"device-communication/src/model"

	"gorm.io/gorm"
)

type DeviceRepository interface {
	GetAllDevicesByUserId(ctx context.Context, userId uint64) ([]*model.MainDevice, error)
	GetMainDeviceCount(ctx context.Context, userId uint64) (int64, error)
	GetSubDeviceCount(ctx context.Context, userId uint64, mainDeviceId uint64) (int64, error)
	CheckRepeatedDevice(ctx context.Context, platform string, version string, deviceId string) (bool, error)
	CheckMainDeviceBinding(ctx context.Context, userId uint64, mainDeviceId uint64) (bool, error)
	CheckSubDeviceBinding(ctx context.Context, userId uint64, mainDeviceId uint64, subDeviceId uint64) (bool, error)
	BindMainDevice(ctx context.Context, userId uint64, platform string, version string, deviceId string) (*model.MainDevice, error)
	UnbindMainDevice(ctx context.Context, userId uint64, mainDeviceId uint64) (bool, error)
	BindSubDevice(ctx context.Context, mainDeviceId uint64, platform string, version string, deviceId string) (*model.SubDevice, error)
	UnbindSubDevice(mctx context.Context, mainDeviceId uint64, subDeviceId uint64) error
}

type deviceRepositoryImpl struct {
	DB *gorm.DB
}

func (d *deviceRepositoryImpl) GetAllDevicesByUserId(ctx context.Context, userId uint64) ([]*model.MainDevice, error) {
	tx := GetTxContext(ctx, d.DB)
	var devices []*model.MainDevice
	result := tx.Preload("SubDevices", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "device_id", "platform", "version", "main_device_id")
	}).Where("user_id = ?", userId).Find(&devices)

	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (d *deviceRepositoryImpl) GetMainDeviceCount(ctx context.Context, userId uint64) (int64, error) {
	tx := GetTxContext(ctx, d.DB)
	var count int64
	err := tx.Model(&model.MainDevice{}).Where("user_id=?", userId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *deviceRepositoryImpl) GetSubDeviceCount(ctx context.Context, userId uint64, mainDeviceId uint64) (int64, error) {
	tx := GetTxContext(ctx, d.DB)
	var count int64
	err := tx.Model(&model.SubDevice{}).
		Joins("JOIN main_devices ON main_devices.id = sub_devices.main_device_id").
		Where("main_devices.user_id = ? AND main_devices.id = ?", userId, mainDeviceId).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *deviceRepositoryImpl) CheckRepeatedDevice(ctx context.Context, platform string, version string, deviceId string) (bool, error) {
	tx := GetTxContext(ctx, d.DB)
	var count int64
	err := tx.Model(&model.MainDevice{}).Where("platform=? and version=? and device_id=?", platform, version, deviceId).Count(&count).Error
	if err != nil {
		return false, err
	} else if count > 0 {
		return false, nil
	}

	err = tx.Model(&model.SubDevice{}).Where("platform=? and version=? and device_id=?", platform, version, deviceId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (d *deviceRepositoryImpl) BindMainDevice(ctx context.Context, userId uint64, platform string, version string, deviceId string) (*model.MainDevice, error) {
	tx := GetTxContext(ctx, d.DB)
	mainDevice := model.MainDevice{
		UserId:   userId,
		Platform: platform,
		Version:  version,
		DeviceId: deviceId,
	}

	result := tx.Create(&mainDevice)
	if result.Error != nil {
		return nil, result.Error
	}
	return &mainDevice, nil
}

func (d *deviceRepositoryImpl) BindSubDevice(ctx context.Context, mainDeviceId uint64, platform string, version string, deviceId string) (*model.SubDevice, error) {
	tx := GetTxContext(ctx, d.DB)
	subDevice := model.SubDevice{
		MainDeviceId: mainDeviceId,
		Platform:     platform,
		Version:      version,
		DeviceId:     deviceId,
	}

	result := tx.Create(&subDevice)
	if result.Error != nil {
		return nil, result.Error
	}
	return &subDevice, nil
}

func (d *deviceRepositoryImpl) CheckMainDeviceBinding(ctx context.Context, userId uint64, mainDeviceId uint64) (bool, error) {
	tx := GetTxContext(ctx, d.DB)
	var count int64
	result := tx.Model(&model.MainDevice{}).Where("user_id = ? AND id = ?", userId, mainDeviceId).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (d *deviceRepositoryImpl) CheckSubDeviceBinding(ctx context.Context, userId uint64, mainDeviceId uint64, subDeviceId uint64) (bool, error) {
	tx := GetTxContext(ctx, d.DB)
	var count int64
	result := tx.Table("main_devices m").Joins("JOIN sub_devices s ON m.id = s.main_device_id").
		Where("m.user_id = ? AND m.id=? AND s.id= ?", userId, mainDeviceId, subDeviceId).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (d *deviceRepositoryImpl) UnbindMainDevice(ctx context.Context, userId uint64, mainDeviceId uint64) (bool, error) {
	tx := GetTxContext(ctx, d.DB)
	result := tx.Where("main_device_id = ?", userId, mainDeviceId).Delete(&model.SubDevice{})
	if result.Error != nil {
		return false, result.Error
	}

	result = tx.Where("user_id = ? AND id = ?", userId, mainDeviceId).Delete(&model.MainDevice{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (d *deviceRepositoryImpl) UnbindSubDevice(ctx context.Context, mainDeviceId uint64, subDeviceId uint64) error {
	tx := GetTxContext(ctx, d.DB)
	result := tx.Where("main_device_id = ? AND id = ?", mainDeviceId, subDeviceId).Delete(&model.SubDevice{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

var device DeviceRepository

func init() {
	device = &deviceRepositoryImpl{
		DB: config.GlobalConfig.DB,
	}
}

func GetDeviceRepository() DeviceRepository {
	return device
}
