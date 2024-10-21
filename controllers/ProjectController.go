package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
)

func NewProject(ctx *gin.Context) {
	var user models.User
	var project = models.NewProject()

	form, _ := ctx.MultipartForm()

	files := form.File["files"]

	for _, file := range files {
		var Documents = models.NewDocument()

		if err := ctx.SaveUploadedFile(file, "uploads/"+file.Filename); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao salvar arquivo"})
			return
		}
		Documents.ProjectID = project.ID
		Documents.Name = file.Filename
		Documents.Link = "uploads/" + file.Filename
		database.DB.Create(&Documents)
		project.Documents = append(project.Documents, *Documents)
	}

	if err := database.DB.Where("id = ?", ctx.PostForm("userId")).First(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuário não encontrado"})
		return
	}

	project.UserID = user.ID
	project.Name = ctx.PostForm("name")
	project.Company = ctx.PostForm("company")
	project.CNPJ = ctx.PostForm("cnpj")

	value, err := strconv.ParseFloat(ctx.PostForm("value"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido"})
		return
	}

	project.Deadline = ctx.PostForm("deadline")
	project.Value = value

	if err := database.DB.Create(&project).Error; err != nil {
		log.Println(err)
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

	database.DB.Model(&models.Project{}).Preload("Documents").Find(&project)

	ctx.JSON(http.StatusOK, project)
}

func GetProjectByUserId(ctx *gin.Context) {
	currentUser, _ := ctx.Get("currentUser")
	id := currentUser.(models.User).ID

	var project []models.Project

	database.DB.Where("user_id = ?", id).Preload("Documents").Find(&project)

	ctx.JSON(http.StatusOK, project)
}
