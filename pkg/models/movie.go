package models

type Movie struct {
	ID              int64     `json:"id" gorm:"primary_key"`
	ExternalID      int64     `json:"external_id"`
	MovieTypeID     int64     `json:"movie_type_id"`
	Name            string    `json:"name"`
	AlternativeName string    `json:"alternative_name"`
	Description     string    `json:"description"`
	Year            int64     `json:"year"`
	StatusID        int64     `json:"status_id"`
	Duration        int64     `json:"duration"`
	Score           float32   `json:"score" gorm:"precision:3"`
	Votes           int64     `json:"votes"`
	AgeRating       int64     `json:"age_rating"`
	CountryID       int64     `json:"country_id"`
	Status          Status    `json:"-"`
	MovieType       MovieType `json:"-"`
	Country         Country   `json:"-"`
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
	ID           int64      `json:"id" gorm:"primary_key"`
	Url          string     `json:"url"`
	MovieID      int64      `json:"movie_id"`
	PosterTypeID int64      `json:"poster_type_id"`
	Movie        Movie      `json:"-"`
	PosterType   PosterType `json:"-"`
}

type PersonalRating struct {
	ID      int64 `json:"id" gorm:"primary_key"`
	MovieID int64 `json:"movie_id"`
	UserID  int64 `json:"user_id"`
	Score   int64 `json:"score"`
	User    User  `json:"-"`
	Movie   Movie `json:"-"`
}

type Area struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique"`
}

type Fees struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	MovieID  int64  `json:"movie_id"`
	Value    int64  `json:"value"`
	Currency string `json:"currency"`
	AreaID   int64  `json:"area_id"`
	Movie    Movie  `json:"-"`
	Area     Area   `json:"-"`
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
	Genre   Genre `json:"-"`
	Movie   Movie `json:"-"`
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
	ID           int64      `json:"id" gorm:"primary_key"`
	MovieID      int64      `json:"movie_id"`
	PersonID     int64      `json:"person_id"`
	ProfessionID int64      `json:"profession_id"`
	Person       Person     `json:"-"`
	Profession   Profession `json:"-"`
	Movie        Movie      `json:"-"`
}
