package modbus

import (
	"errors"

	"github.com/simonvetter/modbus"
)

type Registers map[string]uint16

type ModbusDevice struct {
	Client                string
	UnitId                uint8
	FirstHoldingRegister  uint16
	HoldingRegistersCount uint16
	InputRegisters        Registers
}

var SDM230 = ModbusDevice{
	Client: "usb",
	UnitId: 11,
	InputRegisters: Registers{
		"Voltage":       0,
		"Current":       6,
		"ActivePower":   36,
		"ApparentPower": 18,
		"ReactivePower": 24,
		"Frequency":     70,
		"TotalEnergy":   342,
	},
}

var N4DIH32 = ModbusDevice{
	Client:                "serial",
	UnitId:                1,
	FirstHoldingRegister:  128,
	HoldingRegistersCount: 20,
}

func (d ModbusDevice) ReadFloatInput(name string) (float32, error) {
	Clients[d.Client].SetUnitId(d.UnitId)
	input, ok := d.InputRegisters[name]
	if !ok {
		return 0, errors.New("Unknown input")
	}
	return Clients[d.Client].ReadFloat32(input, modbus.INPUT_REGISTER)
}

func (d ModbusDevice) ReadHoldingRegisters() ([]uint16, error) {
	Clients[d.Client].SetUnitId(d.UnitId)
	return Clients[d.Client].ReadRegisters(d.FirstHoldingRegister, d.HoldingRegistersCount, modbus.HOLDING_REGISTER)
}
