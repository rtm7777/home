package main

import (
	"fmt"
	"home/modbus"
)

func main() {
	modbus.InitClients()

	var total float32
	var err error

	voltage, _ := modbus.SDM230.ReadInput("voltage")
	fmt.Println(voltage, "Volts")

	freq, _ := modbus.SDM230.ReadInput("frequency")
	fmt.Println(freq, "Hz")

	if total, err = modbus.SDM230.ReadInput("total"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(total, "kwh")
	}
}
