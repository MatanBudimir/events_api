package jwt

import (
	"b2match_api/pkg/database/models"
	"b2match_api/utils"
	"errors"
	"fmt"
	jwtPackage "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	User models.User `json:"user"`
	jwtPackage.RegisteredClaims
}

func Parse(token string) (*jwtPackage.Token, error) {
	return jwtPackage.ParseWithClaims(token, &Claims{}, func(token *jwtPackage.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtPackage.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := utils.GetEnv("JWT_SECRET", "")

		if len(secret) == 0 {
			return nil, errors.New("no secret provided")
		}

		return []byte(secret), nil
	})
}
