package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
	routes "github.com/jgb27/nidus-controller-projects/routes"
)

func init() {
	database.ConnectToDatabase()
	database.DB.AutoMigrate(models.Admin{}, models.User{}, models.Project{}, models.Document{})
}

func main() {
	app := gin.Default()

	routes.Routes(app)

	app.Run(":3000")
}
