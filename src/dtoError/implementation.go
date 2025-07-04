package dtoError

import (
	"fmt"
	"net/http"
)

var s ServiceErrorWarpper = &ServiceErrorWarpperImpl{}

type ServiceErrorWarpperImpl struct{}

func (s *ServiceErrorWarpperImpl) NewRoomCreateFailedError(reason string) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusBadRequest,
		InternalError:  nil,
		ExtrenalReason: reason,
	}
}

func (s *ServiceErrorWarpperImpl) NewWebsocketUpgradeFailedError(err error) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusInternalServerError,
		InternalError:  nil,
		ExtrenalReason: "websocket upgrade failed",
	}
}

func (s *ServiceErrorWarpperImpl) NewMainDeviceNotBindingError() *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusNotFound,
		InternalError:  nil,
		ExtrenalReason: "device not binding",
	}
}

func (s *ServiceErrorWarpperImpl) NewSubDeviceNotBindingError() *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusNotFound,
		InternalError:  nil,
		ExtrenalReason: "device not binding",
	}
}

func (s *ServiceErrorWarpperImpl) NewRepeatDeviceError() *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusConflict,
		InternalError:  nil,
		ExtrenalReason: "device has already binded",
	}
}

func (s *ServiceErrorWarpperImpl) NewMainDeviceTooManyError(count int64) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusConflict,
		InternalError:  nil,
		ExtrenalReason: fmt.Sprintf("Allow %d max devices per user", count),
	}
}

func (s *ServiceErrorWarpperImpl) NewSubDeviceTooManyError(count int64) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusConflict,
		InternalError:  nil,
		ExtrenalReason: fmt.Sprintf("Allow %d max sub devices per main device", count),
	}
}

func (s *ServiceErrorWarpperImpl) NewParseParametersFailedError(err error) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusBadRequest,
		InternalError:  err,
		ExtrenalReason: err.Error(),
	}
}

func (s *ServiceErrorWarpperImpl) NewPasswordInvaildError(err error) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusBadRequest,
		InternalError:  err,
		ExtrenalReason: err.Error(),
	}
}

func (s *ServiceErrorWarpperImpl) NewDBCommitServiceError(err error) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusInternalServerError,
		InternalError:  err,
		ExtrenalReason: "Service Temporary Unavailable",
	}
}

func (s *ServiceErrorWarpperImpl) NewDBNoAffectedServiceError() *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusOK,
		InternalError:  nil,
		ExtrenalReason: "No Affected Data",
	}
}

func (s *ServiceErrorWarpperImpl) NewDBServiceError(err error) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusInternalServerError,
		InternalError:  err,
		ExtrenalReason: "Service Temporary Unavailable",
	}
}

func (s *ServiceErrorWarpperImpl) NewLoginFailedServiceError(err error) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusUnauthorized,
		InternalError:  err,
		ExtrenalReason: "LoginFailed",
	}
}

func (s *ServiceErrorWarpperImpl) NewRessetPasswordServiceError() *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusUnauthorized,
		InternalError:  nil,
		ExtrenalReason: "ResetPasswordFailed",
	}
}

func (s *ServiceErrorWarpperImpl) NewUserHasRegisterdError(username string) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusConflict,
		InternalError:  nil,
		ExtrenalReason: fmt.Sprintf("user %s has already registered", username),
	}
}

func (s *ServiceErrorWarpperImpl) NewUserNotExist(Id uint64) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusUnauthorized,
		InternalError:  nil,
		ExtrenalReason: fmt.Sprintf("user %d does not exist", Id),
	}
}

func (s *ServiceErrorWarpperImpl) NewUsernameExist(username string) *ServiceError {
	return &ServiceError{
		StatusCode:     http.StatusConflict,
		InternalError:  nil,
		ExtrenalReason: fmt.Sprintf("username %s is used", username),
	}
}

func GetServiceErrorWarpper() ServiceErrorWarpper {
	return s
}
