package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
)

func LoginAdmin(ctx *gin.Context) {
	var admin models.Admin
	ctx.BindJSON(&admin)

	if admin.Username == "" || admin.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Preencha todos os campos"})
		return
	}

	var checkAdmin models.Admin
	database.DB.Where("username = ? AND password = ?", admin.Username, admin.Password).First(&checkAdmin)

	if checkAdmin.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario ou senha invalidos"})
		return
	}

	checkAdmin.Password = ""

	ctx.JSON(http.StatusOK, checkAdmin)
}

func NewAdmin(ctx *gin.Context) {
	admin := models.NewAdmin()
	ctx.BindJSON(&admin)

	if admin.Name == "" || admin.Username == "" || admin.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Preencha todos os campos"})
		return
	}

	var checkAdmin models.Admin
	database.DB.Where("username = ?", admin.Username).First(&checkAdmin)
	if checkAdmin.Username != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario ja existe"})
		return
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao criar usuario"})
		return
	}

	admin.Password = ""

	ctx.JSON(http.StatusOK, admin)
}

func EditAdmin(ctx *gin.Context) {
	var NewAdmin models.Admin
	ctx.BindJSON(&NewAdmin)

	if NewAdmin.Name == "" || NewAdmin.Username == "" || NewAdmin.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Preencha todos os campos"})
		return
	}

	var admin models.Admin
	database.DB.Where("id = ?", NewAdmin.ID).First(&admin)

	if admin.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario nao encontrado"})
		return
	}

	if err := database.DB.Model(&admin).Updates(NewAdmin).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao atualizar usuario"})
		return
	}

	admin.Password = ""

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario atualizado com sucesso"})
}

func DeleteAdmin(ctx *gin.Context) {
	id := ctx.Param("id")

	fmt.Printf("ID: %s\n", id)

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID nao encontrado"})
		return
	}

	var admin models.Admin

	database.DB.Where("id = ?", id).First(&admin)

	if admin.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario nao encontrado"})
		return
	}

	if err := database.DB.Delete(&admin).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao deletar usuario"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario deletado com sucesso"})
}
