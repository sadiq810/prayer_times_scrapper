package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/sadiq810/prayer_times_scrapper/models"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

type NewCountryController struct {
}

type CountryType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CountryResponse struct {
	Model      []CountryType `json:"model"`
	HasError   bool          `json:"hasError"`
	Message    *string       `json:"message"`
	StatusCode int           `json:"statusCode"`
	Errors     *string       `json:"errors"`
}

func (cc NewCountryController) Handle(db *gorm.DB) {
	resp, err := http.Get("https://time.my-masjid.com/api/Country/GetAllCountries")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)

	var response CountryResponse

	err = json.Unmarshal([]byte(sb), &response)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	for _, c := range response.Model {
		exist := models.Country{}

		db.Where("title = ?", c.Name).First(&exist)

		if exist.ID > 0 {
			exist.FStatus = 0
			exist.FId = c.ID

			db.Save(&exist)
		} else {
			country := models.Country{Title: c.Name, FId: c.ID, FStatus: 0}

			db.Create(&country)
		}
	}

	fmt.Printf("Countries inserted: %+v\n", len(response.Model))
}
