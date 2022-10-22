package main

import (
  //"net/http"
  "log"
  "os"
  "github.com/gin-gonic/gin"
  "github.com/MovieDB_backend/pkg/config"
  "github.com/MovieDB_backend/pkg/models"
  "github.com/MovieDB_backend/pkg/controllers"
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

  // Подключаемся к ДБ
  models.ConnectDatabase(c.SqlUrl)

  // Создаем сервер
  r := gin.Default()

  // Эндпоинты
  r.GET("/movies", controllers.FindMovies)
  r.GET("/movies/:id", controllers.FindMovie)
  r.POST("/movies", controllers.CreateMovie)
  r.PATCH("/movies/:id", controllers.UpdateMovie)
  r.DELETE("/movies/:id", controllers.DeleteMovie)

  // Запускаем сервер
  r.Run(c.Port)
}