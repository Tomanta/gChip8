package chip8

import "fmt"

type Chip8 struct {
	memory      [4096]byte
	Display     [64][32]bool // A 64 x 32 matrix of which pixels are turned on
	pc          uint16       // Program counter
	index       uint16       // Index register, points to memory locations
	stack       []uint16
	timer       int      // Decrements 60 times per second until reaching 0
	sound_timer int8     // Gives beep as long as not 0
	variables   [16]int8 // Variable registers, may need to change this
}

func NewChip8FromFile(filepath string) (Chip8, error) {
	return Chip8{}, fmt.Errorf("not yet implemented")
}

func NewChip8FromByte(rom []byte) (Chip8, error) {
	if len(rom) == 0 {
		return Chip8{}, fmt.Errorf("no rom data provided")
	}

	c := Chip8{
		pc: 0x200,
	}
	for i, byt := range rom {
		c.memory[0x200+i] = byt
	}
	return c, nil
}
