package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	ID          string `gorm:"type:uuid;primary_key" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//BookedDateTime string        `json:"booked_datetime"` // 使用逗号分隔的字符串
	Reservations []Reservation `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID for new room
	r.ID = uuid.New().String()
	return
}
