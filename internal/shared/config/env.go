package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	SecretKey   string
}

func NewEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró archivo .env, usando variables de sistema")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL no está definida en el entorno")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET no configurado")
	}

	return &Config{
		DatabaseURL: dbUrl,
		Port:        port,
		SecretKey:   jwtSecret,
	}
}
