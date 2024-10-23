package database

import (
	"fmt"
	"log"

	"github.com/jgb27/nidus-controller-projects/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Apenas conectando ao banco de dados
func ConnectToDatabase() {
	var err error

	// carregando variáveis de ambiente
	pUser := services.LoadEnv("POSTGRES_USER")
	pPass := services.LoadEnv("POSTGRES_PASSWORD")
	pDB := services.LoadEnv("POSTGRES_DB")
	pHost := services.LoadEnv("POSTGRES_HOST")

	// gerando url de conexão psql
	urlDB := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432", pHost, pUser, pPass, pDB)

	DB, err = gorm.Open(postgres.Open(urlDB), &gorm.Config{}) // conectando ao banco de dados

	if err != nil {
		log.Panic(err.Error())
	}

	log.Println("Your database is connected")
}
