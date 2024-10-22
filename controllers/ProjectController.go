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

func uploadDocument(ctx *gin.Context, project *models.Project) {
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	for _, file := range files {
		var Documents = models.NewDocument()

		if file.Size > 1000000 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo muito grande"})
			return
		}

		if file.Header.Get("Content-Type") != "application/pdf" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não é um PDF"})
			return
		}

		if err := services.Upload(file, project.ID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao fazer upload do arquivo"})
		}

		Documents.ProjectID = project.ID
		Documents.Name = file.Filename
		Documents.Link = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s/%s", AWS_BUCKET_NAME, AWS_REGION, ENV, file.Filename)
		database.DB.Create(&Documents)
		project.Documents = append(project.Documents, *Documents)
	}
}

func NewProject(ctx *gin.Context) {
	var user models.User
	var project = models.NewProject()

	uploadDocument(ctx, project)

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
	var documents []models.Document
	var project models.Project
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID nao encontrado"})
		return
	}

	database.DB.Where("project_id = ?", id).Find(&documents)

	for _, doc := range documents {
		for _, file := range files {
			if doc.Name == file.Filename {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo já existe"})
				return
			}
		}
	}

	if err := database.DB.Where("id = ?", id).First(&project).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Projeto não encontrado"})
		return
	}

	if len(files) > 0 {
		uploadDocument(ctx, &project)
	}

	database.DB.Where("user_id = ?", id).Preload("Documents").Find(&project)

	if ctx.PostForm("name") != "" {
		project.Name = ctx.PostForm("name")
	}

	if ctx.PostForm("company") != "" {
		project.Company = ctx.PostForm("company")
	}

	if ctx.PostForm("cnpj") != "" {
		project.CNPJ = ctx.PostForm("cnpj")
	}

	if ctx.PostForm("value") != "" {
		value, err := strconv.ParseFloat(ctx.PostForm("value"), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido"})
			return
		}
		project.Value = value
	}

	if ctx.PostForm("deadline") != "" {
		project.Deadline = ctx.PostForm("deadline")
	}

	if err := database.DB.Save(&project).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao se comunicar com o DB"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Projeto atualizado com sucesso", "project": project})
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
