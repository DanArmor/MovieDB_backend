package controllers

import (
	"net/http"

	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/DanArmor/MovieDB_backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type LoginUserInput struct {
	Email string `json:"email" binding:"required"`
	Pass  string `json:"pass" binding:"required"`
}

func (self *Service) LoginUser(context *gin.Context) {
	var user models.User

	var input LoginUserInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if result := self.DB.Where(&models.User{Email: input.Email}).First(&user); result.Error == nil {
		match := utils.CheckPasswordHash(input.Pass, user.Pass)
		if !match {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password!"})
			return
		}
		token, _ := self.Jwt.GenerateToken(user)
		context.JSON(http.StatusOK, gin.H{"token": token})
		return
	}
	user.Email = input.Email
	user.Pass = utils.HashPassword(input.Pass)

	self.DB.Create(&user)
	token, _ := self.Jwt.GenerateToken(user)
	context.JSON(http.StatusOK, gin.H{"token": token})
	return
}

func (self *Service) ValidateToken(context *gin.Context) {
	token := context.Request.Header.Get("token")
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token in header!"})
		return
	}
	claims, err := self.Jwt.ValidateToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	var user models.User
	if result := self.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found!"})
		return
	}

	context.Next()
}

func (self *Service) ValidateAdmin(context *gin.Context) {
	pass := context.Request.Header.Get("pass")
	if pass == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No admin password!"})
		return
	}

	if utils.CheckPasswordHash(pass, self.AdminPass) == false {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Wrong admin password!"})
		return
	}

	context.Next()
}
