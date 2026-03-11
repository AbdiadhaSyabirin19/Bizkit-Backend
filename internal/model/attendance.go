package model

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	UserID   uint       `json:"user_id"`
	PhotoIn  string     `json:"photo_in"`
	PhotoOut string     `json:"photo_out"`
	CheckIn  time.Time  `json:"check_in"`
	CheckOut *time.Time `json:"check_out"` // pointer → NULL saat belum checkout
	User     User       `json:"user" gorm:"foreignKey:UserID"`
}