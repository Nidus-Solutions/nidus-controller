package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
	routes "github.com/jgb27/nidus-controller-projects/routes"
)

func main() {
	app := gin.Default()

	database.ConnectToDatabase()

	database.DB.AutoMigrate(models.Admin{})
	database.DB.AutoMigrate(models.User{})
	database.DB.AutoMigrate(models.Project{})

	routes.Routes(app)

	app.Run(":3000")
}
