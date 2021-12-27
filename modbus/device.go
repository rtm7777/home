package modbus

import (
	"errors"

	"github.com/simonvetter/modbus"
)

type Registers map[string]uint16

type ModbusDevice struct {
	Client string
	UnitId uint8
	Registers
}

var SDM230 = ModbusDevice{
	Client: "usb",
	UnitId: 11,
	Registers: Registers{
		"voltage":        0,
		"current":        6,
		"active_power":   12,
		"apparent_power": 18,
		"reactive_power": 24,
		"frequency":      70,
		"total_energy":   342,
	},
}

func (d ModbusDevice) ReadInput(name string) (float32, error) {
	Clients[d.Client].SetUnitId(d.UnitId)
	input, ok := d.Registers[name]
	if !ok {
		return 0, errors.New("Unknown input")
	}
	return Clients[d.Client].ReadFloat32(input, modbus.INPUT_REGISTER)
}
