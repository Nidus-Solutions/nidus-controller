package main

import (
	"github.com/gin-gonic/gin"
	db "github.com/jgb27/nidus-controller-projects/database"
	routes "github.com/jgb27/nidus-controller-projects/routes"
)

func main() {
	app := gin.Default()

	db.ConnectToDatabase()

	routes.Routes(app)

	app.Run(":8080")
}
