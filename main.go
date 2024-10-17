package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sadiq810/prayer_times_scrapper/configs"
	"github.com/sadiq810/prayer_times_scrapper/controllers"
	"log"
	"time"
)

func main() {
	startTime := time.Now()
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := configs.SetupDB()

	fmt.Println("Fetch Countries List......")
	countryController := controllers.NewCountryController{}

	countryController.Handle(db)

	fmt.Println("Fetch Cities List......")
	cityController := controllers.NewCityController{}
	cityController.Handle(db)

	fmt.Println("Fetch Masjid List......")
	masjidController := controllers.NewMasjidController{}
	masjidController.Handle(db)

	fmt.Println("Fetch Masjid Timings List......")
	timingController := controllers.NewTimingController{}
	timingController.Handle(db)

	endTime := time.Now()

	fmt.Println("All Done, Total Script Execution Time: ", endTime.Sub(startTime))
}
