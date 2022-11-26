package sdm

import (
	"encoding/json"
	"home/database"
	"home/database/models"
	"home/modbus"
	"home/poller"
	"log"
	"time"
)

type SDMDevice struct {
	modbus.ModbusDevice
	inputRegisters map[string]uint16
}

var inputs = [...]string{"ActivePower", "Voltage", "Current", "Frequency", "PhaseAngle", "TotalEnergy"}

var SDMDevices []SDMDevice

func Init() {
	var SDMModbusModules []models.ModbusModule

	database.DB.Model(&models.ModbusModule{}).Preload("ModbusClient").Where("name IN ?", []string{"SDM230"}).Find(&SDMModbusModules)

	for _, module := range SDMModbusModules {
		var config map[string]uint16
		json.Unmarshal([]byte(module.Config), &config)
		SDMDevices = append(SDMDevices, SDMDevice{
			ModbusDevice: modbus.ModbusDevice{
				Client: module.ModbusClient.Name,
				UnitId: module.UnitId,
			},
			inputRegisters: config,
		})
	}
}

func StartPolling() {
	tick := make(chan time.Time)
	go poller.NewPoller(time.Minute * 5).Run(tick)

	for t := range tick {
		go readSDMDevices(t)
	}
}

func readSDMDevices(t time.Time) {
	for _, device := range SDMDevices {
		var values []float32
		for _, input := range inputs {
			value, err := device.ReadFloatInput(device.inputRegisters[input])
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
			PhaseAngle:  values[4],
			TotalEnergy: values[5],
		})
	}
}
