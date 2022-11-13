package controllers

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	longmodels "github.com/DanArmor/MovieDB_backend/pkg/long_models"
	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
)

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
	var movieLong longmodels.Movie
	self.DB.Preload("Country").Preload("MovieType").
			Preload("Posters", "poster_type_id = ?", self.PreviewID).
			Preload("Genres").Preload("Status").Preload("Fees.Area").
			Where("id = ?", context.Param("id")).First(&movieLong)
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
	dptr.Find(&movieLong.Persons)

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
				pdf.Text(fmt.Sprintf("%s/%s - %s", person.Name, person.NameEn, person.ProfessionNameEn))
			} else {
				pdf.Text(fmt.Sprintf("%s - %s", person.NameEn, person.ProfessionNameEn))
			}
			BrPDF(&pdf)
		}
	}

	pdfName := TempFileNamePDf()
	pdf.WritePdf("./res/pdf/" + pdfName)
	time.AfterFunc(1*time.Hour, func() { os.Remove("./res/pdf/" + pdfName) })

	context.JSON(http.StatusOK, gin.H{"pdf": self.BaseUrl + "/res/pdf/" + pdfName})
}