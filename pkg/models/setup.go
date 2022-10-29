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

	database.AutoMigrate(&User{})

	database.AutoMigrate(&Movie{})
	database.AutoMigrate(&MovieType{})
	database.AutoMigrate(&Poster{})
	database.AutoMigrate(&PersonalRating{})
	database.AutoMigrate(&Rater{})
	database.AutoMigrate(&Rating{})
	database.AutoMigrate(&Budget{})
	database.AutoMigrate(&Fees{})
	database.AutoMigrate(&Status{})
	database.AutoMigrate(&Genre{})
	database.AutoMigrate(&MovieGenres{})
	database.AutoMigrate(&Country{})
	database.AutoMigrate(&Person{})
	database.AutoMigrate(&ProductionCompany{})
	database.AutoMigrate(&Profession{})
	database.AutoMigrate(&PersonInMovie{})

	return database
}
