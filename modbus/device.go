package modbus

import (
	"github.com/simonvetter/modbus"
)

type Registers map[string]uint16
type RegisterStates map[string]uint16

type ModbusDevice struct {
	Client string
	UnitId uint8
}

func (d ModbusDevice) ReadFloatInput(addr uint16) (float32, error) {
	err := Clients[d.Client].SetUnitId(d.UnitId)
	if err != nil {
		return 0, err
	}
	return Clients[d.Client].ReadFloat32(addr, modbus.INPUT_REGISTER)
}

func (d ModbusDevice) ReadHoldingRegister(register uint16) (uint16, error) {
	err := Clients[d.Client].SetUnitId(d.UnitId)
	if err != nil {
		return 0, err
	}
	return Clients[d.Client].ReadRegister(register, modbus.HOLDING_REGISTER)
}

func (d ModbusDevice) WriteHoldingRegister(addr uint16, value uint16) error {
	err := Clients[d.Client].SetUnitId(d.UnitId)
	if err != nil {
		return err
	}
	return Clients[d.Client].WriteRegister(addr, value)
}

func (d ModbusDevice) ReadHoldingRegisters(addr uint16, quantity uint16) ([]uint16, error) {
	err := Clients[d.Client].SetUnitId(d.UnitId)
	if err != nil {
		return []uint16{}, err
	}
	return Clients[d.Client].ReadRegisters(addr, quantity, modbus.HOLDING_REGISTER)
}
