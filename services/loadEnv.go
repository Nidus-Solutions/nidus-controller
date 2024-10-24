// Esse arquivos apenas carrega as variáveis de ambiente do arquivo .env
package services

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(key string) string {

	if godotenv.Load(".env") != nil {
		if godotenv.Load("/etc/secrets/.env") != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return os.Getenv(key)
}
