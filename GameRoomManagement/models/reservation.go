package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Reservation struct {
	ID             string    `gorm:"type:uuid;primary_key" json:"id"` // UUID 存储为字符串
	RoomID         string    `json:"room_id"`
	WholeHourStart time.Time `json:"whole_hour_start"` // 整点开始时间，包括日期和时间
	Player         string    `json:"player"`
}

func (r *Reservation) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New().String()
	return
}
