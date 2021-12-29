package main

import (
	"fmt"
	"home/database"
	"home/database/models"
)

type Switch struct {
	name  string
	state bool
}

var Switches = map[int]*Switch{
	0: {name: "kitchen", state: false},
	1: {name: "kitchen_backlight", state: false},
	2: {name: "bathroom", state: false},
	3: {name: "bathroom_ventilation", state: false},
	4: {name: "bathroom_towel_dryer", state: false},
	5: {name: "livingroom_entrance", state: false},
	6: {name: "livingroom_middle", state: false},
	7: {name: "livingroom_fireplace", state: false},
	8: {name: "storageroom", state: false},
	9: {name: "stairs", state: false},
}

var RegistersPrevState = []uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func InitSwitcher() {

}

func HandleRegisters(registers []uint16) {
	for i, r := range registers {
		if RegistersPrevState[i] == 1 && r == 0 {
			sw, ok := Switches[i]
			if !ok {
				fmt.Println("unknown switch, id:", i)
			} else {
				Switches[i].state = !sw.state
				database.DB.Create(&models.ModbusSwitcher{
					Name:  Switches[i].name,
					State: Switches[i].state,
				})
				fmt.Println("changing switch state:", Switches[i].name, "to", Switches[i].state)
			}
		}
	}
	RegistersPrevState = registers
}