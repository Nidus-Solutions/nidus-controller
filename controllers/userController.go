package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	database "github.com/jgb27/nidus-controller-projects/database"
	entities "github.com/jgb27/nidus-controller-projects/entities"
	user "github.com/jgb27/nidus-controller-projects/entities"
)

var db *sql.DB = database.ConnectToDatabase()

func handleDBError(ctx *gin.Context, err error, message string) {
	fmt.Println("Error:", err)
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": message})
}

func CreateUser(ctx *gin.Context) {
	user := user.NewUser()

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO clients (id, name, phone, mail, cpf, createdAt) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Name, user.Phone, user.Mail, user.CPF, user.CreatedAt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func FindAllUser(ctx *gin.Context) {
	rows, err := db.Query("SELECT * FROM clients")

	if err != nil {
		fmt.Println("Erro ao executar a query:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var users []entities.User

	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.ID, &user.Name, &user.Phone, &user.Mail, &user.CPF, &user.CreatedAt)

		if err != nil {
			fmt.Print("Error ao escanerar os resultados:", err)
			continue
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erro geral ao iterar sobre as linhas:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func FindUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	stmt, err := db.Prepare("SELECT * FROM clients WHERE id = $1")

	if err != nil {
		handleDBError(ctx, err, "failed to prepare statement")
		return
	}

	defer stmt.Close()

	row := stmt.QueryRow(id)
	if row.Err() != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	var user entities.User
	err = row.Scan(&user.ID, &user.Name, &user.Phone, &user.Mail, &user.CPF, &user.CreatedAt)
	if err != nil {
		handleDBError(ctx, err, "failed to scan user")
		return
	}

	ctx.JSON(http.StatusOK, user)
}
func DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.Exec("DELETE FROM clients where id = $1", id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
