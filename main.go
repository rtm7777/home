package main

import (
	"fmt"
	"home/modbus"
	"home/poller"
)

func main() {
	modbus.InitClients()

	power, _ := modbus.SDM230.ReadFloatInput("power")
	fmt.Println(power, "W")

	voltage, _ := modbus.SDM230.ReadFloatInput("voltage")
	fmt.Println(voltage, "V")

	freq, _ := modbus.SDM230.ReadFloatInput("frequency")
	fmt.Println(freq, "Hz")

	if total, err := modbus.SDM230.ReadFloatInput("total_energy"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(total, "kwh")
	}

	registers := make(chan []uint16)
	go poller.NewPoller().Run(&modbus.N4DIH32, registers)
	for range registers {
		fmt.Println(<-registers)
	}
}
