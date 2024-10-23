// Arquivo responsável por controlar as rotas de usuários
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
	"golang.org/x/crypto/bcrypt"
)

// Função responsável por retornar todos os usuários cadastrados no banco de dados
func GetAllUsers(ctx *gin.Context) {
	var users []models.User

	database.DB.Find(&users)

	ctx.JSON(http.StatusOK, users)
}

// Função responsável pelo login do usuário
func LoginUser(ctx *gin.Context) {
	var user models.User // Variável para armazenar o usuário que vamos logar
	ctx.BindJSON(&user)

	// Verifica simples para ver se os campos estão preenchidos
	if user.UserName == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Preencha todos os campos"})
		return
	}

	var checkUser models.User // Variável para verificar se o usuário existe
	database.DB.Where("user_name = ?", user.UserName).First(&checkUser)

	// verifica se o usuário existe
	if bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(user.Password)) != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Usuario ou senha invalidos"})
		return
	}

	if checkUser.UserName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario ou senha invalidos"})
		return
	}

	// retornando apenas o necessario
	ctx.JSON(http.StatusOK, gin.H{
		"token":    GenerateTokenUser(&checkUser),
		"name":     checkUser.Name,
		"username": checkUser.UserName,
	})
}

// criando um novo usuário
func NewUser(ctx *gin.Context) {
	cost := bcrypt.DefaultCost  // Uma constate para poder gerar o hash da senha
	newUser := models.NewUser() // obj que vai receber os dados do novo usuário
	ctx.BindJSON(&newUser)      // pegando os dados enviado pelo usuário, em formato json

	// Verficação para saber se o usuário já existe
	var checkUser models.User
	database.DB.Where("username = ?", newUser.UserName).First(&checkUser)
	if checkUser.UserName != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuario já existe"})
		return
	}

	// Gerando o hash da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), cost)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newUser.Password = string(hash) // Salvanado a senha já com o hash

	if err := database.DB.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao criar usuario"})
		return
	}

	// Obj que vai receber apenas os dados necessarios para gerar o token
	// não sei pq, mas assim funciona
	resultUser := models.User{}

	resultUser.ID = newUser.ID
	resultUser.Name = newUser.Name
	resultUser.UserName = newUser.UserName

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Usuario criado com sucesso",
		"user":    resultUser,
		"token":   GenerateTokenUser(&resultUser),
	})
}

// Função responsável por retornar um usuário específico
func EditUser(ctx *gin.Context) {
	// Salvando apenas o id do usuário
	currentUser, _ := ctx.Get("currentUser")
	id := currentUser.(models.User).ID

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID nao encontrado"})
		return
	}

	// Variável para armazenar o usuário que vamos atualizar
	var user models.User
	ctx.BindJSON(&user)

	// Verificação para saber se o usuário existe
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

	user.Password = "" // Não retornar a senha

	ctx.JSON(http.StatusOK, gin.H{"message": "Usuario atualizado com sucesso"})
}

// remove o usuário
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
