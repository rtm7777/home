package models

import "gorm.io/gorm"

type Sdm230 struct {
	gorm.Model
	ActivePower float32
	Voltage     float32
	Current     float32
	Frequency   float32
	TotalEnergy float32
}
