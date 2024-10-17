package models

import (
	"gorm.io/gorm"
	"time"
)

type Tabler interface {
	TableName() string
}

type Country struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `json:"title"`
	Image     string         `json:"image"`
	FId       int            `json:"f_id"`
	FStatus   int            `json:"f_status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Country) TableName() string {
	return "countries"
}
