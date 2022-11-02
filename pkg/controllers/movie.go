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

type FeesLong struct {
	models.Fees
	AreaName string `json:"area_name"`
}

func MovieToShort(movie models.Movie) MovieInfoShort {
	return MovieInfoShort{
		ID: movie.ID, ExternalID: movie.ExternalID,
		Name: movie.Name, AlternativeName: movie.AlternativeName,
		Year: movie.Year, Score: movie.Score,
		Votes: movie.Votes,
	}
}

func (s *Service) GetUserID(c *gin.Context) int64 {
	token := c.Request.Header.Get("token")
	claims, err := s.Jwt.ValidateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return -1
	}
	var user models.User
	if result := s.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found!"})
		return -1
	}
	return user.Id
}

func (s *Service) GetMovieShortInfo(movie models.Movie, user_id int64) MovieInfoShort {
	id := movie.ID
	info := MovieToShort(movie)

	s.DB.Where("movie_id = ?", id).Where("poster_type_id = ?", s.PreviewID).First(&info.Preview)

	var mg []models.MovieGenres
	s.DB.Where("movie_id = ?", id).Find(&mg)
	for _, mg_link := range mg {
		info.Genres = append(info.Genres, models.Genre{ID: mg_link.GenreID, Name: s.MapGenre[mg_link.GenreID]})
	}

	s.DB.Where("movie_id = ?", id).Where("user_id = ?", user_id).First(&info.PersonalRating)

	info.Country = models.Country{ID: movie.CountryID, Name: s.MapCountry[movie.CountryID]}
	info.MovieType = models.MovieType{ID: movie.MovieTypeID, Name: s.MapType[movie.MovieTypeID]}

	return info
}

func (s *Service) GetMovieLongInfo(movie models.Movie, user_id int64) MovieInfoLong {
	id := movie.ID
	info := MovieInfoLong{MovieInfoShort: s.GetMovieShortInfo(movie, user_id)}

	var poster models.Poster
	if err := s.DB.Where("movie_id = ?", id).Where("poster_type_id = ?", s.BackdropID).First(&poster).Error; err != nil {
		info.Backdrop = info.Preview
	} else {
		info.Backdrop = poster
	}

	info.Description = movie.Description
	info.Duration = movie.Duration
	info.Status = models.Status{ID: movie.StatusID, Name: s.MapStatus[movie.StatusID]}
	info.AgeRating = movie.AgeRating

	var fees []models.Fees
	s.DB.Where("movie_id = ?", id).Find(&fees)
	for _, fee := range fees {
		info.Fees = append(info.Fees, FeesLong{fee, s.MapArea[fee.AreaID]})
	}

	var pims []models.PersonInMovie
	if err := s.DB.Where("movie_id = ?", id).Find(&pims).Error; err != nil {
		return MovieInfoLong{}
	}
	for _, pim := range pims {
		var person models.Person
		if err := s.DB.Where("id = ?", pim.PersonID).First(&person).Error; err != nil {
			return MovieInfoLong{}
		}
		personLong := PersonLong{person, s.MapProfs[pim.ProfessionID]}
		info.Persons = append(info.Persons, personLong)
	}

	return info
}

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
	if len(genresIDs) != 0 {
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

	if gjson.Get(string(jsonData), "avgRateFrom").Exists() {
		dptr = dptr.Where("score >= ?", gjson.Get(string(jsonData), "avgRateFrom").Int())
	}
	if gjson.Get(string(jsonData), "avgRateTo").Exists() {
		dptr = dptr.Where("score <= ?", gjson.Get(string(jsonData), "avgRateTo").Int())
	}
	if gjson.Get(string(jsonData), "searchName").Exists() {
		dptr = dptr.Where("name LIKE ?", "%"+gjson.Get(string(jsonData), "searchName").String()+"%")
	}

	dptr = dptr.Group("id, name")
	if sort == 0 {
		dptr = dptr.Order("year")
	} else {
		dptr = dptr.Order("score")
	}

	err = dptr.Find(&movies).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id := s.GetUserID(c)
	var moviesShorts []MovieInfoShort
	for _, movie := range movies {
		moviesShorts = append(moviesShorts, s.GetMovieShortInfo(movie, user_id))
	}

	c.JSON(http.StatusOK, gin.H{"movies": moviesShorts})
}

func (s *Service) FindMovie(c *gin.Context) {
	var movie models.Movie
	if err := s.DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// user id
	token := c.Request.Header.Get("token")
	claims, err := s.Jwt.ValidateToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	var user models.User
	if result := s.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found!"})
		return
	}

	movieLong := s.GetMovieLongInfo(movie, user.Id)

	c.JSON(http.StatusOK, gin.H{"movie": movieLong})
}

func (s *Service) UpdatePersonalScore(c *gin.Context) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error during reading json data!"})
		return
	}

	score := gjson.Get(string(jsonData), "score").Int()
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user_id := s.GetUserID(c)
	var rating models.PersonalRating
	if err := s.DB.Where("movie_id = ?", id).Where("user_id = ?", user_id).First(&rating).Error; err != nil {
		rating = models.PersonalRating{MovieID:  id, UserID: user_id, Score: score}
		s.DB.Create(&rating)
	} else {
		rating.Score = score
		s.DB.Updates(&rating)
	}

	c.JSON(http.StatusOK, rating)
}