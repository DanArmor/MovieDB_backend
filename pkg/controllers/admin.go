package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"

	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/gin-gonic/gin"
)

type CreateFeesInput struct {
	MovieID  int64  `json:"movie_id" binding:"required"`
	Value    int64  `json:"value" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	AreaID   int64  `json:"area_id" binding:"required"`
}

func (self *Service) CreateFees(context *gin.Context) {
	var input CreateFeesInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	area := models.Fees{MovieID: input.MovieID, Value: input.Value, Currency: input.Currency, AreaID: input.AreaID}
	self.DB.Create(&area)

	context.JSON(http.StatusOK, area)
}

type CreateMovieGenreLinkInput struct {
	MovieID int64 `json:"movie_id" binding:"required"`
	GenreID int64 `json:"genre_id" binding:"required"`
}

func (self *Service) CreateMovieGenreLink(context *gin.Context) {
	var input CreateMovieGenreLinkInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movieGenreLink := models.MovieGenres{MovieID: input.MovieID, GenreID: input.GenreID}
	self.DB.Create(&movieGenreLink)

	context.JSON(http.StatusOK, movieGenreLink)
}

type CreateMovieInput struct {
	ExternalID      int64   `json:"external_id" binding:"required"`
	MovieTypeID     int64   `json:"movie_type_id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	AlternativeName string  `json:"alternative_name" binding:"required"`
	Description     string  `json:"description" binding:"required"`
	Year            int64   `json:"year"`
	StatusID        int64   `json:"status_id" binding:"required"`
	Duration        int64   `json:"duration"`
	Score           float32 `json:"score"`
	Votes           int64   `json:"votes" binding:"required"`
	AgeRating       int64   `json:"age_rating"`
	CountryID       int64   `json:"country_id" binding:"required"`
}

type UpdateMovieInput struct {
	ExternalID      int64   `json:"external_id"`
	MovieTypeID     int64   `json:"movie_type_id"`
	Name            string  `json:"name"`
	AlternativeName string  `json:"alternative_name"`
	Description     string  `json:"description"`
	Year            int64   `json:"year"`
	StatusID        int64   `json:"status_id"`
	Duration        int64   `json:"duration"`
	Score           float32 `json:"score"`
	Votes           int64   `json:"votes"`
	AgeRating       int64   `json:"age_rating"`
	CountryID       int64   `json:"country_id"`
}

// PATCH /movies/:id
// Update a movie
func (self *Service) UpdateMovie(context *gin.Context) {
	var movie models.Movie
	if err := self.DB.Where("external_id = ?", context.Query("external_id")).First(&movie).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate
	var input UpdateMovieInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateMovie := models.Movie{Name: input.Name, Description: input.Description,
		Year: input.Year, StatusID: input.StatusID,
		Duration: input.Duration,
		Score:    input.Score, Votes: input.Votes,
		ExternalID: input.ExternalID, AlternativeName: input.AlternativeName,
		CountryID: input.CountryID, MovieTypeID: input.MovieTypeID,
		AgeRating: input.AgeRating}
	if err := self.DB.Model(&movie).Updates(updateMovie).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, movie)
}

// POST /movies
func (self *Service) CreateMovie(context *gin.Context) {
	// Validate
	var input CreateMovieInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create
	movie := models.Movie{Name: input.Name, Description: input.Description,
		Year: input.Year, StatusID: input.StatusID,
		Duration: input.Duration,
		Score:    input.Score, Votes: input.Votes,
		ExternalID: input.ExternalID, AlternativeName: input.AlternativeName,
		CountryID: input.CountryID, MovieTypeID: input.MovieTypeID,
		AgeRating: input.AgeRating}
	self.DB.Create(&movie)

	context.JSON(http.StatusOK, movie)
}

type CreatePersonInput struct {
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
}

func (self *Service) CreatePerson(context *gin.Context) {
	var input CreatePersonInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := models.Person{Name: input.Name, NameEn: input.NameEn}
	self.DB.Create(&person)

	context.JSON(http.StatusOK, person)
}

type CreatePersonInMovieInput struct {
	MovieID      int64 `json:"movie_id" binding:"required"`
	PersonID     int64 `json:"person_id" binding:"required"`
	ProfessionID int64 `json:"profession_id" binding:"required"`
}

func (self *Service) CreatePersonInMovie(context *gin.Context) {
	var input CreatePersonInMovieInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := models.PersonInMovie{MovieID: input.MovieID, PersonID: input.PersonID, ProfessionID: input.ProfessionID}
	self.DB.Create(&person)

	context.JSON(http.StatusOK, person)
}

type CreateRatingInput struct {
	MovieID int64   `json:"movie_id" binding:"required"`
	RaterID int64   `json:"rating_id" binding:"required"`
	Score   float32 `json:"score" gorm:"precision:1" binding:"required"`
}

type CreatePosterInput struct {
	MovieID      int64  `json:"movie_id" binding:"required"`
	Url          string `json:"url" binding:"required"`
	PosterTypeID int64  `json:"poster_type_id" binding:"required"`
}

func (self *Service) CreatePoster(context *gin.Context) {
	var input CreatePosterInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	poster := models.Poster{Url: input.Url, MovieID: input.MovieID, PosterTypeID: input.PosterTypeID}
	self.DB.Create(&poster)

	context.JSON(http.StatusOK, poster)
}

type CreateSimpleInput struct {
	Type string `json:"type" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (self *Service) CreateSimpleData(context *gin.Context) {
	var input CreateSimpleInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch {
	case input.Type == "areas":
		data := models.Area{Name: input.Name}
		self.DB.Create(&data)
		context.JSON(http.StatusOK, data)
	case input.Type == "countries":
		data := models.Country{Name: input.Name}
		self.DB.Create(&data)
		context.JSON(http.StatusOK, data)
	case input.Type == "genres":
		data := models.Genre{Name: input.Name}
		self.DB.Create(&data)
		context.JSON(http.StatusOK, data)
	case input.Type == "movie_types":
		data := models.MovieType{Name: input.Name}
		self.DB.Create(&data)
		context.JSON(http.StatusOK, data)
	case input.Type == "professions":
		data := models.Profession{NameEn: input.Name}
		self.DB.Create(&data)
		context.JSON(http.StatusOK, data)
	case input.Type == "poster_types":
		data := models.PosterType{Name: input.Name}
		self.DB.Create(&data)
		context.JSON(http.StatusOK, data)
	case input.Type == "statuses":
		data := models.Status{Name: input.Name}
		self.DB.Create(&data)
		context.JSON(http.StatusOK, data)
	}
}

func (self *Service) FindSimple(context *gin.Context) {
	result := map[string]interface{}{}
	var err error
	err = nil

	if _, has := context.GetQuery("type"); has == false {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No type to find"})
		return
	}
	if _, has := context.GetQuery("field"); has == false {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No field to find"})
		return
	}
	if _, has := context.GetQuery("value"); has == false {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No value to find"})
		return
	}

	t := context.Query("type")
	field := context.Query("field")
	value := context.Query("value")
	err = self.DB.Table(t).Take(&result, fmt.Sprintf("%s = ?", field), value).Error

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found! " + err.Error()})
		return
	}

	context.JSON(http.StatusOK, result)
}

func (self *Service) FindSimpleAll(context *gin.Context) {
	result := []map[string]interface{}{}
	var err error
	err = nil
	var hasField bool

	if _, has := context.GetQuery("type"); has == false {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No type to find"})
		return
	}
	if _, has := context.GetQuery("field"); has == false {
		hasField = has
	}
	if hasField {
		if _, has := context.GetQuery("value"); has == false {
			context.JSON(http.StatusBadRequest, gin.H{"error": "No value to find"})
			return
		}
	}

	t := context.Query("type")
	if hasField {
		field := context.Query("field")
		value := context.Query("value")
		err = self.DB.Table(t).Find(&result, fmt.Sprintf("%s = ?", field), value).Error
	} else {
		err = self.DB.Table(t).Find(&result).Error
	}

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found! " + err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{t: result})
}

func (self *Service) FindAdv(context *gin.Context) {
	result := map[string]interface{}{}
	jsonData, err := ioutil.ReadAll(context.Request.Body)
	t := gjson.Get(string(jsonData), "type")
	fields := gjson.Get(string(jsonData), "fieldNames")
	values := gjson.Get(string(jsonData), "values")

	dptr := self.DB.Table(t.String())
	for i := 0; i < len(fields.Array()); i++ {
		dptr = dptr.Where(fmt.Sprintf("%s = ?", fields.Array()[i].String()), values.Array()[i].String())
	}

	err = dptr.Take(&result).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found! " + err.Error()})
		return
	}

	context.JSON(http.StatusOK, result)
}
