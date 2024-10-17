package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/sadiq810/prayer_times_scrapper/models"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

type NewMasjidController struct {
}

type MasjidType struct {
	ID      uint
	Name    string
	GuidId  string
	Address string
}

type MasjidResponse struct {
	Model      []MasjidType
	HasError   bool
	Message    string
	StatusCode int
	Errors     *error
}

func (mc NewMasjidController) Handle(db *gorm.DB) {
	var countries []models.Country

	db.Model(&models.Country{}).Where("f_status=0").Where("f_id is not null").Find(&countries)

	totalMasjid := 0

	for _, country := range countries {
		var cities []models.City

		db.Where("f_status=0").Where("f_id is not null").Where("country_id = ?", country.ID).Find(&cities)

		for _, city := range cities {
			fmt.Printf("Fetching masjids for country: %v and City: %v \n", country.ID, city.ID)

			var response MasjidResponse

			url := fmt.Sprintf("https://time.my-masjid.com/api/Masjid/SearchMasjidByLocation?CountryId=%v&CityId=%v", country.FId, city.FId)

			resp, err := http.Get(url)

			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				log.Fatal(err)
			}

			bodyString := string(body)

			err = json.Unmarshal([]byte(bodyString), &response)

			if err != nil {
				log.Fatalf("Error unmarshalling JSON: %v", err)
			}

			for _, masjid := range response.Model {
				exist := models.Masjid{}

				db.Where("title = ?", masjid.Name).Where("country_id = ?", country.ID).Where("city_id = ?", city.ID).First(&exist)

				if exist.ID > 0 {
					exist.FStatus = 0
					exist.FId = masjid.ID
					exist.FGuid = masjid.GuidId
					exist.Address = masjid.Address
					exist.Image = "masjid-default.png"

					db.Save(&exist)
				} else {
					msjd := models.Masjid{CountryId: country.ID, CityId: city.ID, Title: masjid.Name, FId: masjid.ID, FGuid: masjid.GuidId, Address: masjid.Address, Image: "masjid-default.png", FStatus: 0}

					db.Create(&msjd)

					totalMasjid++
				}
			} // end of masjid loop

			city.FId = 0
			city.FStatus = 1

			db.Save(&city)

			// APIs restrict more than 4 calls in a second.
			time.Sleep(time.Millisecond * 400)
		} // end of city loop

		country.FId = 0
		country.FStatus = 1

		db.Save(&country)
	} // end of country loop

	fmt.Println("Total Masjid created: ", totalMasjid)
}
