package models

import "time"

type Sdm230 struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	ActivePower float32
	Voltage     float32
	Current     float32
	Frequency   float32
	TotalEnergy float32
}
