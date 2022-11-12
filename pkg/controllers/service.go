package controllers

import (
	"github.com/DanArmor/MovieDB_backend/pkg/utils"
	"gorm.io/gorm"
)

type Service struct {
	Jwt        utils.JwtWrapper
	DB         *gorm.DB
	AdminPass  string
	BackdropID int64
	PreviewID  int64
	Domain     string
	BaseUrl    string
}
