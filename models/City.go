package models

import (
	"gorm.io/gorm"
	"time"
)

type City struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CountryId uint           `json:"country_id"`
	Title     string         `json:"title"`
	Image     string         `json:"image"`
	FId       uint           `json:"f_id"`
	FStatus   int            `json:"f_status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (City) TableName() string {
	return "cities"
}
