package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/jgb27/nidus-controller-projects/controllers"
)

func Routes(router *gin.Engine) *gin.RouterGroup {
	routes := router.Group("/users")
	{
		routes.POST("/", controllers.CreateUser)
		routes.GET("/", controllers.FindAllUser)
		routes.GET("/:id", controllers.FindUserById)
		routes.DELETE("/:id", controllers.DeleteUser)
	}

	return routes
}
