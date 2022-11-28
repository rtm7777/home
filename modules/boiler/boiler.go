package boiler

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"home/database"
	"home/database/models"
	"home/modbus"
	"home/poller"
)

type BoilerState struct {
	desiredTemp         float32
	isON                bool
	isBoiling           bool
	floorPumpActive     bool
	tempReadErrorsCount int
	tempPoller          poller.Poller
	roomTemperature     [200]int
	heatFlowTemperature [200]int
}

type TempSensorDevice struct {
	modbus.ModbusDevice
	sensors map[string]uint16
}

type RelayDevice struct {
	modbus.ModbusDevice
	loads map[string]uint16
}

var Boiler BoilerState

var TempSensor TempSensorDevice
var Relay RelayDevice

func Init() {
	Boiler.isON = false
	Boiler.desiredTemp = 19.5
	Boiler.tempPoller = poller.NewPoller(time.Second * 1)

	var tempSensorModule models.ModbusModule
	var relayModule models.ModbusModule
	var tempSensorConfig map[string]uint16
	var relayConfig map[string]uint16

	database.DB.Model(&models.ModbusModule{}).Preload("ModbusClient").Where("name = ?", "R4DCB08").First(&tempSensorModule)
	json.Unmarshal([]byte(tempSensorModule.Config), &tempSensorConfig)
	TempSensor = TempSensorDevice{
		ModbusDevice: modbus.ModbusDevice{
			Client: tempSensorModule.ModbusClient.Name,
			UnitId: tempSensorModule.UnitId,
		},
		sensors: tempSensorConfig,
	}

	database.DB.Model(&models.ModbusModule{}).Preload("ModbusClient").Where("name = ?", "N4DIG08").First(&relayModule)
	json.Unmarshal([]byte(relayModule.Config), &relayConfig)
	Relay = RelayDevice{
		ModbusDevice: modbus.ModbusDevice{
			Client: relayModule.ModbusClient.Name,
			UnitId: relayModule.UnitId,
		},
		loads: relayConfig,
	}
}

func StartTempPolling() {
	go Boiler.tempPoller.Init()
	for t := range Boiler.tempPoller.Tick {
		if Boiler.isON {
			go readTempSensor(t)
		}
	}
}

func PowerON() {
	Boiler.isON = true
}

func PowerOFF() {
	Boiler.isON = false
}

func GetState() (state struct {
	IsON              bool `json:"isON"`
	IsBoiling         bool `json:"isBoiling"`
	IsFloorPumpActive bool `json:"isFloorPumpActive"`
}) {
	state.IsON = Boiler.isON
	state.IsBoiling = Boiler.isBoiling
	state.IsFloorPumpActive = Boiler.floorPumpActive
	return
}

func readTempSensor(t time.Time) {
	fmt.Println(Boiler.tempReadErrorsCount)
	if registers, err := TempSensor.ReadHoldingRegisters(0, 8); err != nil {
		log.Println(err.Error())
		Boiler.tempReadErrorsCount++
		if Boiler.tempReadErrorsCount > 10 {
			Boiler.tempPoller.Throttle(true)
		}
	} else {
		Boiler.tempReadErrorsCount = 0
		fmt.Println(registers)
	}
}
