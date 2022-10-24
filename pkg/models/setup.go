// models/setup.go

package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase(sqlurl string) *gorm.DB {
	database, err := gorm.Open(mysql.Open(sqlurl), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Movie{})
	database.AutoMigrate(&User{})

	return database
}
