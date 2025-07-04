package service

import (
	"context"
	"device-communication/src/dto"
	"device-communication/src/dtoError"
	logger "device-communication/src/log"
	"device-communication/src/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	UserRegisterService(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, *dtoError.ServiceError)
	UserLoginService(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, *dtoError.ServiceError)
	ResetPasswordService(ctx context.Context, req *dto.ResetPasswordRequest) *dtoError.ServiceError
}

type userServiceImpl struct {
	userRepo   repository.UserRepository
	errWarpper dtoError.ServiceErrorWarpper
	logger     logger.Logger
}

var user UserService

func init() {
	user = &userServiceImpl{
		userRepo:   repository.GetuserRepository(),
		errWarpper: dtoError.GetServiceErrorWarpper(),
		logger:     logger.NewInfoLogger(),
	}
}

func GetUserService() UserService {
	return user
}

type Password struct {
	hashedByte []byte
}

func newPasswordByRaw(rawPassword string) (Password, error) {
	if len(rawPassword) < 5 || len(rawPassword) > 50 {
		return Password{}, errors.New("password length must be between 5 and 50 characters")
	}

	b, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}
	return Password{hashedByte: b}, nil
}

func newPasswordByHashed(hashedPassword string) Password {
	return Password{hashedByte: []byte(hashedPassword)}
}

func (p *Password) Hashed() string {
	return string(p.hashedByte)
}

func (p *Password) Check(password string) bool {
	err := bcrypt.CompareHashAndPassword(p.hashedByte, []byte(password))
	return err == nil
}

func (u *userServiceImpl) UserRegisterService(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, *dtoError.ServiceError) {
	password, err := newPasswordByRaw(req.Password)
	data := map[string]any{
		"username": req.Username,
		"name":     req.Name,
		"email":    req.Email,
	}

	if err != nil {
		return nil, u.errWarpper.NewPasswordInvaildError(err)
	}

	userModel, ok, err := u.userRepo.UserRegister(ctx, req.Username, password.Hashed(), req.Name, req.Email)
	if err != nil {
		u.logger.Error("", "u.userRepo.UserRegister", data, err)
		return nil, u.errWarpper.NewDBServiceError(err)
	} else if !ok {
		u.logger.Info("", "u.userRepo.UserRegister", data, nil)
		return nil, u.errWarpper.NewUserHasRegisterdError(req.Username)
	}

	u.logger.Info("", "UserRegisterService.end", data, err)
	return &dto.UserRegisterResponse{ID: userModel.Id}, nil
}

func (u *userServiceImpl) UserLoginService(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, *dtoError.ServiceError) {
	userModel, exist, err := u.userRepo.SelectUserByName(ctx, req.Username)
	data := map[string]any{"username": req.Username}

	if err != nil {
		u.logger.Error("", "u.userRepo.SelectUserByName", data, err)
		return nil, u.errWarpper.NewDBServiceError(err)
	} else if !exist {
		u.logger.Info("", "u.userRepo.SelectUserByName", data, nil)
		return nil, u.errWarpper.NewLoginFailedServiceError(nil)
	}

	password := newPasswordByHashed(userModel.Password)
	passwordMatch := password.Check(req.Password)
	if !passwordMatch {
		u.logger.Info("", "password.Check", data, nil)
		return nil, u.errWarpper.NewLoginFailedServiceError(err)
	}

	u.logger.Info("", "UserLoginService.end", data, nil)
	return &dto.UserLoginResponse{
		ID:       userModel.Id,
		Username: userModel.Username,
	}, nil
}

func (u *userServiceImpl) ResetPasswordService(ctx context.Context, req *dto.ResetPasswordRequest) *dtoError.ServiceError {
	txContext, tx := repository.SetTxContext(ctx)
	user, ok, err := u.userRepo.SelectUserByName(txContext, req.Username)
	data := map[string]any{"username": req.Username}

	if err != nil {
		u.logger.Error("", "u.userRepo.SelectUserByName", data, err)
		tx.Rollback()
		return u.errWarpper.NewDBServiceError(err)
	} else if !ok {
		u.logger.Info("", "u.userRepo.SelectUserByName", data, nil)
		tx.Rollback()
		return u.errWarpper.NewRessetPasswordServiceError()
	}

	password := newPasswordByHashed(user.Password)
	passwordMatch := password.Check(req.Password)
	if !passwordMatch {
		u.logger.Info("", "password.Check", data, nil)
		return u.errWarpper.NewLoginFailedServiceError(err)
	}

	newPassword, err := newPasswordByRaw(req.NewPassword)
	if err != nil {
		u.logger.Error("", "newPasswordByRaw", data, err)
		tx.Rollback()
		return u.errWarpper.NewPasswordInvaildError(err)
	}

	ok, err = u.userRepo.UpdatePassword(txContext, user.Id, newPassword.Hashed())
	if err != nil {
		u.logger.Error("", "u.userRepo.UpdatePassword", data, err)
		tx.Rollback()
		return u.errWarpper.NewDBServiceError(err)
	} else if !ok {
		u.logger.Info("", "u.userRepo.UpdatePassword", data, nil)
		tx.Rollback()
		return u.errWarpper.NewDBNoAffectedServiceError()
	}

	err = tx.Commit().Error
	if err != nil {
		u.logger.Error("", "tx.Commit", data, err)
		return u.errWarpper.NewDBCommitServiceError(err)
	}
	return nil
}
