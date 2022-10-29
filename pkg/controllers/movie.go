package controllers

import (
	"net/http"

	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/gin-gonic/gin"
)

type PostPersonalRatingInput struct {
	MovieID int64 `json:"movie_id" binding:"required"`
	UserID  int64 `json:"user_id" binding:"require"`
	Score   int64 `json:"score" binding:"require"`
}

func (s *Service) FindMovies(c *gin.Context) {
	var movies []models.Movie
	s.DB.Find(&movies)

	c.JSON(http.StatusOK, gin.H{"data": movies})
}

// GET /movies
// Find a movie
func (s *Service) FindMovie(c *gin.Context) {
	var movie models.Movie

	if err := s.DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": movie})
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
