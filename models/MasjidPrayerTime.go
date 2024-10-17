package models

import (
	"time"
)

type MasjidPrayerTime struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	MasjidId     uint      `json:"masjid_id"`
	PrayerNameId uint      `json:"prayer_name_id"`
	AdhanTime    string    `json:"adhan_time"`
	IqamahTime   string    `json:"iqamah_time"`
	Month        int       `json:"month"`
	Day          int       `json:"day"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (MasjidPrayerTime) TableName() string {
	return "masjid_prayer_times"
}
