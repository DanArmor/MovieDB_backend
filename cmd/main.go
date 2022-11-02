package main

import (
	//"net/http"

	"log"

	"github.com/DanArmor/MovieDB_backend/pkg/config"
	"github.com/DanArmor/MovieDB_backend/pkg/controllers"
	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/DanArmor/MovieDB_backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func SetupDataCache(s *controllers.Service) {
	s.MapCountry = make(map[int64]string)
	s.MapGenre = make(map[int64]string)
	s.MapProfs = make(map[int64]string)
	s.MapStatus = make(map[int64]string)
	s.MapType = make(map[int64]string)
	s.MapArea = make(map[int64]string)
	var g []models.Genre
	if err := s.DB.Find(&g).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	for _, genre := range g {
		s.MapGenre[genre.ID] = genre.Name
	}

	var mt []models.MovieType
	if err := s.DB.Find(&mt).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	for _, mtype := range mt {
		s.MapType[mtype.ID] = mtype.Name
	}

	var c []models.Country
	if err := s.DB.Find(&c).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	for _, country := range c {
		s.MapCountry[country.ID] = country.Name
	}

	var st []models.Status
	if err := s.DB.Find(&st).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	for _, status := range st {
		s.MapStatus[status.ID] = status.Name
	}

	var profs []models.Profession
	if err := s.DB.Find(&profs).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	for _, prof := range profs {
		s.MapProfs[prof.ID] = prof.NameEn
	}

	var areas []models.Area
	if err := s.DB.Find(&areas).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	for _, area := range areas {
		s.MapArea[area.ID] = area.Name
	}

	var posterType models.PosterType
	if err := s.DB.Where("name = ?", "preview").First(&posterType).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	s.PreviewID = posterType.ID
	posterType = models.PosterType{}
	if err := s.DB.Where("name = ?", "backdrop").First(&posterType).Error; err != nil {
		log.Fatalln("Error during cache setup")
	}
	s.BackdropID = posterType.ID

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
	}

	//Setup data
	SetupDataCache(&service)

	// Эндпоинты
	private := router.Group("/api")
	private.Use(service.ValidateToken)
	private.GET("/movies", service.FindMovies)
	private.GET("/movies/:id", service.FindMovie)
	private.POST("/rating/:id", service.UpdatePersonalScore)
	private.Static("res/img", "res/img")

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
	router.Run(config.Port)
}
