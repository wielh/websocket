package service

import (
	"context"
	"device-communication/src/dto"
	"device-communication/src/dtoError"
	logger "device-communication/src/log"
	"device-communication/src/repository"
)

type DeviceService interface {
	BindMainDevice(ctx context.Context, req *dto.BindMainDeviceRequest) (*dto.BindMainDeviceResponse, *dtoError.ServiceError)
	UnBindMainDevice(ctx context.Context, req *dto.UnbindMainDeviceRequest) (*dto.UnbindMainDeviceResponse, *dtoError.ServiceError)
	BindSubDevice(ctx context.Context, req *dto.BindSubDeviceRequest) (*dto.BindSubDeviceResponse, *dtoError.ServiceError)
	UnBindSubDevice(ctx context.Context, req *dto.UnbindSubDeviceRequest) (*dto.UnbindSubDeviceResponse, *dtoError.ServiceError)
	GetDevicesByUserId(ctx context.Context, req *dto.GetDevicesByUserIdRequest) (*dto.GetDevicesByUserIdResponse, *dtoError.ServiceError)
}

type deviceServiceImpl struct {
	userRepo              repository.UserRepository
	deviceRepo            repository.DeviceRepository
	errWarpper            dtoError.ServiceErrorWarpper
	MAX_MAIN_DEVICE_COUNT int64
	MAX_SUB_DEVICE_COUNT  int64
	logger                logger.Logger
}

var device DeviceService

func init() {
	device = &deviceServiceImpl{
		userRepo:              repository.GetuserRepository(),
		deviceRepo:            repository.GetDeviceRepository(),
		errWarpper:            dtoError.GetServiceErrorWarpper(),
		MAX_MAIN_DEVICE_COUNT: 1,
		MAX_SUB_DEVICE_COUNT:  1,
		logger:                logger.NewInfoLogger(),
	}
}

func (d *deviceServiceImpl) BindMainDevice(ctx context.Context, req *dto.BindMainDeviceRequest) (*dto.BindMainDeviceResponse, *dtoError.ServiceError) {
	txContext, tx := repository.SetTxContext(ctx)
	ok, err := d.deviceRepo.CheckRepeatedDevice(txContext, req.Platform, req.Version, req.DeviceId)
	if err != nil {
		d.logger.Error("", "d.deviceRepo.CheckRepeatedDevice", req, err)
		tx.Rollback()
		return nil, d.errWarpper.NewDBServiceError(err)
	} else if !ok {
		d.logger.Info("", "d.deviceRepo.CheckRepeatedDevice", req, nil)
		tx.Rollback()
		return nil, d.errWarpper.NewRepeatDeviceError()
	}

	count, err := d.deviceRepo.GetMainDeviceCount(txContext, req.UserId)
	if err != nil {
		d.logger.Error("", "d.deviceRepo.GetMainDeviceCount", req, err)
		tx.Rollback()
		return nil, d.errWarpper.NewDBServiceError(err)
	} else if count > d.MAX_MAIN_DEVICE_COUNT {
		d.logger.Info("", "d.deviceRepo.GetMainDeviceCount", req, nil)
		tx.Rollback()
		return nil, d.errWarpper.NewMainDeviceTooManyError(d.MAX_MAIN_DEVICE_COUNT)
	}

	device, err := d.deviceRepo.BindMainDevice(txContext, req.UserId, req.Platform, req.Version, req.DeviceId)
	err = tx.Commit().Error
	if err != nil {
		d.logger.Error("", "tx.Commit", req, err)
		return nil, d.errWarpper.NewDBCommitServiceError(err)
	}

	d.logger.Info("", "BindMainDevice.end", req, nil)
	return &dto.BindMainDeviceResponse{
		MainDeviceId: device.Id,
		Ok:           true,
	}, nil
}

func (d *deviceServiceImpl) BindSubDevice(ctx context.Context, req *dto.BindSubDeviceRequest) (*dto.BindSubDeviceResponse, *dtoError.ServiceError) {
	txContext, tx := repository.SetTxContext(ctx)
	ok, err := d.deviceRepo.CheckRepeatedDevice(txContext, req.Platform, req.Version, req.DeviceId)
	if err != nil {
		d.logger.Error("", "d.deviceRepo.CheckRepeatedDevice", req, err)
		tx.Rollback()
		return nil, d.errWarpper.NewDBServiceError(err)
	} else if !ok {
		d.logger.Info("", "d.deviceRepo.CheckRepeatedDevice", req, nil)
		tx.Rollback()
		return nil, d.errWarpper.NewRepeatDeviceError()
	}

	count, err := d.deviceRepo.GetSubDeviceCount(txContext, req.UserId, req.MainDeviceId)
	if err != nil {
		d.logger.Error("", "d.deviceRepo.GetSubDeviceCount", req, err)
		tx.Rollback()
		return nil, d.errWarpper.NewDBServiceError(err)
	} else if count >= d.MAX_SUB_DEVICE_COUNT {
		d.logger.Info("", "d.deviceRepo.GetSubDeviceCount", req, nil)
		tx.Rollback()
		return nil, d.errWarpper.NewSubDeviceTooManyError(d.MAX_SUB_DEVICE_COUNT)
	}

	device, err := d.deviceRepo.BindSubDevice(txContext, req.MainDeviceId, req.Platform, req.Version, req.DeviceId)
	err = tx.Commit().Error
	if err != nil {
		d.logger.Error("", "d.deviceRepo.BindSubDevice", req, err)
		return nil, d.errWarpper.NewDBCommitServiceError(err)
	}

	d.logger.Info("", "BindSubDevice.end", req, nil)
	return &dto.BindSubDeviceResponse{
		SubDeviceId: device.Id,
		Ok:          true,
	}, nil
}

func (d *deviceServiceImpl) UnBindMainDevice(ctx context.Context, req *dto.UnbindMainDeviceRequest) (*dto.UnbindMainDeviceResponse, *dtoError.ServiceError) {
	txContext, tx := repository.SetTxContext(ctx)
	binding, err := d.deviceRepo.CheckMainDeviceBinding(txContext, req.UserId, req.MainDeviceId)
	if err != nil {
		d.logger.Error("", "d.deviceRepo.BindSubDevice", req, err)
		tx.Rollback()
		return nil, d.errWarpper.NewDBServiceError(err)
	} else if !binding {
		d.logger.Info("", "d.deviceRepo.BindSubDevice", req, nil)
		tx.Rollback()
		return nil, d.errWarpper.NewMainDeviceNotBindingError()
	}

	ok, err := d.deviceRepo.UnbindMainDevice(ctx, req.UserId, req.MainDeviceId)
	if err != nil {
		d.logger.Error("", "d.deviceRepo.UnbindMainDevice", req, err)
		return nil, d.errWarpper.NewDBServiceError(err)
	}

	d.logger.Info("", "UnBindMainDevice.end", req, nil)
	return &dto.UnbindMainDeviceResponse{Ok: ok}, nil
}

func (d *deviceServiceImpl) UnBindSubDevice(ctx context.Context, req *dto.UnbindSubDeviceRequest) (*dto.UnbindSubDeviceResponse, *dtoError.ServiceError) {
	txContext, tx := repository.SetTxContext(ctx)
	binding, err := d.deviceRepo.CheckMainDeviceBinding(txContext, req.UserId, req.MainDeviceId)
	if err != nil {
		d.logger.Error("", "d.deviceRepo.CheckMainDeviceBinding", req, err)
		tx.Rollback()
		return nil, d.errWarpper.NewDBServiceError(err)
	} else if !binding {
		d.logger.Info("", "d.deviceRepo.CheckMainDeviceBinding", req, nil)
		tx.Rollback()
		return nil, d.errWarpper.NewMainDeviceNotBindingError()
	}

	err = d.deviceRepo.UnbindSubDevice(ctx, req.MainDeviceId, req.SubDeviceId)
	if err != nil {
		d.logger.Error("", "d.deviceRepo.UnbindSubDevice", req, err)
		return nil, d.errWarpper.NewDBServiceError(err)
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, d.errWarpper.NewDBCommitServiceError(err)
	}

	d.logger.Info("", "UnBindSubDevice.end", req, nil)
	return &dto.UnbindSubDeviceResponse{}, nil
}

func (d *deviceServiceImpl) GetDevicesByUserId(ctx context.Context, req *dto.GetDevicesByUserIdRequest) (*dto.GetDevicesByUserIdResponse, *dtoError.ServiceError) {
	mainDevices, err := d.deviceRepo.GetAllDevicesByUserId(ctx, req.UserId)
	if err != nil {
		d.logger.Error("", "d.deviceRepo.GetAllDevicesByUserId", req, err)
		return nil, d.errWarpper.NewDBServiceError(err)
	}

	response := &dto.GetDevicesByUserIdResponse{
		MainDevices: make([]*dto.MainDevice, 0, len(mainDevices)),
	}

	for _, device := range mainDevices {
		subDevices := make([]*dto.SubDevice, 0, len(device.SubDevices))
		for _, subDevice := range device.SubDevices {
			subDevices = append(subDevices, &dto.SubDevice{
				Id:           subDevice.Id,
				MainDeviceId: subDevice.MainDeviceId,
				Platform:     subDevice.Platform,
				Version:      subDevice.Version,
				DeviceId:     subDevice.DeviceId,
			})
		}

		response.MainDevices = append(response.MainDevices, &dto.MainDevice{
			Id:         device.Id,
			UserId:     device.UserId,
			Platform:   device.Platform,
			Version:    device.Version,
			DeviceId:   device.DeviceId,
			SubDevices: subDevices,
		})
	}

	return response, nil
}

func GetDeviceService() DeviceService {
	return device
}
