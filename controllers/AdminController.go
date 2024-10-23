// Arquivo responsável por controlar as rotas de Admins
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
	"golang.org/x/crypto/bcrypt"
)

// Função responsável por retornar todos os Admins que temos no banco de dados
func GetAllAdmin(ctx *gin.Context) {
	var admins []models.Admin
	database.DB.Find(&admins)

	// Remove a senha do retorno de todos os admins
	for i := range admins {
		admins[i].Password = ""
	}

	ctx.JSON(http.StatusOK, admins)
}

// Função responsável por gerar o token de autenticação do Admin
func LoginAdmin(ctx *gin.Context) {
	var admin models.Admin
	ctx.BindJSON(&admin)

	if admin.Username == "" || admin.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Preencha todos os campos"})
		return
	}

	var checkAdmin models.Admin
	database.DB.Where("username = ?", admin.Username).First(&checkAdmin)

	// Verifica se o usuario existe e se a senha está correta
	if bcrypt.CompareHashAndPassword([]byte(checkAdmin.Password), []byte(admin.Password)) != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Usuario ou senha invalidos"})
		return
	}

	if checkAdmin.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario ou senha invalidos"})
		return
	}

	//Retornando apenas o necessario para o token
	ctx.JSON(http.StatusOK, gin.H{
		"token":    GenerateTokenAdmin(&checkAdmin),
		"name":     checkAdmin.Name,
		"username": checkAdmin.Username,
	})
}

// Função responsável por criar um novo Admin
func NewAdmin(ctx *gin.Context) {
	admin := models.NewAdmin()
	ctx.BindJSON(&admin)

	cost := bcrypt.DefaultCost

	// Verifica se todos os campos foram preenchidos
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

	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), cost) // Gera o hash da senha

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao gerar hash"})
		return
	}

	admin.Password = string(hash) // Salvando a senha criptografada

	if err := database.DB.Create(&admin).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao criar usuario"})
		return
	}

	admin.Password = "" // Remove a senha do retorno

	ctx.JSON(http.StatusOK, admin)
}

// Função responsável por editar um Admin
func EditAdmin(ctx *gin.Context) {
	// Pega o ID do Admin logado
	currentAdmin, _ := ctx.Get("currentAdmin")
	id := currentAdmin.(models.Admin).ID

	// Capturando os dados que vão ser alterados do Admin
	var NewAdmin models.Admin
	ctx.BindJSON(&NewAdmin)

	// Campturando os dados, do admin que vai ser alterado, no banco de dados
	var admin models.Admin
	database.DB.Where("id = ?", id).First(&admin)

	// Verifica se o admin existe
	if admin.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario nao encontrado"})
		return
	}

	// Fazendo a alteração apenas dos dados que foram passados
	if err := database.DB.Model(&admin).Updates(NewAdmin).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao atualizar usuario"})
		return
	}

	admin.Password = ""

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario atualizado com sucesso"})
}

// Função responsável por deletar um Admin, nada demais, só remove mesmo
func DeleteAdmin(ctx *gin.Context) {
	id := ctx.Param("id")

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
