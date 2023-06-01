package main

import (
	"b2match_api/pkg/database"
	"b2match_api/pkg/database/models"
	jwt2 "b2match_api/pkg/jwt"
	"b2match_api/utils"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("please provide action")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	switch os.Args[1] {
	case "jwt":
		token, err := generateJWT()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Token -", token)
	default:
		fmt.Println("Invalid action provided")
	}
}

func generateJWT() (string, error) {
	if err := database.Connect(); err != nil {
		log.Fatalln(err)
	}

	// get the first user
	row := database.DB.QueryRow("SELECT * FROM users LIMIT 1")

	if row.Err() != nil {
		log.Fatalln(row.Err())
	}

	user := models.User{}

	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.OrganizationID); err != nil {
		log.Fatalln(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt2.Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    utils.GetEnv("JWT_ISSUER", "b2match"),
		},
	})

	secret := utils.GetEnv("JWT_SECRET", "")

	if len(secret) == 0 {
		return "", errors.New("no secret was provided")
	}

	return token.SignedString([]byte(secret))
}
