package controller

import (
	"device-communication/src/config"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var loginFilter gin.HandlerFunc
var customRecoveryFilter gin.HandlerFunc
var readLoginSession gin.HandlerFunc

func commonMiddleware(g *gin.RouterGroup) {
	g.Use(
		customRecoveryFilter,
		readLoginSession,
	)
}

func init() {
	loginFilter = func(c *gin.Context) {
		ok, _, _ := GetSessionValue(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
			c.Abort()
			return
		}
		c.Next()
	}

	customRecoveryFilter = func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var panicMessage string
				switch v := err.(type) {
				case string:
					panicMessage = v
				case error:
					panicMessage = v.Error()
				default:
					panicMessage = fmt.Sprintf("Unknown panic: %v", v)
				}
				fmt.Printf("Panic recovered: %s\n", panicMessage)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
			}
		}()
		c.Next()
	}

	readLoginSession = sessions.Sessions("login", config.GlobalConfig.RedisSession)
}

func SetSessionValue(c *gin.Context, ID uint64, username string) (string, error) {
	session := sessions.Default(c)
	session.Set("id", strconv.FormatUint(ID, 10))
	session.Set("username", username)
	err := session.Save()
	if err != nil {
		return "", err
	}
	return session.ID(), nil
}

func GetSessionValue(c *gin.Context) (bool, uint64, string) {
	session := sessions.Default(c)
	idStr, ok1 := session.Get("id").(string)
	username, ok2 := session.Get("username").(string)

	if !ok1 || !ok2 {
		return false, 0, ""
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return false, 0, ""
	}
	return true, id, username
}

func GetLoginFilter() func(*gin.Context) {
	return loginFilter
}
