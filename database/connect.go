package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	return os.Getenv(key)
}

func ConnectToDatabase() {
	var err error
	pUser := goDotEnvVariable("POSTGRES_USER")
	pPass := goDotEnvVariable("POSTGRES_PASSWORD")
	pDB := goDotEnvVariable("POSTGRES_DB")
	pHost := goDotEnvVariable("POSTGRES_HOST")

	urlDB := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432", pHost, pUser, pPass, pDB)

	db, err = gorm.Open(postgres.Open(urlDB), &gorm.Config{})

	if err != nil {
		log.Panic(err.Error())
	}

	log.Println("Your database is connected")
}
