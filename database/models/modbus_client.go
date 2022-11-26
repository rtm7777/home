package models

import (
	"time"
)

type ModbusClient struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	Name      string
	URL       string
	Speed     uint
	StopBits  uint
	Timeout   uint
}
