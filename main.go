package main

import (
	"log"
	"os"

	"home/database"
	"home/modbus"
	"home/modules/boiler"
	"home/modules/counters/sdm"
	"home/modules/switcher"
	"home/server"
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

	database.Connect()
	database.Migrate()

	modbus.InitClients()

	switcher.Init()
	sdm.Init()
	boiler.Init()

	go switcher.StartPolling()
	go sdm.StartPolling()
	// go boiler.StartTempPolling()

	server.Run()

	<-exit
}
