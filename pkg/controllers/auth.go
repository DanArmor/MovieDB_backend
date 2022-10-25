package controllers

import (
	"net/http"

	"github.com/MovieDB_backend/pkg/models"
	"github.com/MovieDB_backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type LoginUserInput struct {
	Email string `json:"email" binding:"required"`
	Pass  string `json:"pass" binding:"required"`
}

func (s *Service) LoginUser(c *gin.Context) {
	var user models.User

	var input LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if result := s.DB.Where(&models.User{Email: input.Email}).First(&user); result.Error == nil {
		match := utils.CheckPasswordHash(input.Pass, user.Pass)
		if !match {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password!"})
			return
		}
		token, _ := s.Jwt.GenerateToken(user)
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}
	user.Email = input.Email
	user.Pass = utils.HashPassword(input.Pass)

	s.DB.Create(&user)
	token, _ := s.Jwt.GenerateToken(user)
	c.JSON(http.StatusOK, gin.H{"token": token})
	return
}

func (s *Service) ValidateToken(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token in header!"})
		return
	}
	claims, err := s.Jwt.ValidateToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	if result := s.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found!"})
		return
	}

	c.Next()
}
