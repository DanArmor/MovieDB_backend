package controllers

import (
	"github.com/MovieDB_backend/pkg/utils"
	"gorm.io/gorm"
)

type Service struct {
	Jwt utils.JwtWrapper
	DB  *gorm.DB
}
