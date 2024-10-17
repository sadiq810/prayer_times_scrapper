package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/sadiq810/prayer_times_scrapper/models"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type NewCityController struct {
}

type CityType struct {
	ID   uint
	Name string
}

type CityType2 struct {
	ID    uint
	Title string
	Cid   uint
}

type CityResponse struct {
	Model      []CityType
	HasError   bool
	Message    string
	StatusCode int
	Errors     *error
}

var totalCities int = 0

var lock sync.Mutex

func (cc NewCityController) Handle(db *gorm.DB) {
	startTime := time.Now()

	var wg sync.WaitGroup

	var countries []models.Country

	db.Model(&models.Country{}).Where("f_status=0").Where("f_id is not null").Find(&countries)

	countryChannel := make(chan models.Country)
	resultChannel := make(chan []CityType2)

	wg.Add(1)

	go fetchCities(countryChannel, &wg, resultChannel)

	numberOfGoroutine := 1

	for i := 0; i < numberOfGoroutine; i++ {
		wg.Add(1)
		go saveCities(db, &wg, resultChannel)
	}

	for _, c := range countries {
		countryChannel <- c
	}

	close(countryChannel)

	wg.Wait()

	endTime := time.Now()

	duration := endTime.Sub(startTime)

	fmt.Printf("Total Cities Inserted: %v with duration: %v\n", totalCities, duration)
}

func saveCities(db *gorm.DB, wg *sync.WaitGroup, results chan []CityType2) {
	defer wg.Done()

	for cities := range results {

		for _, ct := range cities {
			exist := models.City{}

			db.Where("title = ?", ct.Title).Where("country_id = ?", ct.Cid).First(&exist)

			if exist.ID > 0 {
				exist.FStatus = 0
				exist.FId = ct.ID

				db.Save(&exist)
			} else {
				city := models.City{CountryId: ct.Cid, Title: ct.Title, FId: ct.ID, FStatus: 0}

				db.Create(&city)

				lock.Lock()
				totalCities++
				lock.Unlock()
			}
		}
	}
}

func fetchCities(countries chan models.Country, wg *sync.WaitGroup, results chan<- []CityType2) {
	defer wg.Done()
	defer close(results)

	for c := range countries {

		fmt.Printf("Fetching Cities for Country ID: %v \n", c.ID)

		var response CityResponse

		url := fmt.Sprintf("https://time.my-masjid.com/api/City/GetCitiesByCountryId?CountryId=%v", c.FId)

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

		r := []CityType2{}

		for _, ct := range response.Model {
			r = append(r, CityType2{
				Cid:   c.ID,
				ID:    ct.ID,
				Title: ct.Name,
			})
		}

		results <- r

		// APIs restrict more than 4 calls in a second.
		time.Sleep(time.Millisecond * 400)
	}
}
