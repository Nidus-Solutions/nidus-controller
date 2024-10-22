package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
	"github.com/jgb27/nidus-controller-projects/services"
)

var (
	AWS_REGION      = services.LoadEnv("AWS_REGION")
	AWS_BUCKET_NAME = services.LoadEnv("AWS_BUCKET_NAME")
	ENV             = services.LoadEnv("ENV")
)

func NewProject(ctx *gin.Context) {
	var user models.User
	var project = models.NewProject()
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	for _, file := range files {
		var Documents = models.NewDocument()

		if err := services.Upload(file); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao fazer upload do arquivo"})
		}

		Documents.ProjectID = project.ID
		Documents.Name = file.Filename
		Documents.Link = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/%s", AWS_BUCKET_NAME, AWS_REGION, ENV, file.Filename)
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

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID nao encontrado"})
		return
	}

	var project models.Project

	if err := database.DB.Where("id = ?", id).First(&project).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Projeto não encontrado"})
		return
	}

	if ctx.PostForm("file") != "" {
		var document []models.Document

		database.DB.Where("project_id = ?", id).Find(&document)

		log.Println(document)
	}

	value, err := strconv.ParseFloat(ctx.PostForm("value"), 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido"})
		return
	}

	database.DB.Model(&project).Updates(models.Project{
		Name:     ctx.PostForm("name"),
		Company:  ctx.PostForm("company"),
		CNPJ:     ctx.PostForm("cnpj"),
		Deadline: ctx.PostForm("deadline"),
		Value:    value,
	})

	ctx.JSON(http.StatusOK, project)
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
