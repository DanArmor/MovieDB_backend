package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	longmodels "github.com/DanArmor/MovieDB_backend/pkg/long_models"
	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/DanArmor/MovieDB_backend/pkg/short_models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

const SORT_YEAR = 0
const SORT_SCORE = 1

type PostPersonalRatingInput struct {
	MovieID int64 `json:"movie_id" binding:"required"`
	UserID  int64 `json:"user_id" binding:"require"`
	Score   int64 `json:"score" binding:"require"`
}

func (self *Service) GetUserID(context *gin.Context) int64 {
	token := context.Request.Header.Get("token")
	claims, err := self.Jwt.ValidateToken(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return -1
	}
	var user models.User
	if result := self.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found!"})
		return -1
	}
	return user.Id
}

// GET /movies
// Find a movie
func (self *Service) FindMovies(context *gin.Context) {
	defer context.Request.Body.Close()
	jsonBytes, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error during reading json data!"})
		return
	}
	jsonStr := string(jsonBytes)

	sort := gjson.Get(jsonStr, "sort").Int()
	genres := gjson.Get(jsonStr, "genres").Array()
	genresIDs := []int64{}
	for _, genreID := range genres {
		genresIDs = append(genresIDs, genreID.Int())
	}
	var subQuery *gorm.DB
	subQuery = nil
	if len(genresIDs) != 0 {
		subQuery = self.DB.Select("movie_id").Where("genre_id in ?", genresIDs).Group("movie_id").Having("COUNT(distinct genre_id) = ?", len(genresIDs)).Model(&models.MovieGenres{})
	}

	limit := gjson.Get(jsonStr, "limit").Int()
	from := gjson.Get(jsonStr, "from").Int()
	var dptr *gorm.DB
	if subQuery != nil {
		dptr = self.DB.Limit(int(limit)).Offset(int(from)).Select("*").Joins("INNER JOIN (?) AS g ON g.movie_id = id", subQuery)
	} else {
		dptr = self.DB.Limit(int(limit)).Offset(int(from)).Select("*")
	}

	if gjson.Get(string(jsonStr), "avgRateFrom").Exists() {
		dptr = dptr.Where("score >= ?", gjson.Get(jsonStr, "avgRateFrom").Int())
	}
	if gjson.Get(string(jsonStr), "avgRateTo").Exists() {
		dptr = dptr.Where("score <= ?", gjson.Get(jsonStr, "avgRateTo").Int())
	}
	dptr = dptr.Group("id, name")
	if sort == SORT_YEAR {
		dptr = dptr.Order("year DESC")
	} else if sort == SORT_SCORE {
		dptr = dptr.Order("score DESC")
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Wrong sort enum"})
		return
	}
	dptr = dptr.Order("id")
	if gjson.Get(string(jsonStr), "searchName").Exists() {
		dptr = dptr.Where("name LIKE ?", "%"+gjson.Get(jsonStr, "searchName").String()+"%").Or("alternative_name LIKE ?", "%"+gjson.Get(jsonStr, "searchName").String()+"%")
	}

	user_id := self.GetUserID(context)
	var movies []shortmodels.Movie
	dptr = dptr.Preload("Country").Preload("MovieType").
			Preload("Posters", "poster_type_id = ?", self.PreviewID).
			Preload("Genres").Preload("PersonalRating", "user_id = ?", user_id).Find(&movies)
	err = dptr.Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	context.JSON(http.StatusOK, gin.H{"movies": movies})
}

func (self *Service) FindMovie(context *gin.Context) {
	var movie longmodels.Movie
	user_id := self.GetUserID(context)
	self.DB.Preload("Country").Preload("MovieType").
			Preload("Posters", "poster_type_id = ?", self.PreviewID).
			Preload("Genres").Preload("PersonalRating", "user_id = ?", user_id).
			Preload("Status").Preload("Fees.Area").Where("id = ?", context.Param("id")).First(&movie)
	dptr := self.DB.Table("person_in_movies pim").Where("pim.movie_id = ?", context.Param("id")).
	Joins("JOIN professions prof on prof.id = pim.profession_id").
	Joins("JOIN people ppl on ppl.id = pim.person_id").
	Select("ppl.id as id, ppl.name as name, ppl.name_en as name_en, prof.name_en as profession_name_en")
	num, err := strconv.Atoi(context.Query("persons_count"))
	if err != nil {
		num = 0
	}
	if(num != 0){
		dptr = dptr.Limit(num)
	}
	dptr.Find(&movie.Persons)

	context.JSON(http.StatusOK, gin.H{"movie": movie})
}

func (self *Service) UpdatePersonalScore(context *gin.Context) {
	defer context.Request.Body.Close()
	jsonData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error during reading json data!"})
		return
	}

	score := gjson.Get(string(jsonData), "score").Int()
	movie_id, _ := strconv.ParseInt(context.Param("id"), 10, 64)
	user_id := self.GetUserID(context)
	var ratingArr []models.PersonalRating
	if err := self.DB.Where("movie_id = ?", movie_id).Where("user_id = ?", user_id).First(&ratingArr).Error; err != nil {
		rating := models.PersonalRating{MovieID: movie_id, UserID: user_id, Score: score}
		ratingArr = append(ratingArr, rating)
		self.DB.Create(&rating)
	} else {
		ratingArr[0].Score = score
		self.DB.Updates(&ratingArr[0])
	}

	context.JSON(http.StatusOK, ratingArr)
}

func (self *Service) GetGenres(context *gin.Context) {
	var genres []models.Genre
	self.DB.Find(&genres)

	context.JSON(http.StatusOK, gin.H{"genres": genres})
}
