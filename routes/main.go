package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/controllers"
)

func Routes(router *gin.Engine) *gin.RouterGroup {
	routes := router.Group("/admin")
	{
		routes.POST("/", controllers.NewAdmin)
		routes.POST("/login", controllers.LoginAdmin)
		routes.PUT("/", controllers.EditAdmin)
		routes.DELETE("/:id", controllers.DeleteAdmin)
	}

	routes = router.Group("/user")
	{
		routes.POST("/", controllers.LoginUser)
	}

	return routes
}
