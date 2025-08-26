package chip8

import (
	"fmt"
	"time"
)

type Chip8 struct {
	Memory       [4096]byte
	Display      [64][32]bool // A 64 x 32 matrix of which pixels are turned on
	PC           uint16       // Program counter
	Index        uint16       // Index register, points to memory locations
	Stack        [16]uint16
	cpuTimer     uint8 // Decrements 60 times per second until reaching 0
	delayTimer   uint8 // Gives beep as long as not 0
	timeStart    time.Time
	tickDuration time.Duration
	Registers    [16]uint8 // Variable registers, may need to change this

	stackPointer int
}

// NewChip8FromByte takes a slice of bytes and returns a Chip8 emulator with default settings
// and the ROM loaded into memory
func NewChip8FromByte(rom []byte) (Chip8, error) {
	if len(rom) == 0 {
		return Chip8{}, fmt.Errorf("no rom data provided")
	}

	c := Chip8{
		PC:           0x200,
		tickDuration: time.Second / 60,
	}

	for i, byt := range rom {
		c.Memory[0x200+i] = byt
	}
	return c, nil
}

// Update will process the next instruction. If more than a second has passed since the last tick
// it will advance the delay and sound timers. It is recommended to run this loop around 700 times
// per second for most purposes but it should be configured. This does not handle exact cycle timing.
// Note that on a very slow process such as stepping through instructions the timers will still only
// count down at most once per execution.
func (c *Chip8) Update() error {
	if time.Since(c.timeStart) > c.tickDuration {
		if c.cpuTimer > 0 {
			c.cpuTimer -= 1
		}

		if c.delayTimer > 0 {
			c.delayTimer -= 1
		}
		c.timeStart = time.Now() // start the new tick
	}
	instruction, err := c.fetch()
	if err != nil {
		return err
	}

	err = c.execute((instruction))
	if err != nil {
		return err
	}

	return nil
}

// fetch the next instruction
func (c *Chip8) fetch() (uint16, error) {
	if (int)(c.PC+2) > len(c.Memory) {
		return 0, fmt.Errorf("out of memory! program counter at: %d", c.PC)
	}

	instruction := uint16FromTwoBytes(c.Memory[c.PC], c.Memory[c.PC+1])
	c.PC = c.PC + 2 // Increment the program counter to the next instruction
	return instruction, nil
}

// process the instruction
func (c *Chip8) execute(instruction uint16) error {
	x := (uint8)(instruction & 0x0F00 >> 8)
	//	y := (uint8)(instruction & 0x00F0 >> 4)
	// n := (uint8)(instruction & 0x000F)
	nn := (uint8)(instruction & 0x00FF)
	nnn := (uint16)(instruction & 0x0FFF)

	switch instruction & 0xF000 {
	case 0x0000:
		switch instruction {
		case 0x00E0:
			c.clearDisplay()
		case 0x00EE:
			c.stackPop()
		default: // We explicity ignore any other 0x000 instruction
			return nil
		}
	case 0x1000:
		c.jump(nnn)
	case 0x2000:
		c.stackPush(c.PC)
		c.PC = nnn
	case 0x6000:
		c.setRegister((int)(x), nn)
	default:
		return fmt.Errorf("unknown instruction: %04X", instruction)
	}
	return nil
}
