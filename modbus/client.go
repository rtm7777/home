package modbus

import (
	"fmt"
	"time"

	"github.com/simonvetter/modbus"
)

var Clients map[string]*modbus.ModbusClient

var clientConfigurations = map[string]modbus.ClientConfiguration{
	"usb": {
		URL:      "rtu:///dev/ttyUSB0",
		Speed:    9600,
		StopBits: 1,
		Timeout:  300 * time.Millisecond,
	},
	"serial": {
		URL:      "rtu:///dev/ttyUSB0",
		Speed:    19200,
		StopBits: 1,
		Timeout:  300 * time.Millisecond,
	},
}

func InitClients() {
	Clients = make(map[string]*modbus.ModbusClient)
	for name, configuration := range clientConfigurations {
		client, err := modbus.NewClient(&configuration)
		if err != nil {
			fmt.Println("error client creation")
		}

		err = client.Open()
		if err != nil {
			fmt.Println("error connection to:", name)
		}

		Clients[name] = client
	}
}
