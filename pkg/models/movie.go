package models

import "time"

type Movie struct {
	ID                  int64  `json:"id" gorm:"primary_key"`
	MovieTypeID         int64  `json:"movie_type_id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Year                int64  `json:"year"`
	StatusID            int64  `json:"status_id"`
	MovieLength         int64  `json:"length"`
	ProductionCompanyID int64  `json:"production_company_id"`
	Status              Status
	MovieType           MovieType
	ProductionCompany   ProductionCompany
}

type MovieType struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type Poster struct {
	ID      int64  `json:"id" gorm:"primary_key"`
	Url     string `json:"url"`
	MovieID int64  `json:"movie_id"`
	Movie   Movie
}

type Rater struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type Rating struct {
	ID      int64   `json:"id" gorm:"primary_key"`
	MovieID int64   `json:"movie_id"`
	RaterID int64   `json:"rating_id"`
	Score   float32 `json:"score" gorm:"precision:1"`
	Movie   Movie
	Rater   Rater
}

type Budget struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	MovieID  int64  `json:"movie_id"`
	Value    int64  `json:"value"`
	Currency string `json:"currency"`
	Movie    Movie
}

type Fees struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	MovieID  int64  `json:"movie_id"`
	Value    int64  `json:"value"`
	Currency string `json:"currency"`
	Area     string `json:"area"`
	Movie    Movie
}

type PremierType struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type Premier struct {
	ID            int64     `json:"id" gorm:"primary_key"`
	MovieID       int64     `json:"movie_id"`
	PremierTypeID int64     `json:"premier_type"`
	Date          time.Time `json:"date"`
	PremierType   PremierType
	Movie         Movie
}

type Status struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type Genre struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type MovieGenres struct {
	ID      int64 `json:"id" gorm:"primary_key"`
	GenreID int64 `json:"genre_id"`
	MovieID int64 `json:"movie_id"`
	Genre   Genre
	Movie   Movie
}

type Country struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type Person struct {
	ID          int64  `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	NameEn      string `json:"name_en"`
	PhotoUrl    string `json:"photo_url"`
	Description string `json:"description"`
}

type ProductionCompany struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type Profession struct {
	ID     int64  `json:"id" gorm:"primary_key"`
	NameEn string `json:"name_en"`
}

type PersonInMovie struct {
	ID           int64  `json:"id" gorm:"primary_key"`
	MovieID      int64  `json:"movie_id"`
	PersonID     string `json:"name"`
	ProfessionID string `json:"name_en"`
	Description  string `json:"description"`
	Person       Person
	Profession   Profession
	Movie        Movie
}
