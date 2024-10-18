package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jgb27/nidus-controller-projects/models"
)

func LoginUser(ctx *gin.Context) {
	var user models.User
	ctx.BindJSON(&user)
	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}

func NewUser(ctx *gin.Context) {
	var user models.User
	ctx.BindJSON(&user)
	ctx.JSON(http.StatusCreated, user)
}

func EditUser(ctx *gin.Context) {
	var user models.User
	ctx.BindJSON(&user)
	ctx.JSON(http.StatusOK, user)
}

func DeleteUser(ctx *gin.Context) {
	var user models.User
	ctx.BindJSON(&user)
	ctx.JSON(http.StatusOK, nil)
}
