package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sadiq810/prayer_times_scrapper/models"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type NewTimingController struct {
}

type SalahTimingType struct {
	Fajr          string
	Shouruq       string
	Zuhr          string
	Asr           string
	Maghrib       string
	Isha          string
	Day           int
	Month         int
	IqamahFajr    string `json:"iqamah_Fajr"`
	IqamahZuhr    string `json:"iqamah_Zuhr"`
	IqamahAsr     string `json:"iqamah_Asr"`
	IqamahMaghrib string `json:"iqamah_Maghrib"`
	IqamahIsha    string `json:"iqamah_Isha"`
}

type MasjidTimingResponse struct {
	Model      Model   `json:"model"`
	HasError   bool    `json:"hasError"`
	Message    *string `json:"message"` // use pointer to handle null
	StatusCode int     `json:"statusCode"`
	Errors     *string `json:"errors"` // use pointer to handle null
	Masjid     models.Masjid
}

type Model struct {
	//MasjidDetails           interface{}       `json:"masjidDetails"`
	MasjidSettings MasjidSettings    `json:"masjidSettings"`
	SalahTimings   []SalahTimingType `json:"salahTimings"`
	/*IqamahTimings           interface{}       `json:"iqamahTimings"`
	PrimaryLanguage         interface{}       `json:"primaryLanguage"`
	SecondaryLanguage       interface{}       `json:"secondaryLanguage"`
	JumahSalahIqamahTimings []interface{}     `json:"jumahSalahIqamahTimings"`
	LastUpdatedAt           string            `json:"lastUpdatedAt"`*/
}

type MasjidSettings struct {
	IsDstOn                    bool   `json:"isDstOn"`
	JumahTime                  string `json:"jumahTime"`
	IsTimingsUploaded          bool   `json:"isTimingsUploaded"`
	ShowJumahTime              bool   `json:"showJumahTime"`
	HijriOffset                int    `json:"hijriOffset"`
	JummahTimeEqualsZuhrTime   bool   `json:"jummahTimeEqualsZuhrTime"`
	ShowTomorrowIqamahTimes    bool   `json:"showTomorrowIqamahTimes"`
	StandardIqamahCalculation  bool   `json:"standardIqamahCalculation"`
	ShowIqamahMinutesAsTime    bool   `json:"showIqamahMinutesasTime"`
	ShowHijriCalender          bool   `json:"showHijriCalender"`
	DisplayTimeIn12HourFormat  bool   `json:"displayTimeIn12HourFormat"`
	EnableMarkers              bool   `json:"enableMarkers"`
	PlaySoundBeforeIqamah      bool   `json:"playSoundBeforeIqamah"`
	EnableMultipleSalahTimings bool   `json:"enableMultipleSalahTimings"`
	EnableArc                  bool   `json:"enableArc"`
}

func (tc NewTimingController) Handle(db *gorm.DB) {
	var masjids []models.Masjid
	var prayerNames []models.PrayerName

	db.Select("id, title").Order("sort_order asc").Find(&prayerNames)

	db.Where("f_status=0").Where("f_id is not null").Order("id asc").Find(&masjids)

	timingChannel := make(chan MasjidTimingResponse)

	var wg sync.WaitGroup

	numberOfRoutines := runtime.NumCPU()

	// limit number of goroutines upto 15.
	if numberOfRoutines > 15 {
		numberOfRoutines = 15
	}

	for i := 0; i < numberOfRoutines; i++ {
		wg.Add(1)

		go saveTimings(timingChannel, prayerNames, db, &wg)
	}

	for _, masjid := range masjids {
		fmt.Printf("Fetching timings for masjid: %v \n", masjid.ID)

		var response MasjidTimingResponse

		requestStartTime := time.Now()

		url := fmt.Sprintf("https://time.my-masjid.com/api/TimingsInfoScreen/GetMasjidTimings?GuidId=%v", masjid.FGuid)

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

		fmt.Printf("Saving timings for masjid: %v & total records: %v\n", masjid.ID, len(response.Model.SalahTimings))

		response.Masjid = masjid

		timingChannel <- response

		requestEnd := time.Now()

		d := requestEnd.Sub(requestStartTime)

		tm := d.Milliseconds()

		masjid.FId = 0
		masjid.FStatus = 1
		masjid.FGuid = ""
		db.Save(&masjid)

		if tm < 400 {
			time.Sleep(time.Millisecond * 300)
		}
	} // end loop of masjid

	close(timingChannel)

	wg.Wait()
}

func saveTimings(responseChannel <-chan MasjidTimingResponse, prayerNames []models.PrayerName, db *gorm.DB, wg *sync.WaitGroup) {

	defer wg.Done()

	for response := range responseChannel {
		masjid := response.Masjid

		jummaTiming := response.Model.MasjidSettings.JumahTime

		var exist models.MasjidPrayerTime

		db.Where("masjid_id=?", masjid.ID).Where("prayer_name_id=?", 6).Where("masjid_id=?", masjid.ID).First(&exist)

		if exist.ID > 0 {
			exist.IqamahTime = jummaTiming

			db.Save(&exist)
		} else {
			jt := models.MasjidPrayerTime{MasjidId: masjid.ID, PrayerNameId: 6, IqamahTime: jummaTiming}

			db.Create(&jt)
		}

		for _, timing := range response.Model.SalahTimings {
			if timing.Month > 12 {
				continue
			}

			for _, prayerName := range prayerNames {
				var exist models.MasjidPrayerTime

				db.Where("masjid_id=?", masjid.ID).Where("prayer_name_id=?", prayerName.ID).Where("day=?", timing.Day).Where("month=?", timing.Month).First(&exist)

				if exist.ID != 0 {
					err := setFields(timing, &exist, prayerName.Title)

					if err != nil {
						continue
					}

					db.Save(&exist)
				} else {
					var prayerTime models.MasjidPrayerTime

					prayerTime.MasjidId = masjid.ID
					prayerTime.PrayerNameId = prayerName.ID
					prayerTime.Day = timing.Day
					prayerTime.Month = timing.Month

					err := setFields(timing, &prayerTime, prayerName.Title)

					if err != nil {
						continue
					} else {
						db.Create(&prayerTime)
					}
				}
			} // end prayer name loop.
		}
	} // end of channel loop
}

func setFields(timing SalahTimingType, record *models.MasjidPrayerTime, prayerName string) error {
	//fmt.Printf("%#v\n", timing)
	if prayerName == "Fajr" {
		record.AdhanTime = timing.Fajr
		record.IqamahTime = timing.IqamahFajr
	} else if prayerName == "shouruq" {
		record.AdhanTime = timing.Shouruq
		record.IqamahTime = ""
	} else if prayerName == "Dhuhr" {
		record.AdhanTime = timing.Zuhr
		record.IqamahTime = timing.IqamahZuhr
	} else if prayerName == "Asr" {
		record.AdhanTime = timing.Asr
		record.IqamahTime = timing.IqamahAsr
	} else if prayerName == "Maghrib" {
		record.AdhanTime = timing.Maghrib
		record.IqamahTime = timing.IqamahMaghrib
	} else if prayerName == "Isha" {
		record.AdhanTime = timing.Isha
		record.IqamahTime = timing.IqamahIsha
	} else {
		return errors.New("Invalid Prayer Time")
	}

	return nil
}
