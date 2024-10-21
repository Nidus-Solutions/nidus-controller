package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(ctx *gin.Context) {
	var users []models.User

	database.DB.Find(&users)

	ctx.JSON(http.StatusOK, users)
}

func LoginUser(ctx *gin.Context) {
	var user models.User
	ctx.BindJSON(&user)

	if user.UserName == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Preencha todos os campos"})
		return
	}

	var checkUser models.User
	database.DB.Where("user_name = ?", user.UserName).First(&checkUser)

	if bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(user.Password)) != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Usuario ou senha invalidos"})
		return
	}

	if checkUser.UserName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario ou senha invalidos"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"name":     checkUser.Name,
		"username": checkUser.UserName,
		"id":       checkUser.ID,
	})
}

func NewUser(ctx *gin.Context) {
	newUser := models.NewUser()
	ctx.BindJSON(&newUser)

	cost := bcrypt.DefaultCost

	var checkUser models.User

	database.DB.Where("username = ?", newUser.UserName).First(&checkUser)
	if checkUser.UserName != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario já existe"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), cost)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newUser.Password = string(hash)

	if err := database.DB.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao criar usuario"})
		return
	}

	ctx.JSON(http.StatusOK, newUser)
}

func EditUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	ctx.BindJSON(&user)

	var checkUser models.User
	database.DB.Where("id = ?", id).First(&checkUser)

	if checkUser.UserName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario nao encontrado"})
		return
	}

	if err := database.DB.Model(&checkUser).Updates(user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao atualizar usuario"})
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario atualizado com sucesso"})
}

func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID nao encontrado"})
		return
	}

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	if err := database.DB.Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao deletar usuario"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario deletado com sucesso"})
}
