package modbus

import (
	"errors"

	"github.com/simonvetter/modbus"
)

type Registers map[string]uint16
type RegisterStates map[string]uint16

type ModbusDevice struct {
	Client                string
	HoldingRegisters      Registers
	HoldingRegistersRange []uint16
	HoldingRegisterStates RegisterStates
	InputRegisters        Registers
	UnitId                uint8
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
	HoldingRegistersRange: []uint16{128, 20},
}

var R4D1C32 = ModbusDevice{
	Client: "usb",
	UnitId: 1,
	HoldingRegisters: Registers{
		"kitchen":              1,
		"kitchen_backlight":    2,
		"bathroom":             3,
		"bathroom_ventilation": 4,
		"bathroom_towel_dryer": 5,
		"livingroom_entrance":  6,
		"livingroom_middle":    7,
		"livingroom_fireplace": 8,
		"storageroom":          9,
		"stairs":               10,
	},
	HoldingRegisterStates: RegisterStates{
		"open":   0x100,
		"close":  0x200,
		"toggle": 0x300,
	},
}

func (d ModbusDevice) ReadFloatInput(name string) (float32, error) {
	err := Clients[d.Client].SetUnitId(d.UnitId)
	if err != nil {
		return 0, err
	}
	input, ok := d.InputRegisters[name]
	if !ok {
		return 0, errors.New("Unknown input name: " + name)
	}
	return Clients[d.Client].ReadFloat32(input, modbus.INPUT_REGISTER)
}

func (d ModbusDevice) ReadHoldingRegister(register uint16) (uint16, error) {
	err := Clients[d.Client].SetUnitId(d.UnitId)
	if err != nil {
		return 0, err
	}
	return Clients[d.Client].ReadRegister(register, modbus.HOLDING_REGISTER)
}

func (d ModbusDevice) writeRegister(register uint16, value uint16) error {
	err := Clients[d.Client].SetUnitId(d.UnitId)
	if err != nil {
		return err
	}
	return Clients[d.Client].WriteRegister(register, value)
}

func (d ModbusDevice) WriteHoldingRegister(name string, value uint16) error {
	registerAddress, ok := d.HoldingRegisters[name]
	if !ok {
		return errors.New("Unknown register name: " + name)
	}
	return d.writeRegister(registerAddress, value)
}

func (d ModbusDevice) ReadHoldingRegisters() ([]uint16, error) {
	err := Clients[d.Client].SetUnitId(d.UnitId)
	if err != nil {
		return []uint16{}, err
	}
	return Clients[d.Client].ReadRegisters(d.HoldingRegistersRange[0], d.HoldingRegistersRange[1], modbus.HOLDING_REGISTER)
}
