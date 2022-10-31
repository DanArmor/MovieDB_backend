// models/setup.go

package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(sqlurl string) *gorm.DB {
	database, err := gorm.Open(mysql.Open(sqlurl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&User{})

	database.AutoMigrate(&Area{})
	database.AutoMigrate(&Movie{})
	database.AutoMigrate(&MovieType{})
	database.AutoMigrate(&PosterType{})
	database.AutoMigrate(&Poster{})
	database.AutoMigrate(&PersonalRating{})
	database.AutoMigrate(&Fees{})
	database.AutoMigrate(&Status{})
	database.AutoMigrate(&Genre{})
	database.AutoMigrate(&MovieGenres{})
	database.AutoMigrate(&Country{})
	database.AutoMigrate(&Person{})
	database.AutoMigrate(&Profession{})
	database.AutoMigrate(&PersonInMovie{})

	return database
}
