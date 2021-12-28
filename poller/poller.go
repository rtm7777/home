package poller

import (
	"fmt"
	"time"

	"home/modbus"
)

type poller struct {
	ticker *time.Ticker
}

func NewPoller() *poller {
	return &poller{
		ticker: time.NewTicker(time.Millisecond * 250),
	}
}

func (p *poller) Run(device *modbus.ModbusDevice, registers chan<- []uint16) {
	for {
		select {
		case <-p.ticker.C:
			if hRegs, err := device.ReadHoldingRegisters(); err != nil {
				fmt.Println(err.Error())
			} else {
				registers <- hRegs
			}
		}
	}
}
