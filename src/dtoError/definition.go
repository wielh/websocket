package dtoError

import (
	"github.com/gin-gonic/gin"
)

type ServiceError struct {
	StatusCode     int
	InternalError  error
	ExtrenalReason string
}

func (s *ServiceError) ToJsonResponse() (statusCode int, H *gin.H) {
	statusCode = s.StatusCode
	H = &gin.H{
		"reason": s.ExtrenalReason,
	}
	return
}

type ServiceErrorWarpper interface {
	websocketErrorWarpper
	commonErrorWarpper
	userErrorWarpper
	dbErrorWarpper
	deviceErrorWarpper
}

type websocketErrorWarpper interface {
	NewWebsocketUpgradeFailedError(err error) *ServiceError
	NewRoomCreateFailedError(reason string) *ServiceError
}

type commonErrorWarpper interface {
	NewParseParametersFailedError(err error) *ServiceError
}

type userErrorWarpper interface {
	NewLoginFailedServiceError(err error) *ServiceError
	NewRessetPasswordServiceError() *ServiceError
	NewUserHasRegisterdError(username string) *ServiceError
	NewUsernameExist(username string) *ServiceError
	NewUserNotExist(Id uint64) *ServiceError
	NewPasswordInvaildError(err error) *ServiceError
}

type dbErrorWarpper interface {
	NewDBServiceError(err error) *ServiceError
	NewDBNoAffectedServiceError() *ServiceError
	NewDBCommitServiceError(err error) *ServiceError
}

type deviceErrorWarpper interface {
	NewRepeatDeviceError() *ServiceError
	NewMainDeviceTooManyError(count int64) *ServiceError
	NewSubDeviceTooManyError(count int64) *ServiceError
	NewMainDeviceNotBindingError() *ServiceError
	NewSubDeviceNotBindingError() *ServiceError
}
