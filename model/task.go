package model

import (
	"time"
)

type Task struct {
	Id          uint64    `gorm:"primary_key;auto_increment;not null" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:varchar(1024);not null" json:"description"`
	StartTime   time.Time `gorm:"type:timestamp;not null" json:"start_time"`
	EndTime     time.Time `gorm:"type:timestamp;not null" json:"end_time"`
}
