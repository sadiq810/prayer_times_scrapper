package models

import (
	"time"
)

type Masjid struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CountryId uint      `json:"country_id"`
	CityId    uint      `json:"city_id"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	Address   string    `json:"address"`
	Lat       float32   `json:"lat"`
	Lng       float32   `json:"lng"`
	FId       uint      `json:"f_id"`
	FStatus   int       `json:"f_status"`
	FGuid     string    `json:"f_guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Masjid) TableName() string {
	return "masjids"
}
