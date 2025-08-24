package chip8

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

const (
	StartProgramCounter uint16 = 0x200
	StartRomMemory      uint16 = 0x200
)

func openTestRom() []byte {
	data, err := os.ReadFile("../roms/ibm_logo.ch8")
	if err != nil {
		panic(err)
	}
	return data
}

func TestCanCreateFromBytes(t *testing.T) {
	romData := openTestRom()
	romDataLength := uint16(len(romData))
	fmt.Printf("want length: %d\n\n", len(romData))
	got, err := NewChip8FromByte(romData)

	if err != nil {
		t.Errorf("Received error creating file: %q", err)
	}

	gotr := got.memory[StartRomMemory : StartRomMemory+romDataLength]
	fmt.Printf("gotr length: %d", len(gotr))

	if !reflect.DeepEqual(gotr, romData) {
		t.Errorf("ROM not loaded to memory. want %X, got %X", romData, gotr)
	}

	if got.pc != StartProgramCounter {
		t.Errorf("Program counter not set: want %X, got %X", StartProgramCounter, got.pc)
	}

}
