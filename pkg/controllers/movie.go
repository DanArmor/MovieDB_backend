package controllers

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	longmodels "github.com/DanArmor/MovieDB_backend/pkg/long_models"
	"github.com/DanArmor/MovieDB_backend/pkg/models"
	"github.com/DanArmor/MovieDB_backend/pkg/short_models"
	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
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
			Preload("Status").Preload("Fees").Preload("Persons", func(tx *gorm.DB) *gorm.DB{
				num, err := strconv.Atoi(context.Query("persons_count"))
				if err != nil {
					num = 0
				}
				return tx.Limit(num)
			}).Preload("Persons.Profession").First(&movie)
				fmt.Print(movie)

	context.JSON(http.StatusOK, gin.H{"movie": movie})
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

	context.JSON(http.StatusOK, gin.H{"genres": genres})
}

func TempFileNamePDf() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return hex.EncodeToString(randBytes) + ".pdf"
}

func BrPDF(pdf *gopdf.GoPdf) {
	pdf.Br(20)
	pdf.SetX(20)
}

func (self *Service) GetPDF(context *gin.Context) {
	var movie models.Movie
	if err := self.DB.Where("id = ?", context.Param("id")).First(&movie).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var movieLong longmodels.Movie
	var err error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	pdf.SetXY(20, 20)
	pdf.AddTTFFont("anpro", "/usr/share/fonts/truetype/anonymous-pro/Anonymous Pro.ttf")
	pdf.SetFont("anpro", "", 14)

	pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - 60, H: gopdf.PageSizeA4.H},
		fmt.Sprintf("Название: %s/%s", movieLong.Name, movieLong.AlternativeName))
	BrPDF(&pdf)
	if movieLong.Year != 0 {
		pdf.Text(fmt.Sprintf("Год: %d    Страна: %s", movieLong.Year, movieLong.Country.Name))
	} else {
		pdf.Text(fmt.Sprintf("Год: Нет информации    Страна: %s", movieLong.Country.Name))
	}
	BrPDF(&pdf)
	pdf.Text(fmt.Sprintf("Средняя оценка: %.2f", movieLong.Score))
	BrPDF(&pdf)
	pdf.Text(fmt.Sprintf("Количество оценок: %d", movieLong.Votes))
	pdf.Image(movieLong.Posters[0].Url, gopdf.PageSizeA4.W/2-240, 130, &gopdf.Rect{W: 480, H: 600})
	pdf.SetXY(20, 780)
	if movieLong.Status.Name != "undefined" {
		pdf.Text(fmt.Sprintf("Длительность: %d    Статус: %s", movieLong.Duration, movieLong.Status.Name))
	} else {
		pdf.Text(fmt.Sprintf("Длительность: %d    Статус: %s", movieLong.Duration, "Неизвестно"))
	}
	BrPDF(&pdf)
	if movieLong.AgeRating != 0 {
		pdf.Text(fmt.Sprintf("Возрастной рейтинг: %d+", movieLong.AgeRating))
	} else {
		pdf.Text(fmt.Sprintf("Возрастной рейтинг: Неизвестно"))
	}
	pdf.AddPage()
	pdf.SetXY(20, 20)
	pdf.MultiCell(&gopdf.Rect{W: gopdf.PageSizeA4.W - 60, H: gopdf.PageSizeA4.H},
		fmt.Sprintf("Описание: %s", movieLong.Description))
	var genresNames []string
	for _, genre := range movieLong.Genres {
		genresNames = append(genresNames, genre.Name)
	}
	BrPDF(&pdf)
	pdf.Text(fmt.Sprintf("Жанры: %s", strings.Join(genresNames, ", ")))
	if len(movieLong.Fees) != 0 {
		BrPDF(&pdf)
		pdf.Text("Сборы:")
		for _, fee := range movieLong.Fees {
			BrPDF(&pdf)
			pdf.Text(fmt.Sprintf("%s : %d%s", fee.Area.Name, fee.Value, fee.Currency))
		}
	}
	if len(movieLong.Persons) != 0 {
		pdf.AddPage()
		pdf.SetXY(20, 20)
		pdf.Text("Участники:")
		BrPDF(&pdf)
		for index, person := range movieLong.Persons {
			if index != 0 && index%40 == 0 {
				pdf.AddPage()
				pdf.SetXY(20, 20)
			}
			if person.Name != "" {
				pdf.Text(fmt.Sprintf("%s/%s - %s", person.Name, person.NameEn, person.Profession.NameEn))
			} else {
				pdf.Text(fmt.Sprintf("%s - %s", person.NameEn, person.Profession.NameEn))
			}
			BrPDF(&pdf)
		}
	}

	pdfName := TempFileNamePDf()
	pdf.WritePdf("./res/pdf/" + pdfName)
	time.AfterFunc(1*time.Hour, func() { os.Remove("./res/pdf/" + pdfName) })

	context.JSON(http.StatusOK, gin.H{"pdf": self.BaseUrl + "/res/pdf/" + pdfName})
}
