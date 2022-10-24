package models

type Movie struct {
	ID     int64  `json:"id" gorm:"primary_key"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
