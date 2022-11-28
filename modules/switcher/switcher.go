package switcher

import (
	"encoding/json"
	"log"
	"time"

	"home/database"
	"home/database/models"
	"home/modbus"
	"home/poller"
)

type InputDevice struct {
	modbus.ModbusDevice
	holdingRegistersRange map[string]uint16
}

type RelayDevice struct {
	modbus.ModbusDevice
	holdingRegisterStates map[string]uint16
}

var SwitchStates = map[bool]string{
	true:  "open",
	false: "close",
}

var Input InputDevice
var Relay RelayDevice

var InputLoads map[uint16][]models.Load
var LoadStates map[uint16]bool
var RegistersPrevState = []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func Init() {
	InputLoads = make(map[uint16][]models.Load)
	LoadStates = make(map[uint16]bool)

	var inputModule models.ModbusModule
	var relayModule models.ModbusModule
	var inputConfig map[string]uint16
	var relayConfig map[string]uint16
	var inputs []models.Input
	var loads []models.Load

	database.DB.Model(&models.ModbusModule{}).Preload("ModbusClient").Where("name = ?", "N4DIH32").First(&inputModule)
	json.Unmarshal([]byte(inputModule.Config), &inputConfig)
	Input = InputDevice{
		ModbusDevice: modbus.ModbusDevice{
			Client: inputModule.ModbusClient.Name,
			UnitId: inputModule.UnitId,
		},
		holdingRegistersRange: inputConfig,
	}

	database.DB.Model(&models.ModbusModule{}).Preload("ModbusClient").Where("name = ?", "R4D1C32").First(&relayModule)
	json.Unmarshal([]byte(relayModule.Config), &relayConfig)
	Relay = RelayDevice{
		ModbusDevice: modbus.ModbusDevice{
			Client: relayModule.ModbusClient.Name,
			UnitId: relayModule.UnitId,
		},
		holdingRegisterStates: relayConfig,
	}

	database.DB.Model(&models.Input{}).Preload("Loads").Find(&inputs)
	for _, input := range inputs {
		InputLoads[input.RegisterIndex] = input.Loads
	}

	database.DB.Model(&models.Load{}).Find(&loads)
	for _, load := range loads {
		LoadStates[load.RegisterIndex] = false
	}
}

func StartPolling() {
	switchesPoller := poller.NewPoller(time.Millisecond * 200)
	go switchesPoller.Init()

	for t := range switchesPoller.Tick {
		go handleSwitches(t)
	}
}

func handleSwitches(t time.Time) {
	if registers, err := Input.ReadHoldingRegisters(Input.holdingRegistersRange["addr"], Input.holdingRegistersRange["quantity"]); err != nil {
		log.Println(err.Error())
	} else {
		for i, r := range registers {
			idx := uint16(i)
			if RegistersPrevState[idx] == 0 && r == 1 {
				for _, load := range InputLoads[idx] {
					LoadStates[load.RegisterIndex] = !LoadStates[load.RegisterIndex]
					log.Println("ticker time:", t)
					log.Println("before register write:", time.Now())
					Relay.WriteHoldingRegister(load.RegisterIndex, Relay.holdingRegisterStates[SwitchStates[LoadStates[load.RegisterIndex]]])

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
}
