package main

import (
	"fmt"
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

	database.DB.Debug().Model(&models.Input{}).Preload("Loads").Find(&Inputs)
	database.DB.Debug().Model(&models.Load{}).Find(&Loads)

	for _, load := range Loads {
		LoadStates[load.RegisterIndex] = false
	}
	for _, input := range Inputs {
		InputLoads[input.RegisterIndex] = input.Loads
	}
}

func HandleSwitches(registers []uint16) {
	for i, r := range registers {
		if RegistersPrevState[i] == 1 && r == 0 {
			// find all loads this input controls
			for _, load := range InputLoads[uint16(i)] {
				LoadStates[load.RegisterIndex] = !LoadStates[load.RegisterIndex]
				database.DB.Create(&models.ModbusSwitcher{
					Name:  load.Name,
					State: LoadStates[load.RegisterIndex],
				})
				err := modbus.R4D1C32.WriteHoldingRegister(uint16(i), SwitchStates[LoadStates[uint16(i)]])
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
	RegistersPrevState = registers
}
