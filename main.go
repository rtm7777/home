package main

import (
	"fmt"
	"home/database"
	"home/database/models"
	"home/modbus"
	"home/poller"
	"time"
)

func main() {
	modbus.InitClients()

	database.Connect()
	database.Migrate()

	go startSwitchesPolling()
	go startSDMPolling()

	exit := make(chan bool)
	<-exit
}

func startSwitchesPolling() {
	tick := make(chan time.Time)
	go poller.NewPoller(time.Millisecond * 200).Run(tick)

	for range tick {
		if registers, err := modbus.N4DIH32.ReadHoldingRegisters(); err != nil {
			fmt.Println(err.Error())
		} else {
			HandleSwitches(registers)
		}
	}
}

func startSDMPolling() {
	inputs := []string{"ActivePower", "Voltage", "Current", "Frequency", "TotalEnergy"}
	tick := make(chan time.Time)
	go poller.NewPoller(time.Minute * 5).Run(tick)

	for range tick {
		var values []float32
		for _, input := range inputs {
			value, err := modbus.SDM230.ReadFloatInput(input)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			values = append(values, value)
		}

		database.DB.Create(&models.Sdm230{
			ActivePower: values[0],
			Voltage:     values[1],
			Current:     values[2],
			Frequency:   values[3],
			TotalEnergy: values[4],
		})
	}
}
