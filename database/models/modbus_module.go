package models

import (
	"time"

	"gorm.io/datatypes"
)

type ModbusModule struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	ModbusClientId uint
	ModbusClient   ModbusClient
	Name           string
	Label          string
	UnitId         uint8
	Config         datatypes.JSON
}
