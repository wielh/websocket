package common

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetUUID(c *gin.Context) {
	requestID := uuid.New().String()
	c.Set("RequestID", requestID)
}

func GetUUID(c context.Context) string {
	val := c.Value("RequestID")
	uuid, _ := val.(string)
	return uuid
}
