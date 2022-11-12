package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (self *Service) GetHealth(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "Up and running!"})
}