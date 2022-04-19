package main

import (
	"log"
	"time"

	"home/database"
	"home/database/models"
	"home/modbus"
)

var SwitchStates = map[bool]uint16{
	true:  modbus.R4D1C32.HoldingRegisterStates["open"],
	false: modbus.R4D1C32.HoldingRegisterStates["close"],
}

var Inputs []models.Input
var Loads []models.Load

var InputLoads map[uint16][]models.Load

var LoadStates map[uint16]bool

var RegistersPrevState = []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func InitSwitcher() {
	InputLoads = make(map[uint16][]models.Load)
	LoadStates = make(map[uint16]bool)

	database.DB.Model(&models.Input{}).Preload("Loads").Find(&Inputs)
	database.DB.Model(&models.Load{}).Find(&Loads)

	for _, load := range Loads {
		LoadStates[load.RegisterIndex] = false
	}
	for _, input := range Inputs {
		InputLoads[input.RegisterIndex] = input.Loads
	}
}

func HandleSwitches(registers []uint16) {
	for i, r := range registers {
		idx := uint16(i)
		if RegistersPrevState[idx] == 0 && r == 1 {
			for _, load := range InputLoads[idx] {
				LoadStates[load.RegisterIndex] = !LoadStates[load.RegisterIndex]
				log.Println("before register write:", time.Now())
				err := modbus.R4D1C32.WriteHoldingRegister(load.RegisterIndex, SwitchStates[LoadStates[load.RegisterIndex]])
				if err != nil {
					log.Println(err.Error())
				}

				log.Println("before DB write:", time.Now())
				go database.DB.Create(&models.ModbusSwitcher{
					Name:  load.Name,
					State: LoadStates[load.RegisterIndex],
				})
				log.Println("after DB write:", time.Now())
			}
		}
	}
	RegistersPrevState = registers
}
