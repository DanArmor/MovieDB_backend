package longmodels

import "github.com/DanArmor/MovieDB_backend/pkg/models"

type Movie struct {
	ID              int64                   `json:"id" gorm:"primary_key"`
	ExternalID      int64                   `json:"external_id"`
	Name            string                  `json:"name"`
	AlternativeName string                  `json:"alternative_name"`
	Year            int64                   `json:"year"`
	Score           float32                 `json:"score" gorm:"precision:3"`
	Votes           int64                   `json:"votes"`
	StatusID        int64                   `json:"-"`
	MovieTypeID     int64                   `json:"-"`
	CountryID       int64                   `json:"-"`
	Duration        int64                   `json:"duration"`
	Description     string                  `json:"description"`
	AgeRating       int64                   `json:"age_rating"`
	Genres          []models.Genre          `json:"genres" gorm:"many2many:movie_genres"`
	Posters         []models.Poster         `json:"poster"`
	MovieType       models.MovieType        `json:"movie_type"`
	Country         models.Country          `json:"country"`
	PersonalRating  []models.PersonalRating `json:"personal_rating"`
	Fees            []models.Fees           `json:"fees"`
	Status          models.Status           `json:"status"`
	Persons         []Person                `json:"persons" gorm:"many2many:person_in_movies"`
}

type Person struct {
	ID               int64  `json:"id" gorm:"primary_key"`
	Name             string `json:"name"`
	NameEn           string `json:"name_en"`
	ProfessionID     int64  `json:"-"`
	ProfessionNameEn string `json:"profession_name"`
}
