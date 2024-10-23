// Arquivo responsável por controlar as requisições relacionadas a projetos

package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/database"
	"github.com/jgb27/nidus-controller-projects/models"
	"github.com/jgb27/nidus-controller-projects/services"
)

// Variáveis de ambiente fixas para o bucket da AWS
var (
	AWS_REGION      = services.LoadEnv("AWS_REGION")
	AWS_BUCKET_NAME = services.LoadEnv("AWS_BUCKET_NAME")
	ENV             = services.LoadEnv("ENV")
)

// Função para fazer upload de documentos para o bucket da AWS, tanto na criação quanto na edição de um projeto
func uploadDocument(ctx *gin.Context, project *models.Project) error {
	url := "https://%s.s3.%s.amazonaws.com/%s/%s/%s" // Modelo de URL para acessar o arquivo na AWS
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	// os dados chegam em lista, então é necessário fazer um for para acessar um a um
	for _, file := range files {
		var Documents = models.NewDocument() // crie um novo objeto Document com base no modelo

		// Algumas verificações para garantir que o arquivo é um PDF e não é muito grande
		if file.Size > 5000000 {
			return errors.New("arquivo muito grande")
		}

		if file.Header.Get("Content-Type") != "application/pdf" {
			return errors.New("somente PDF")
		}

		// Se tudo estiver ok, faça o upload do arquivo para a AWS
		if err := services.UploadToAws(file, project.ID); err != nil {
			return errors.New("erro ao fazer upload")
		}

		// Salvando no bando de dados apenas o link do arquivo
		Documents.ProjectID = project.ID
		Documents.Name = file.Filename
		Documents.Link = fmt.Sprintf(url, AWS_BUCKET_NAME, AWS_REGION, ENV, project.ID, file.Filename) // aqui é onde é montado o link do arquivo
		database.DB.Create(&Documents)
		project.Documents = append(project.Documents, *Documents)
	}

	return nil
}

// Criando um novo projeto com base nos dados enviados pelo usuário
func NewProject(ctx *gin.Context) {
	var user models.User
	var project = models.NewProject()

	// Fazendo o upload dos documentos
	if err := uploadDocument(ctx, project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificando se o usuário existe para poder vincular o projeto a ele
	if err := database.DB.Where("id = ?", ctx.PostForm("userId")).First(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Usuário não encontrado"})
		return
	}

	project.UserID = user.ID
	project.Name = ctx.PostForm("name")
	project.Company = ctx.PostForm("company")
	project.CNPJ = ctx.PostForm("cnpj")

	// O value é um float, então é necessário fazer uma conversão
	value, err := strconv.ParseFloat(ctx.PostForm("value"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido"})
		return
	}

	project.Deadline = ctx.PostForm("deadline")
	project.Value = value

	// Salvando o projeto completo no banco de dados
	if err := database.DB.Create(&project).Error; err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao se comunicar com o DB"})
		return
	}

	ctx.JSON(http.StatusCreated, project)
}

// Editando um projeto com base nos dados enviados pelo usuário
func EditProject(ctx *gin.Context) {
	id := ctx.Param("id")
	var project models.Project // Variavel do tipo Project para armazenar o projeto que será editado

	// Criando uma lista com base no models do Document, já que no banco é retornado uma lista
	var documents []models.Document
	form, _ := ctx.MultipartForm()
	files := form.File["files"]

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID nao encontrado"})
		return
	}

	// Salvando todos os documentos referente ao projeto
	database.DB.Where("project_id = ?", id).Find(&documents)

	// Verificando se o arquivo enviado já existe no projeto
	for _, doc := range documents {
		for _, file := range files {
			if doc.Name == file.Filename {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo já existe"})
				return
			}
		}
	}

	// Alterando apenas o projeto que foi solicitado e apenas as partes que foram enviadas
	if err := database.DB.Where("id = ?", id).First(&project).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Projeto não encontrado"})
		return
	}

	// Se houver arquivos, faça o upload
	if len(files) > 0 {
		if err := uploadDocument(ctx, &project); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	database.DB.Where("user_id = ?", id).Preload("Documents").Find(&project) // Buscando o projeto com os documentos referente

	// Como vem em forme, tem que ser feito todos esses 'ifs' para verificar se o campo foi enviado
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

	// Salvando o projeto no banco de dados
	if err := database.DB.Save(&project).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao se comunicar com o DB"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Projeto atualizado com sucesso", "project": project})
}

// Deletando um projeto com base no ID
func DeleteProject(ctx *gin.Context) {
	id := ctx.Param("id")
	var documents []models.Document

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID nao encontrado"})
		return
	}

	database.DB.Where("project_id = ?", id).Find(&documents) // Buscando todos os documentos referente ao projeto

	// Deletando todos os documentos referente ao projeto
	for _, doc := range documents {
		if err := database.DB.Delete(&doc).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao deletar arquivo"})
			return
		}
	}

	// Deletando o projeto em si, tem que ser feito nessa ordem, se não da erro!!!!
	if err := database.DB.Where("id = ?", id).Delete(&models.Project{}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao se comunicar com o DB"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Projeto deletado com sucesso"})
}

// Buscando todos os projetos cadastrados
func GetAllProjects(ctx *gin.Context) {
	var project []models.Project

	database.DB.Model(&models.Project{}).Preload("Documents").Find(&project)

	ctx.JSON(http.StatusOK, project)
}

// Buscando todos os projetos de um usuário específico
func GetProjectByUserId(ctx *gin.Context) {
	currentUser, _ := ctx.Get("currentUser")
	id := currentUser.(models.User).ID

	var project []models.Project

	database.DB.Where("user_id = ?", id).Preload("Documents").Find(&project)

	ctx.JSON(http.StatusOK, project)
}
