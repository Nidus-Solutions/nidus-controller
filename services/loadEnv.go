// Esse arquivos apenas carrega as variáveis de ambiente do arquivo .env
package services

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	return os.Getenv(key)
}
