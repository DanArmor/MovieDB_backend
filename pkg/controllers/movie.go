package controllers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/DanArmor/MovieDB_backend/pkg/models"
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

type MovieInfoShort struct {
	ID              int64                 `json:"id"`
	ExternalID      int64                 `json:"external_id"`
	Name            string                `json:"name"`
	AlternativeName string                `json:"alternative_name"`
	Year            int64                 `json:"year"`
	Score           float32               `json:"score"`
	Votes           int64                 `json:"votes"`
	PersonalRating  models.PersonalRating `json:"personal_rating"`
	MovieType       models.MovieType      `json:"movie_type"`
	Country         models.Country        `json:"country"`
	Genres          []models.Genre        `json:"genres"`
	Preview         models.Poster         `json:"poster"`
}

type PersonLong struct {
	models.Person
	ProfessionName string `json:"profession"`
}

type FeesLong struct {
	models.Fees
	AreaName string `json:"area_name"`
}

type MovieInfoLong struct {
	MovieInfoShort
	Description string        `json:"description"`
	Fees        []FeesLong    `json:"fees"`
	Status      models.Status `json:"status"`
	Duration    int64         `json:"duration"`
	Persons     []PersonLong  `json:"persons"`
	AgeRating   int64         `json:"age_rating"`
	Backdrop    models.Poster `json:"backdrop"`
}

func MovieToShort(movie models.Movie) MovieInfoShort {
	return MovieInfoShort{
		ID: movie.ID, ExternalID: movie.ExternalID,
		Name: movie.Name, AlternativeName: movie.AlternativeName,
		Year: movie.Year, Score: movie.Score,
		Votes: movie.Votes,
	}
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

func (self *Service) GetMovieShortInfo(movie models.Movie, user_id int64) (MovieInfoShort, error) {
	info := MovieToShort(movie)

	if err := self.DB.Where("movie_id = ?", movie.ID).Where("poster_type_id = ?", self.PreviewID).First(&info.Preview).Error; err != nil {
		return MovieInfoShort{}, err
	}

	var mg []models.MovieGenres
	self.DB.Where("movie_id = ?", movie.ID).Find(&mg)
	for _, mg_link := range mg {
		info.Genres = append(info.Genres, models.Genre{ID: mg_link.GenreID, Name: self.MapGenre[mg_link.GenreID]})
	}

	self.DB.Where("movie_id = ?", movie.ID).Where("user_id = ?", user_id).First(&info.PersonalRating)

	info.Country = models.Country{ID: movie.CountryID, Name: self.MapCountry[movie.CountryID]}
	info.MovieType = models.MovieType{ID: movie.MovieTypeID, Name: self.MapType[movie.MovieTypeID]}

	return info, nil
}

func (self *Service) GetMovieLongInfo(movie models.Movie, user_id int64) (MovieInfoLong, error) {
	shortInfo, err := self.GetMovieShortInfo(movie, user_id)
	if err != nil {
		return MovieInfoLong{}, err
	}

	info := MovieInfoLong{MovieInfoShort: shortInfo}

	var poster models.Poster
	if err := self.DB.Where("movie_id = ?", movie.ID).Where("poster_type_id = ?", self.BackdropID).First(&poster).Error; err != nil {
		info.Backdrop = info.Preview
	} else {
		info.Backdrop = poster
	}

	info.Description = movie.Description
	info.Duration = movie.Duration
	info.Status = models.Status{ID: movie.StatusID, Name: self.MapStatus[movie.StatusID]}
	info.AgeRating = movie.AgeRating

	var fees []models.Fees
	self.DB.Where("movie_id = ?", movie.ID).Find(&fees)
	for _, fee := range fees {
		info.Fees = append(info.Fees, FeesLong{fee, self.MapArea[fee.AreaID]})
	}

	var pims []models.PersonInMovie
	if err := self.DB.Where("movie_id = ?", movie.ID).Find(&pims).Error; err != nil {
		return MovieInfoLong{}, err
	}
	for _, pim := range pims {
		var person models.Person
		if err := self.DB.Where("id = ?", pim.PersonID).First(&person).Error; err != nil {
			return MovieInfoLong{}, err
		}
		personLong := PersonLong{person, self.MapProfs[pim.ProfessionID]}
		info.Persons = append(info.Persons, personLong)
	}

	return info, nil
}

// GET /movies
// Find a movie
func (self *Service) FindMovies(context *gin.Context) {
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
	if gjson.Get(string(jsonStr), "searchName").Exists() {
		dptr = dptr.Where("name LIKE ?", "%"+gjson.Get(jsonStr, "searchName").String()+"%")
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

	var movies []models.Movie
	err = dptr.Find(&movies).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id := self.GetUserID(context)
	var moviesShorts []MovieInfoShort
	for _, movie := range movies {
		shortInfo, err := self.GetMovieShortInfo(movie, user_id)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		moviesShorts = append(moviesShorts, shortInfo)
	}

	context.JSON(http.StatusOK, gin.H{"movies": moviesShorts})
}

func (self *Service) FindMovie(context *gin.Context) {
	var movie models.Movie
	if err := self.DB.Where("id = ?", context.Param("id")).First(&movie).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	movieLong, err := self.GetMovieLongInfo(movie, self.GetUserID(context))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"movie": movieLong})
}

func (self *Service) UpdatePersonalScore(context *gin.Context) {
	jsonData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Error during reading json data!"})
		return
	}

	score := gjson.Get(string(jsonData), "score").Int()
	movie_id, _ := strconv.ParseInt(context.Param("id"), 10, 64)
	user_id := self.GetUserID(context)
	var rating models.PersonalRating
	if err := self.DB.Where("movie_id = ?", movie_id).Where("user_id = ?", user_id).First(&rating).Error; err != nil {
		rating = models.PersonalRating{MovieID: movie_id, UserID: user_id, Score: score}
		self.DB.Create(&rating)
	} else {
		rating.Score = score
		self.DB.Updates(&rating)
	}

	context.JSON(http.StatusOK, rating)
}

func (self *Service) GetGenres(context *gin.Context) {
	var genres []models.Genre
	self.DB.Find(&genres)

	context.JSON(http.StatusOK, gin.H{"genres" : genres})
}