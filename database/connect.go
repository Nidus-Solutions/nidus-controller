package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	return os.Getenv(key)
}

func ConnectToDatabase() *sql.DB {
	var err error
	pUser := goDotEnvVariable("POSTGRES_USER")
	pPass := goDotEnvVariable("POSTGRES_PASSWORD")
	pDB := goDotEnvVariable("POSTGRES_DB")

	urlDB := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", pUser, pPass, pDB)

	db, err = sql.Open("postgres", urlDB)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("You connected to yout database.")

	return db
}
