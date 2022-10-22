// models/setup.go

package models

import (
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
)

var DB *gorm.DB

func ConnectDatabase(sqlurl string) {
  database, err := gorm.Open(mysql.Open(sqlurl), &gorm.Config{})

  if err != nil {
    panic("Failed to connect to database!")
  }

  database.AutoMigrate(&Movie{})

  DB = database
}