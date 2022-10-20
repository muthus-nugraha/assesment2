package main

import (
	"assignment2/app/routers"
	"assignment2/config"

	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.Connect()
)

func main() {
	config.Migration(db)
	defer config.Disconnect(db)
	routers.Init()
}
