package modbus

import (
	"log"
	"time"

	"home/database"
	"home/database/models"

	"github.com/simonvetter/modbus"
)

var Clients map[string]*modbus.ModbusClient

var clientsConfigs []models.ModbusClient

func InitClients() {
	database.DB.Model(&models.ModbusClient{}).Find(&clientsConfigs)

	Clients = make(map[string]*modbus.ModbusClient)

	for _, configuration := range clientsConfigs {
		client, err := modbus.NewClient(&modbus.ClientConfiguration{
			URL:      configuration.URL,
			Speed:    configuration.Speed,
			StopBits: configuration.StopBits,
			Timeout:  time.Duration(configuration.Timeout * uint(time.Millisecond)),
		})
		if err != nil {
			log.Println("error client creation")
		}

		err = client.Open()
		if err != nil {
			log.Println("error connection to:", configuration.Name)
		}

		if err == nil {
			log.Println("client init success:", configuration.Name)
			Clients[configuration.Name] = client
		}
	}
}
