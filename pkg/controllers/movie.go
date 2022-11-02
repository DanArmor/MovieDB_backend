package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type PostPersonalRatingInput struct {
	MovieID int64 `json:"movie_id" binding:"required"`
	UserID  int64 `json:"user_id" binding:"require"`
	Score   int64 `json:"score" binding:"require"`
}

func (s *Service) FindMovie(c *gin.Context) {
	// was FindMovies
	var movies []models.Movie
	s.DB.Find(&movies)

	c.JSON(http.StatusOK, gin.H{"data": movies})
}

//func (s *Service) GetMovieByID(id int64) { }

// GET /movies
// Find a movie
func (s *Service) FindMovies(c *gin.Context) {
	var movies []models.Movie

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error during reading json data!"})
		return
	}

	//limit := gjson.Get(string(jsonData), "limit").Int()
	//from := gjson.Get(string(jsonData), "from").Int()
	//sort := gjson.Get(string(jsonData), "sort").Int()
	genres := gjson.Get(string(jsonData), "genres").Array()
	genresIDs := []int64{}
	for _, genreID := range genres {
		genresIDs = append(genresIDs, genreID.Int())
	}
	subQuery := s.DB.Select("movie_id").Where("genre_id in ?", genresIDs).Group("movie_id").Having("COUNT(distinct genre_id) = ?", len(genresIDs)).Model(&models.MovieGenres{})
	dptr := s.DB.Select("*").Joins("INNER JOIN (?) AS g ON g.movie_id = id", subQuery).Group("id, name")
	err = dptr.Find(&movies).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if gjson.Get(string(jsonData), "avgRateFrom").Exists(){

	}
	if gjson.Get(string(jsonData), "avgRateTo").Exists(){

	}
	if gjson.Get(string(jsonData), "searchName").Exists(){

	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}

type UpdateMovieInput struct {
	Name string `json:"name"`
}

// PATCH /movies/:id
// Update a movie
func (s *Service) UpdateMovie(c *gin.Context) {
	var movie models.Movie
	if err := s.DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate
	var input UpdateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateMovie := models.Movie{Name: input.Name}
	s.DB.Model(&movie).Updates(updateMovie)

	c.JSON(http.StatusOK, gin.H{"data": movie})
}

// DELETE /movies/:id
// Delete a movie
func (s *Service) DeleteMovie(c *gin.Context) {
	var movie models.Movie
	if err := s.DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	s.DB.Delete(&movie)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
