package api

import (
	jwt2 "b2match_api/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "authorization header was not provided",
			})
			return
		}

		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "authorization header was not provided",
			})
			return
		}

		token, err := jwt2.Parse(headerParts[1])

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token was provided",
			})
			return
		}

		ctx.Set("user", token.Claims.(*jwt2.Claims).User)

		ctx.Next()
	}
}
