package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	fmt.Print(os.Getenv(key))
	return os.Getenv(key)
}

func ConnectToDatabase() {
	var err error
	pUser := goDotEnvVariable("POSTGRES_USER")
	pPass := goDotEnvVariable("POSTGRES_PASSWORD")
	pDB := goDotEnvVariable("POSTGRES_DB")
	pHost := goDotEnvVariable("POSTGRES_HOST")

	urlDB := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432", pHost, pUser, pPass, pDB)

	DB, err = gorm.Open(postgres.Open(urlDB), &gorm.Config{})

	if err != nil {
		log.Panic(err.Error())
	}

	log.Println("Your database is connected")
}
