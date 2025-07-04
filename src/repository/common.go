package repository

import (
	"context"
	"device-communication/src/config"

	"gorm.io/gorm"
)

type transectioKey string

const tk transectioKey = "gorm-transection"

func SetTxContext(ctx context.Context) (context.Context, *gorm.DB) {
	tx := config.GlobalConfig.NewTransection()
	return context.WithValue(ctx, tk, tx), tx
}

func GetTxContext(ctx context.Context, defaultTx *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(tk).(*gorm.DB)
	if !ok {
		return defaultTx
	}
	return tx.Session(&gorm.Session{NewDB: true})
	// .Session(&gorm.Session{NewDB: true}) 避免查詢條件汙染
}
