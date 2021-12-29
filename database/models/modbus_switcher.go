package models

import "time"

type ModbusSwitcher struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	Name      string
	State     bool
}
