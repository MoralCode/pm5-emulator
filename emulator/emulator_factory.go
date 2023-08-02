package emulator

import (
	"log"
	"pm5-emulator/sm"

	"github.com/bettercap/gatt"
	"github.com/bettercap/gatt/examples/option"
)

//NewEmulator factory methods initializes emulator
func NewEmulator() *Emulator {
	d, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		log.Fatalf("Failed to open config, err: %s", err)
	}
	return &Emulator{
		device:       d,
		stateMachine: sm.NewStateMachine(),
	}
}