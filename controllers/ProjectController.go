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

func EditProject(ctx *gin.Context) {
	id := ctx.Param("id")
	var project models.Project
	ctx.BindJSON(&project)

	var projectDB models.Project
	database.DB.Where("id = ?", id).First(&projectDB)

	if projectDB.ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Projeto nao encontrado"})
		return
	}

	if err := database.DB.Model(&projectDB).Updates(&project).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao atualizar projeto"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Projeto atualizado com sucesso"})

}

func DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID nao encontrado"})
		return
	}

	if err := database.DB.Where("id = ?", id).Delete(&models.Project{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao se comunicar com o DB"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Projeto deletado com sucesso"})
}

func GetAllProjects(ctx *gin.Context) {
	var project []models.Project

	database.DB.Find(&project)

	ctx.JSON(http.StatusOK, project)
}

func GetProjectByUserId(ctx *gin.Context) {
	id := ctx.Param("id")

	var project []models.Project

	database.DB.Where("user_id = ?", id).Find(&project)

	ctx.JSON(http.StatusOK, project)
}
