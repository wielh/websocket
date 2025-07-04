package repository

import (
	"context"
	"device-communication/src/config"
	"device-communication/src/model"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	UserRegister(ctx context.Context, username string, password string, name string, email string) (*model.User, bool, error)
	SelectUserByName(ctx context.Context, username string) (*model.User, bool, error)
	UpdatePassword(ctx context.Context, ID uint64, newHashedPassword string) (ok bool, err error)
	CheckUserExist(ctx context.Context, ID uint64) (exist bool, err error)
}

type userRepositoryImpl struct {
	DB *gorm.DB
}

var user UserRepository

func init() {
	user = &userRepositoryImpl{
		DB: config.GlobalConfig.DB,
	}
}

func GetuserRepository() UserRepository {
	return user
}

func (a *userRepositoryImpl) UserRegister(ctx context.Context, username string, password string, name string, email string) (*model.User, bool, error) {
	tx := GetTxContext(ctx, a.DB)

	user := model.User{
		Username: username,
		Password: password,
		Name:     name,
		Email:    email,
	}

	result := tx.Where("username=?", username).FirstOrCreate(&user)
	if result.Error != nil {
		return nil, false, result.Error
	}
	return &user, result.RowsAffected > 0, nil
}

func (a *userRepositoryImpl) SelectUserByName(ctx context.Context, username string) (*model.User, bool, error) {
	tx := GetTxContext(ctx, a.DB)
	var user = model.User{Username: username}
	result := tx.Select("id", "username", "password").Where("username=?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, result.Error
	}
	return &user, true, nil
}

func (a *userRepositoryImpl) UpdatePassword(ctx context.Context, ID uint64, newHashedPassword string) (bool, error) {
	tx := GetTxContext(ctx, a.DB)
	result := tx.Model(&model.User{}).Where("id=?", ID).Updates(map[string]interface{}{"password": newHashedPassword})
	if result.Error != nil {
		return false, result.Error
	} else if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (a *userRepositoryImpl) CheckUserExist(ctx context.Context, ID uint64) (exist bool, err error) {
	tx := GetTxContext(ctx, a.DB)
	var user model.User
	result := tx.Select("username").Where("id=?", ID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
