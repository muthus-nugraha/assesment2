package repository

import (
	"gorm.io/gorm"
)

type dbConn struct {
	connection *gorm.DB
}
