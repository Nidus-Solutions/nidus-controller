// Esse arquivos apenas carrega as variáveis de ambiente do arquivo .env
package services

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func LoadEnv(key string) string {
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}

	if godotenv.Load() == nil {
		return os.Getenv(key)
	}

	if err := envconfig.Process("", &key); err != nil {
		log.Fatal(err.Error())
	}

	return key
}
