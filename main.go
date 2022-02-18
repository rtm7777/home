package main

import (
	"log"
	"os"
	"time"

	"home/database"
	"home/database/models"
	"home/modbus"
	"home/poller"
)

func main() {
	exit := make(chan bool)

	LOG_FILE := "/mnt/NFS/home_app.log"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	modbus.InitClients()
	if len(modbus.Clients) != 2 {
		exit <- true
	}

	database.Connect()
	database.Migrate()

	InitSwitcher()

	go startSwitchesPolling()
	go startSDMPolling()

	<-exit
}

func startSwitchesPolling() {
	tick := make(chan time.Time)
	go poller.NewPoller(time.Millisecond * 200).Run(tick)

	for range tick {
		if registers, err := modbus.N4DIH32.ReadHoldingRegisters(); err != nil {
			log.Println(err.Error())
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
				log.Println(err.Error())
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
