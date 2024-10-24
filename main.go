package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
	routes "github.com/jgb27/nidus-controller-projects/routes"
	"github.com/jgb27/nidus-controller-projects/services"
)

// Inicindo tudo necessari antes da aplicação rodar
func init() {
	database.ConnectToDatabase()
	database.DB.AutoMigrate(models.Admin{}, models.User{}, models.Project{}, models.Document{})
}

func main() {
	port := services.LoadEnv("PORT") // Carregando a porta do arquivo .env
	app := gin.Default()

	routes.Routes(app) // Iniciando as rotas

	log.Println("Server is running in mode: " + services.LoadEnv("ENV"))
	log.Println("Server is running on port: " + port)
	app.Run(":" + port)
}
