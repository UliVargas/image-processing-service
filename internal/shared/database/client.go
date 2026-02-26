package database

import (
	"image-processing-service/internal/modules/session"
	"image-processing-service/internal/modules/user"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConection(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	db.AutoMigrate(user.User{}, session.Session{})

	log.Println("Conexión a Postgres exitosa")
	return db
}
