package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
)

func NewProject(ctx *gin.Context) {
	var project = models.NewProject()
	ctx.BindJSON(&project)

	if err := database.DB.Create(&project).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao se comunicar com o DB"})
		return
	}

	ctx.JSON(http.StatusCreated, project)
}

func EditProject() {
}

func DeleteProject() {
}

func GetAllProjects(ctx *gin.Context) {
	var project []models.Project

	database.DB.Raw("SELECT * FROM projects").Scan(&project)

	ctx.JSON(http.StatusOK, project)
}

func GetProjectById(ctx *gin.Context) {
	id := ctx.Param("id")
	var project []models.Project

	database.DB.Raw(`
		SELECT projects.id AS project_id, projects.name AS project_name, users.name AS user_name
		FROM projects JOIN users ON projects.user_id = users.id
		WHERE users.id = ? ORDER BY projects.name;`,
		id).Scan(&project)

	ctx.JSON(http.StatusOK, project)
}
