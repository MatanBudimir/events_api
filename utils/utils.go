package utils

import (
	"b2match_api/pkg/database/models"
	"github.com/gin-gonic/gin"
	"os"
)

func GetEnv(key, defaultValue string) string {
	if val := os.Getenv(key); len(val) > 0 {
		return val
	}

	return defaultValue
}

func GetUser(ctx *gin.Context) models.User {
	return ctx.MustGet("user").(models.User)
}
