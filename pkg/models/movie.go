package models

type Movie struct {
	ID          int64   `json:"id" gorm:"primary_key"`
	ExternalID  int64   `json:"external_id"`
	MovieTypeID int64   `json:"movie_type_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Year        int64   `json:"year"`
	StatusID    int64   `json:"status_id"`
	Duration    int64   `json:"duration"`
	Score       float32 `json:"my_rate" gorm:"precision:1"`
	Votes       int64   `json:"votes"`
	AgeRating   int64   `json:"age_rating"`
	Status      Status
	MovieType   MovieType
}

type MovieType struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
}

type PosterType struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
}

type Poster struct {
	ID           int64  `json:"id" gorm:"primary_key"`
	Url          string `json:"url"`
	MovieID      int64  `json:"movie_id"`
	PosterTypeID int64  `json:"poster_type_id"`
	Movie        Movie
	PosterType   PosterType
}

type PersonalRating struct {
	ID      int64 `json:"id" gorm:"primary_key"`
	MovieID int64 `json:"movie_id"`
	UserID  int64 `json:"user_id"`
	Score   int64 `json:"score"`
	User    User
	Movie   Movie
}

type Fees struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	MovieID  int64  `json:"movie_id"`
	Value    int64  `json:"value"`
	Currency string `json:"currency"`
	Area     string `json:"area"`
	Movie    Movie
}

type Status struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
}

type Genre struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
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
	Name string `json:"name" gorm:"unique"`
}

type Person struct {
	ID     int64  `json:"id" gorm:"primary_key"`
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
}

type Profession struct {
	ID     int64  `json:"id" gorm:"primary_key"`
	NameEn string `json:"name_en" gorm:"unique"`
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
