// Criando um token JWT para o usuário e para o admin,
package controllers

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jgb27/nidus-controller-projects/models"
	"github.com/jgb27/nidus-controller-projects/services"
)

// gerando um token JWT para o usuário
func GenerateTokenUser(user *models.User) string {
	secret := services.LoadEnv("JWT_SECRETE_KEY")

	// Criando um token JWT com o ID do usuário e a data de expiração em 24 horas
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

// Gere um token JWT para o admin
func GenerateTokenAdmin(admin *models.Admin) string {
	secret := services.LoadEnv("JWT_SECRETE_KEY")

	// Criando um token JWT com o ID do admin, se ele é um admin e a data de expiração em 1 hora
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      admin.ID,
		"isAdmin": admin.IsAdmin,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}
