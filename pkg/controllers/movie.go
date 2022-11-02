package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type PostPersonalRatingInput struct {
	MovieID int64 `json:"movie_id" binding:"required"`
	UserID  int64 `json:"user_id" binding:"require"`
	Score   int64 `json:"score" binding:"require"`
}

func (s *Service) FindMovie(c *gin.Context) {
	var movie models.Movie
	if err := s.DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	s.DB.First(&movie)

	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

//func (s *Service) GetMovieByID(id int64) { }

// GET /movies
// Find a movie
func (s *Service) FindMovies(c *gin.Context) {
	var movies []models.Movie

	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error during reading json data!"})
		return
	}

	sort := gjson.Get(string(jsonData), "sort").Int()
	genres := gjson.Get(string(jsonData), "genres").Array()
	genresIDs := []int64{}
	for _, genreID := range genres {
		genresIDs = append(genresIDs, genreID.Int())
	}
	var subQuery *gorm.DB
	subQuery = nil
	if len(genresIDs) != 0{
		subQuery = s.DB.Select("movie_id").Where("genre_id in ?", genresIDs).Group("movie_id").Having("COUNT(distinct genre_id) = ?", len(genresIDs)).Model(&models.MovieGenres{})
	}

	limit := gjson.Get(string(jsonData), "limit").Int()
	from := gjson.Get(string(jsonData), "from").Int()
	var dptr *gorm.DB
	if subQuery != nil {
		dptr = s.DB.Limit(int(limit)).Offset(int(from)).Select("*").Joins("INNER JOIN (?) AS g ON g.movie_id = id", subQuery)
	} else {
		dptr = s.DB.Limit(int(limit)).Offset(int(from)).Select("*")
	}

	if gjson.Get(string(jsonData), "avgRateFrom").Exists(){
		dptr = dptr.Where("score >= ?", gjson.Get(string(jsonData), "avgRateFrom").Int())
	}
	if gjson.Get(string(jsonData), "avgRateTo").Exists(){
		dptr = dptr.Where("score <= ?", gjson.Get(string(jsonData), "avgRateTo").Int())
	}
	if gjson.Get(string(jsonData), "searchName").Exists(){
		dptr = dptr.Where("name LIKE ?", "%" + gjson.Get(string(jsonData), "searchName").String() + "%")
	}

	dptr = dptr.Group("id, name")
	if sort == 0{
		dptr = dptr.Order("year")
	} else {
		dptr = dptr.Order("score")
	}

	err = dptr.Find(&movies).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movies": movies})
}