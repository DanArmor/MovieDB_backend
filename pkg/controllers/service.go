package controllers

import (
	"github.com/DanArmor/MovieDB_backend/pkg/utils"
	"gorm.io/gorm"
)

type Service struct {
	Jwt        utils.JwtWrapper
	DB         *gorm.DB
	AdminPass  string
	MapStatus  map[int64]string
	MapGenre   map[int64]string
	MapCountry map[int64]string
	MapType    map[int64]string
	MapProfs   map[int64]string
	MapArea    map[int64]string
	BackdropID int64
	PreviewID  int64
}
