// Esse arquivo contém os middlewares de autenticação de usuário e admin
package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
)

// CheckAuth verifica se o token de autenticação é válido, para usuários normais
func CheckAuth(c *gin.Context) {

	authHeader := c.GetHeader("Authorization") // vai pegar o token do header

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// verifica se o token está no formato correto, todos seguem o formato "Bearer token"
	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authToken[1] // pega o token, sem o "Bearer"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRETE_KEY")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// claims é um mapa de strings para interfaces, que é o tipo de token.Claims, usa ele para pegar o id do usuário
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// verifica se o token expirou ou não
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// pega o usuário do banco de dados e salva no contexto
	var user models.User
	database.DB.Where("ID=?", claims["id"]).Find(&user)

	if user.ID == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentUser", user)

	c.Next() // chama o próximo middleware
}

// CheckAuthAdmin verifica se o token de autenticação é válido, para administradores
func CheckAuthAdmin(c *gin.Context) {

	authHeader := c.GetHeader("Authorization") // vai pegar o token do header

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// verifica se o token está no formato correto, todos seguem o formato "Bearer token
	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authToken[1] // pega o token, sem o "Bearer"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRETE_KEY")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// pega o admin do banco de dados e salva no contexto
	var admin models.Admin
	database.DB.Where("ID=?", claims["id"]).Find(&admin)

	if admin.ID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !admin.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not an admin"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentAdmin", admin)

	c.Next() // chama o próximo middleware
}
