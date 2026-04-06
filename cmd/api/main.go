package main

import (
	"context"
	internalapi "image-processing-service/internal/api"
	"image-processing-service/internal/api/middleware"
	"image-processing-service/internal/modules/auth"
	"image-processing-service/internal/modules/file"
	"image-processing-service/internal/modules/session"
	"image-processing-service/internal/modules/user"
	tokenManager "image-processing-service/internal/shared/auth"
	"image-processing-service/internal/shared/config"
	"image-processing-service/internal/shared/database"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// ==========================================
	// Carga de variables de entorno
	// ==========================================
	cfg := config.NewEnv()

	// ==========================================
	// Configuración de S3
	// ==========================================
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(cfg.S3Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3AccessKey,
			cfg.S3SecretKey,
			"",
		)),
	)
	if err != nil {
		log.Fatal("Error cargando configuración de AWS:", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3Endpoint)
		o.UsePathStyle = cfg.S3ForcePath
	})

	// ==========================================
	// Configuración de JWT y base de datos
	// ==========================================
	m := tokenManager.NewTokenManager(cfg.SecretKey, time.Hour*24*7)
	db := database.NewConection(cfg.DatabaseURL, cfg.EnableAutoMigrate)

	// ==========================================
	// Wiring de módulos
	// ==========================================
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo)
	userHdl := user.NewHandler(userSvc)

	sessionRepo := session.NewRepository(db)
	sessionSvc := session.NewService(sessionRepo)

	authSvc := auth.NewService(userRepo, sessionSvc, m)
	authHdl := auth.NewHandler(authSvc)

	fileRepo := file.NewRepository(db)
	storage := file.NewS3Storage(s3Client, cfg.S3Bucket)
	fileSvc := file.NewService(fileRepo, storage)
	fileHdl := file.NewHandler(fileSvc)

	authMW := middleware.NewAuthMiddleware(m, sessionSvc)

	// ==========================================
	// Servidor
	// ==========================================
	addr := ":" + cfg.Port
	log.Printf("Iniciando servidor en el puerto %s", cfg.Port)
	http.ListenAndServe(addr, internalapi.NewRouter(authMW, authHdl, userHdl, fileHdl))
}
