package controllers

import (
	"net/http"

	"github.com/MovieDB_backend/pkg/models"
	"github.com/gin-gonic/gin"
)

func FindMovies(c *gin.Context) {
	var movies []models.Movie
	models.DB.Find(&movies)

	c.JSON(http.StatusOK, gin.H{"data": movies})
}

type CreateMovieInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

// POST /movies
func CreateMovie(c *gin.Context) {
	// Validate
	var input CreateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create
	movie := models.Movie{Title: input.Title, Author: input.Author}
	models.DB.Create(&movie)

	c.JSON(http.StatusOK, gin.H{"data": movie})
}

// GET /movies
// Find a movie
func FindMovie(c *gin.Context) {
  var movie models.Movie

  if err := models.DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
    return
  }

  c.JSON(http.StatusOK, gin.H{"data": movie})
}

type UpdateMovieInput struct {
  Title  string `json:"title"`
  Author string `json:"author"`  
}

// PATCH /movies/:id
// Update a movie
func UpdateMovie(c *gin.Context) {
  var movie models.Movie
  if err := models.DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
    return
  }

  // Validate
  var input UpdateMovieInput
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  updateMovie := models.Movie{Title: input.Title, Author: input.Author}
  models.DB.Model(&movie).Updates(updateMovie)

  c.JSON(http.StatusOK, gin.H{"data": movie})
}

// DELETE /movies/:id
// Delete a movie
func DeleteMovie(c *gin.Context) {
  var movie models.Movie
  if err := models.DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
    return
  }

  models.DB.Delete(&movie)

  c.JSON(http.StatusOK, gin.H{"data": true})
}