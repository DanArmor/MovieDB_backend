package controllers

import (
	"fmt"
	"net/http"

	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/gin-gonic/gin"
)

type CreateBudgetInput struct {
	MovieID  int64  `json:"movie_id" binding:"required"`
	Value    int64  `json:"value" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

func (s *Service) CreateBudget(c *gin.Context) {
	var input CreateBudgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	budget := models.Budget{MovieID: input.MovieID, Value: input.Value, Currency: input.Currency}
	s.DB.Create(&budget)

	c.JSON(http.StatusOK, gin.H{"data": budget})
}

type CreateFeesInput struct {
	MovieID  int64  `json:"movie_id" binding:"required"`
	Value    int64  `json:"value" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Area     string `json:"area" binding:"required"`
}

func (s *Service) CreateFees(c *gin.Context) {
	var input CreateFeesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	area := models.Fees{MovieID: input.MovieID, Value: input.Value, Currency: input.Currency, Area: input.Area}
	s.DB.Create(&area)

	c.JSON(http.StatusOK, gin.H{"data": area})
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

	c.JSON(http.StatusOK, gin.H{"data": movieGenreLink})
}

type CreateMovieInput struct {
	MovieTypeID         int64   `json:"movie_type_id" binding:"required"`
	Name                string  `json:"name" binding:"required"`
	Description         string  `json:"description" binding:"required"`
	Year                int64   `json:"year" binding:"required"`
	StatusID            int64   `json:"status_id" binding:"required"`
	Duration            int64   `json:"duration" binding:"required"`
	ProductionCompanyID int64   `json:"production_company_id" binding:"required"`
	Score               float32 `json:"score" binding:"required"`
	Votes               int64   `json:"votes" binding:"required"`
	AgeRating           int64   `json:"age_rating" binding:"required"`
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
		Duration: input.Duration, ProductionCompanyID: input.ProductionCompanyID,
		Score: input.Score, Votes: input.Votes}
	s.DB.Create(&movie)

	c.JSON(http.StatusOK, gin.H{"data": movie})
}

type CreatePersonInput struct {
	MovieTypeID int64  `json:"movie_type_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	NameEn      string `json:"name_en" binding:"required"`
	PhotoUrl    string `json:"photo_url" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (s *Service) CreatePerson(c *gin.Context) {
	var input CreatePersonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := models.Person{Name: input.Name, NameEn: input.NameEn, PhotoUrl: input.PhotoUrl, Description: input.Description}
	s.DB.Create(&person)

	c.JSON(http.StatusOK, gin.H{"data": person})
}

type CreatePersonInMovieInput struct {
	MovieID      int64  `json:"movie_id" binding:"required"`
	PersonID     string `json:"name" binding:"required"`
	ProfessionID string `json:"name_en" binding:"required"`
	Description  string `json:"description" binding:"required"`
}

func (s *Service) CreatePersonInMovie(c *gin.Context) {
	var input CreatePersonInMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := models.PersonInMovie{MovieID: input.MovieID, PersonID: input.PersonID, ProfessionID: input.ProfessionID, Description: input.Description}
	s.DB.Create(&person)

	c.JSON(http.StatusOK, gin.H{"data": person})
}

type CreateRatingInput struct {
	MovieID int64   `json:"movie_id" binding:"required"`
	RaterID int64   `json:"rating_id" binding:"required"`
	Score   float32 `json:"score" gorm:"precision:1" binding:"required"`
}

func (s *Service) CreateRating(c *gin.Context) {
	var input CreateRatingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating := models.Rating{MovieID: input.MovieID, RaterID: input.RaterID, Score: input.Score}
	s.DB.Create(&rating)

	c.JSON(http.StatusOK, gin.H{"data": rating})
}

type CreatePosterInput struct {
	MovieID int64  `json:"movie_id" binding:"required"`
	Url     string `json:"url" binding:"required"`
}

func (s *Service) CreatePoster(c *gin.Context) {
	var input CreatePosterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	poster := models.Poster{Url: input.Url, MovieID: input.MovieID}
	s.DB.Create(&poster)

	c.JSON(http.StatusOK, gin.H{"data": poster})
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
	case input.Type == "country":
		data := models.Country{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, gin.H{"data": data})
	case input.Type == "genre":
		data := models.Genre{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, gin.H{"data": data})
	case input.Type == "movie_type":
		data := models.MovieType{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, gin.H{"data": data})
	case input.Type == "production_company":
		data := models.ProductionCompany{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, gin.H{"data": data})
	case input.Type == "profession":
		data := models.Profession{NameEn: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, gin.H{"data": data})
	case input.Type == "rater":
		data := models.Rater{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, gin.H{"data": data})
	case input.Type == "status":
		data := models.Status{Name: input.Name}
		s.DB.Create(&data)
		c.JSON(http.StatusOK, gin.H{"data": data})
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

	t := c.Query("type");
	field := c.Query("field")
	value := c.Query("value")
	switch{
	case t == "genre":
		err = s.DB.Model(&models.Genre{}).First(&result, fmt.Sprintf("%s = ?", field), value).Error
	case t == "user":
		err = s.DB.Model(&models.User{}).First(&result, fmt.Sprintf("%s = ?", field), value).Error
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong type!"})
		return
	}
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found! " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}