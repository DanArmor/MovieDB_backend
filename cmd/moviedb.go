package main

import (
	//"net/http"

	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/DanArmor/MovieDB_backend/pkg/config"
	"github.com/DanArmor/MovieDB_backend/pkg/controllers"
	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/DanArmor/MovieDB_backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func SetupDataCache(service *controllers.Service) {
	var posterType models.PosterType
	if err := service.DB.Where("name = ?", "preview").First(&posterType).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	service.PreviewID = posterType.ID
	posterType = models.PosterType{}
	if err := service.DB.Where("name = ?", "backdrop").First(&posterType).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	service.BackdropID = posterType.ID
}

func RemoveContents(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            return err
        }
    }
    return nil
}

func main() {
	// Логер
	// Грузим конфигурацию
	config, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("failed at config parse! ", err)
	}

	config.AdminPass = utils.HashPassword(config.AdminPass)
	// Создаем сервер
	router := gin.Default()
	jwt := utils.JwtWrapper{
		SecretKey:       config.JWTsecret,
		Issuer:          "MovieDB_backend",
		ExpirationHours: 24 * 7,
	}

	// Подключаемся к ДБ и т п
	service := controllers.Service{
		Jwt:       jwt,
		DB:        models.ConnectDatabase(config.SqlUrl),
		AdminPass: config.AdminPass,
		Domain: config.Domain,
		BaseUrl: "https://" + config.Domain,
	}
	router.SetTrustedProxies([]string{"127.0.0.1:80"})
	//Setup data
	SetupDataCache(&service)
	RemoveContents("./res/pdf")

	// Эндпоинты
	api := router.Group("/api")
	api.Use(service.ValidateToken)
	api.POST("/movies", service.FindMovies)
	api.GET("/movies/:id", service.FindMovie)
	api.GET("/genres", service.GetGenres)
	api.POST("/rating/:id", service.UpdatePersonalScore)
	api.GET("/pdf/:id", service.GetPDF)

	auth := router.Group("/auth")
	auth.POST("/login", service.LoginUser)

	public := router.Group("/public")
	public.GET("/health", service.GetHealth)

	admin := router.Group("/admin")
	admin.Use(service.ValidateAdmin)
	admin.POST("/simple", service.CreateSimpleData)

	admin.POST("/fees", service.CreateFees)
	admin.POST("/movie_genres", service.CreateMovieGenreLink)
	admin.POST("/movies", service.CreateMovie)
	admin.POST("/people", service.CreatePerson)
	admin.POST("/person_in_movies", service.CreatePersonInMovie)
	admin.POST("/posters", service.CreatePoster)

	admin.PATCH("/updateMovie", service.UpdateMovie)

	admin.GET("/find", service.FindSimple)
	admin.GET("/findAll", service.FindSimpleAll)
	admin.GET("/findAdv", service.FindAdv)
	srv := &http.Server{
		Addr: "127.0.0.1" + config.Port,
		Handler: router,
	}
	// Запускаем сервер
	go func(){
		if err := srv.ListenAndServeTLS(config.CertPath, config.KeyPath); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGUSR1)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
