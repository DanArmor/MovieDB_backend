package main

import (
	//"net/http"

	"log"
	"os"
	"path/filepath"

	"github.com/DanArmor/MovieDB_backend/pkg/config"
	"github.com/DanArmor/MovieDB_backend/pkg/controllers"
	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/DanArmor/MovieDB_backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func LoadDataCache(service *controllers.Service, tableName string, mapPtr *map[int64]string) {
	result := []map[string]interface{}{}
	if err := service.DB.Table(tableName).Find(&result).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	for _, res := range result {
		id := res["id"].(int64)
		name, ok := res["name"].(string)
		if ok == false {
			name = res["name_en"].(string)
		}
		(*mapPtr)[id] = name
	}
}

func SetupDataCache(service *controllers.Service) {
	service.MapCountry = make(map[int64]string)
	service.MapGenre = make(map[int64]string)
	service.MapProfs = make(map[int64]string)
	service.MapStatus = make(map[int64]string)
	service.MapType = make(map[int64]string)
	service.MapArea = make(map[int64]string)
	LoadDataCache(service, "genres", &service.MapGenre)
	LoadDataCache(service, "movie_types", &service.MapType)
	LoadDataCache(service, "countries", &service.MapCountry)
	LoadDataCache(service, "statuses", &service.MapStatus)
	LoadDataCache(service, "professions", &service.MapProfs)
	LoadDataCache(service, "areas", &service.MapArea)

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
		BaseUrl: "https://" + config.Domain + ":8080",
	}

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


	res := router.Group("/res")
	res.Static("/img", "res/img")
	res.Static("/pdf", "res/pdf")

	public := router.Group("/auth")
	public.POST("/login", service.LoginUser)

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

	// Запускаем сервер
	router.RunTLS(config.Port, config.CertPath, config.KeyPath)
}
