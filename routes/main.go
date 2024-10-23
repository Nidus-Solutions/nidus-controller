// Arquivo de rotas da aplicação
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/controllers"
	"github.com/jgb27/nidus-controller-projects/middlewares"
)

func Routes(router *gin.Engine) *gin.RouterGroup {
	// Grupo de rotas para o controller de admin, apenas o Admin pode acessar
	routes := router.Group("/admin")
	{
		// Admin controller: Aqui tem apenas rotas relacionado ao Admin em si
		routes.GET("/", middlewares.CheckAuthAdmin, controllers.GetAllAdmin)
		routes.POST("/", controllers.LoginAdmin)
		routes.POST("/newadmin", controllers.NewAdmin)
		routes.PUT("/", middlewares.CheckAuthAdmin, controllers.EditAdmin)
		routes.DELETE("/:id", middlewares.CheckAuthAdmin, controllers.DeleteAdmin)
		routes.DELETE("/user/:id", middlewares.CheckAuthAdmin, controllers.DeleteUser)

		// Admin Project: Aqui tem apenas rotas relacionado ao projeto, mas apenas que o admin pode acessar
		routes.POST("/project", middlewares.CheckAuthAdmin, controllers.NewProject)
		routes.GET("/projects", middlewares.CheckAuthAdmin, controllers.GetAllProjects)
		routes.PUT("/project/:id", middlewares.CheckAuthAdmin, controllers.EditProject)
		routes.DELETE("/project/:id", middlewares.CheckAuthAdmin, controllers.DeleteProject)

		// Admin user: Aqui tem apenas rotas relacionado ao user, mas apenas que o admin pode acessar
		routes.POST("/user", middlewares.CheckAuthAdmin, controllers.NewUser)
		routes.GET("/user", middlewares.CheckAuthAdmin, controllers.GetAllUsers)
	}

	// Grupo de rotas para o controller de user: São poucas pq a ideia é o admin controlar tudo
	routes = router.Group("/user")
	{
		routes.POST("/", controllers.LoginUser)
		routes.PUT("/", middlewares.CheckAuth, controllers.EditUser)
		routes.GET("/project", middlewares.CheckAuth, controllers.GetProjectByUserId)
	}
	return routes
}
