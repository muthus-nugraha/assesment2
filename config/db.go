package config

import (
	"assignment2/app/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dbuser := "muthus"
	dbpass := "P@ssw0rd"
	dbhost := "localhost"
	dbname := "swagger"
	dbport := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbhost, dbuser, dbpass, dbname, dbport)
	db, errorDB := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errorDB != nil {
		panic("Connecting DB Failed")
	}
	return db
}

func Disconnect(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Disconect DB Failed")
	}
	dbSQL.Close()
}

func Migration(db *gorm.DB) {
	db.AutoMigrate(
		&models.Order{}, &models.Item{})
	fmt.Println("Migration DB Complated")
}
