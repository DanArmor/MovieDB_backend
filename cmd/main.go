package main

import (
	//"net/http"
	"log"
	"os"

	"github.com/MovieDB_backend/pkg/config"
	"github.com/MovieDB_backend/pkg/controllers"
	"github.com/MovieDB_backend/pkg/models"
	"github.com/MovieDB_backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Логер
	infoLog := log.New(os.Stderr, "\033[32mINFO\033[0m\t", log.Ldate|log.Ltime)
	// Грузим конфигурацию
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed at config parse! ", err)
	}
	infoLog.Println("SqlUrl:", c.SqlUrl)
	infoLog.Println("Port:", c.Port)

	// Создаем сервер
	r := gin.Default()
	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTsecret,
		Issuer:          "MovieDB_backend",
		ExpirationHours: 24 * 7,
	}

	// Подключаемся к ДБ и т п
	service := controllers.Service{
		Jwt: jwt,
		DB:  models.ConnectDatabase(c.SqlUrl),
	}
	// Эндпоинты
	private := r.Group("/api")
	private.Use(service.ValidateToken)
	private.GET("/movies", service.FindMovies)
	private.GET("/movies/:id", service.FindMovie)
	private.POST("/movies", service.CreateMovie)
	private.PATCH("/movies/:id", service.UpdateMovie)
	private.DELETE("/movies/:id", service.DeleteMovie)

	public := r.Group("/auth")
	public.POST("/login", service.LoginUser)

	// Запускаем сервер
	r.Run(c.Port)
}
