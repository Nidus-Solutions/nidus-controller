package main

import (
	"github.com/gin-gonic/gin"
	routes "github.com/jgb27/nidus-controller-projects/routes"
)

func main() {
	app := gin.Default()

	routes.Routes(app)

	app.Run(":8080")
}
