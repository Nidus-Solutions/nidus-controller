package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/controllers"
)

func Routes(router *gin.Engine) *gin.RouterGroup {
	routes := router.Group("/admin")
	{
		// Admin controller
		routes.POST("/", controllers.NewAdmin)
		routes.POST("/login", controllers.LoginAdmin)
		routes.PUT("/", controllers.EditAdmin)
		routes.DELETE("/:id", controllers.DeleteAdmin)
		routes.DELETE("/user/:id", controllers.DeleteUser)

		// Admin Project
		routes.POST("/project", controllers.NewProject)
		routes.GET("/projects", controllers.GetAllProjects)
		routes.PUT("/project/:id", controllers.EditProject)
		routes.DELETE("/project/:id", controllers.DeleteProject)

		// Admin user
		routes.POST("/user", controllers.NewUser)
		routes.GET("/user", controllers.GetAllUsers)
	}

	routes = router.Group("/user")
	{
		routes.POST("/", controllers.LoginUser)
		routes.PUT("/:id", controllers.EditUser) // ID tpm, depois vai pegar pelo token
		routes.GET("/project/:id", controllers.GetProjectByUserId)
	}
	return routes
}
