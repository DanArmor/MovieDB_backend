package shortmodels

import "github.com/DanArmor/MovieDB_backend/pkg/models"

type MovieInfoShort struct {
	ID              int64   `json:"id"`
	ExternalID      int64   `json:"external_id"`
	Name            string  `json:"name"`
	AlternativeName string  `json:"alternative_name"`
	Year            int64   `json:"year"`
	Score           float32 `json:"score"`
	Votes           int64   `json:"votes"`
}

type Movie struct {
	ID              int64                   `json:"id" gorm:"primary_key"`
	ExternalID      int64                   `json:"external_id"`
	Name            string                  `json:"name"`
	AlternativeName string                  `json:"alternative_name"`
	Year            int64                   `json:"year"`
	Score           float32                 `json:"score" gorm:"precision:3"`
	Votes           int64                   `json:"votes"`
	MovieTypeID     int64                   `json:"-"`
	CountryID       int64                   `json:"-"`
	Genres          []models.Genre          `json:"genres" gorm:"many2many:movie_genres"`
	Posters         []models.Poster         `json:"poster"`
	MovieType       models.MovieType        `json:"movie_type"`
	Country         models.Country          `json:"country"`
	PersonalRating  []models.PersonalRating `json:"personal_rating"`
}
