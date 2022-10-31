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

func (s *Service) CreateFees(c *gin.Context) {
	var input CreateFeesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	area := models.Fees{MovieID: input.MovieID, Value: input.Value, Currency: input.Currency, AreaID: input.AreaID}
	s.DB.Create(&area)

	c.JSON(http.StatusOK, area)
}

type CreateMovieGenreLinkInput struct {
	MovieID int64 `json:"movie_id" binding:"required"`
	GenreID int64 `json:"genre_id" binding:"required"`
}

func (s *Service) CreateMovieGenreLink(c *gin.Context) {
	var input CreateMovieGenreLinkInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movieGenreLink := models.MovieGenres{MovieID: input.MovieID, GenreID: input.GenreID}
	s.DB.Create(&movieGenreLink)

	c.JSON(http.StatusOK, movieGenreLink)
}

type CreateMovieInput struct {
	ExternalID      int64   `json:"external_id" binding:"required"`
	MovieTypeID     int64   `json:"movie_type_id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	AlternativeName string  `json:"alternative_name" binding:"required"`
	Description     string  `json:"description" binding:"required"`
	Year            int64   `json:"year" binding:"required"`
	StatusID        int64   `json:"status_id" binding:"required"`
	Duration        int64   `json:"duration"`
	Score           float32 `json:"score" binding:"required"`
	Votes           int64   `json:"votes" binding:"required"`
	AgeRating       int64   `json:"age_rating"`
	CountryID       int64   `json:"country_id" binding:"required"`
}

// POST /movies
func (s *Service) CreateMovie(c *gin.Context) {
	// Validate
	var input CreateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	s.DB.Create(&movie)

	c.JSON(http.StatusOK, movie)
}

type CreatePersonInput struct {
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
}

func (s *Service) CreatePerson(c *gin.Context) {
	var input CreatePersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := models.Person{Name: input.Name, NameEn: input.NameEn}
	s.DB.Create(&person)

	c.JSON(http.StatusOK, person)
}

type CreatePersonInMovieInput struct {
	MovieID      int64 `json:"movie_id" binding:"required"`
	PersonID     int64 `json:"person_id" binding:"required"`
	ProfessionID int64 `json:"profession_id" binding:"required"`
}

func (s *Service) CreatePersonInMovie(c *gin.Context) {
	var input CreatePersonInMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := models.PersonInMovie{MovieID: input.MovieID, PersonID: input.PersonID, ProfessionID: input.ProfessionID}
	s.DB.Create(&person)

	c.JSON(http.StatusOK, person)
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

func (s *Service) CreatePoster(c *gin.Context) {
	var input CreatePosterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	poster := models.Poster{Url: input.Url, MovieID: input.MovieID, PosterTypeID: input.PosterTypeID}
	s.DB.Create(&poster)

	c.JSON(http.StatusOK, poster)
}

type CreateSimpleInput struct {
	Type string `json:"type" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (s *Service) CreateSimpleData(c *gin.Context) {
	var input CreateSimpleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch {
	case input.Type == "areas":
		data := models.Area{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, data)
	case input.Type == "countries":
		data := models.Country{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, data)
	case input.Type == "genres":
		data := models.Genre{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, data)
	case input.Type == "movie_types":
		data := models.MovieType{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, data)
	case input.Type == "professions":
		data := models.Profession{NameEn: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, data)
	case input.Type == "poster_types":
		data := models.PosterType{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, data)
	case input.Type == "statuses":
		data := models.Status{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, data)
	}
}

func (s *Service) FindSimple(c *gin.Context) {
	result := map[string]interface{}{}
	var err error
	err = nil

	if _, has := c.GetQuery("type"); has == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No type to find"})
		return
	}
	if _, has := c.GetQuery("field"); has == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No field to find"})
		return
	}
	if _, has := c.GetQuery("value"); has == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No value to find"})
		return
	}

	t := c.Query("type")
	field := c.Query("field")
	value := c.Query("value")
	err = s.DB.Table(t).Take(&result, fmt.Sprintf("%s = ?", field), value).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found! " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Service) FindSimpleAll(c *gin.Context) {
	result := []map[string]interface{}{}
	var err error
	err = nil
	var hasField bool

	if _, has := c.GetQuery("type"); has == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No type to find"})
		return
	}
	if _, has := c.GetQuery("field"); has == false {
		hasField = has
	}
	if hasField {
		if _, has := c.GetQuery("value"); has == false {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No value to find"})
			return
		}
	}

	t := c.Query("type")
	if hasField {
		field := c.Query("field")
		value := c.Query("value")
		err = s.DB.Table(t).Find(&result, fmt.Sprintf("%s = ?", field), value).Error
	} else {
		err = s.DB.Table(t).Find(&result).Error
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found! " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{t: result})
}

func (s *Service) FindAdv(c *gin.Context) {
	result := map[string]interface{}{}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	t := gjson.Get(string(jsonData), "type")
	fields := gjson.Get(string(jsonData), "fieldNames")
	values := gjson.Get(string(jsonData), "values")

	dptr := s.DB.Table(t.String())
	for i := 0; i < len(fields.Array()); i++ {
		dptr = dptr.Where(fmt.Sprintf("%s = ?", fields.Array()[i].String()), values.Array()[i].String())
	}

	err = dptr.Take(&result).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found! " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
