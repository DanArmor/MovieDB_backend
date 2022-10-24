package models

type User struct {
	Id    int64  `json:"id" gorm:"primary_key"`
	Email string `json:"email"`
	Pass  string `json:"pass"`
}
