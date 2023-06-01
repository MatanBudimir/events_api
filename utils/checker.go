package utils

import "github.com/gin-gonic/gin"

var AcceptedTypes = map[string]string{
	"JSON": "application/json",
}

func ContentChecker(ctx *gin.Context, accept string) bool {
	return ctx.Request.Header.Get("Content-Type") == accept
}
