package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	Port              string
	SecretKey         string
	EnableAutoMigrate bool
	S3Bucket          string
	S3Region          string
	S3Endpoint        string
	S3AccessKey       string
	S3SecretKey       string
	S3ForcePath       bool
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
		DatabaseURL:       dbUrl,
		Port:              port,
		SecretKey:         jwtSecret,
		EnableAutoMigrate: os.Getenv("ENABLE_GORM_AUTOMIGRATE") == "true",
		S3Bucket:          os.Getenv("STORAGE_BUCKET_NAME"),
		S3Region:          os.Getenv("STORAGE_REGION"),
		S3Endpoint:        os.Getenv("STORAGE_ENDPOINT"),
		S3AccessKey:       os.Getenv("STORAGE_ACCESS_KEY_ID"),
		S3SecretKey:       os.Getenv("STORAGE_SECRET_ACCESS_KEY"),
		S3ForcePath:       os.Getenv("STORAGE_FORCE_PATH_STYLE") == "true",
	}
}
