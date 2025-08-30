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
	delayTimer   uint8 // Decrements 60 times per second until reaching 0
	soundTimer   uint8 // Decrements 60 times per second until reaching 0; should beep
	timeStart    time.Time
	tickDuration time.Duration
	Registers    [16]uint8 // Variable registers, may need to change this
	keysPressed  []byte    // Holds a list of all pressed keys

	stackPointer int
	DebugMsg     string
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

	// Copy the rom data into memory
	for i, byt := range rom {
		c.Memory[0x200+i] = byt
	}

	c.loadFonts()
	return c, nil
}

var fonts = [...]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func (c *Chip8) loadFonts() {
	var start_mem int = 0x50
	for i := range fonts {
		c.Memory[start_mem+i] = fonts[i]
	}
}

func (c *Chip8) SetKeysPressed(keys []byte) {
	c.keysPressed = keys
}

// Update will process the next instruction. If more than a second has passed since the last tick
// it will advance the delay and sound timers. It is recommended to run this loop around 700 times
// per second for most purposes but it should be configured. This does not handle exact cycle timing.
// Note that on a very slow process such as stepping through instructions the timers will still only
// count down at most once per execution.
func (c *Chip8) Update() error {
	if time.Since(c.timeStart) > c.tickDuration {
		if c.delayTimer > 0 {
			c.delayTimer -= 1
		}

		if c.soundTimer > 0 {
			c.soundTimer -= 1
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
	X := (uint8)(instruction & 0x0F00 >> 8)
	Y := (uint8)(instruction & 0x00F0 >> 4)
	N := (uint8)(instruction & 0x000F)
	NN := (uint8)(instruction & 0x00FF)
	NNN := (uint16)(instruction & 0x0FFF)

	switch instruction & 0xF000 {
	case 0x0000:
		switch instruction {
		case 0x00E0:
			c.op00E0()
		case 0x00EE:
			c.op00EE()
		default: // We explicity ignore any other 0x000 instruction
			return fmt.Errorf("unknown instruction: %04X", instruction)
		}
	case 0x1000:
		c.op1NNN(NNN)
	case 0x2000:
		c.op2NNN(c.PC)
		c.PC = NNN
	case 0x3000:
		c.op3XNN(X, NN)
	case 0x4000:
		c.op4XNN(X, NN)
	case 0x5000:
		c.op5XY0(X, Y)
	case 0x6000:
		c.op6XNN(X, NN)
	case 0x7000:
		c.op7XNN(X, NN)
	case 0x8000:
		switch N {
		case 0:
			c.op8XY0(X, Y)
		case 1:
			c.op8XY1(X, Y)
		case 2:
			c.op8XY2(X, Y)
		case 3:
			c.op8XY3(X, Y)
		case 4:
			c.op8XY4(X, Y)
		case 5:
			c.op8XY5(X, Y)
		case 6:
			c.op8XY6(X, Y)
		case 7:
			c.op8XY7(X, Y)
		case 0xE:
			c.op8XYE(X, Y)
		default:
			return fmt.Errorf("unknown instruction: %04X", instruction)
		}
	case 0x9000:
		c.op9XY0(X, Y)
	case 0xA000:
		c.opANNN(NNN)
	case 0xB000:
		c.opBNNN(NNN)
	case 0xC000:
		c.opCXNN(X, NN)
	case 0xD000:
		c.opDXYN(X, Y, N)
	case 0xE000:
		switch NN {
		case 0x9E:
			c.opEX9E(X)
		case 0xA1:
			c.opEXA1(X)
		default:
			return fmt.Errorf("unknown instruction: %04X", instruction)
		}
	case 0xF000:
		switch NN {
		case 0x07:
			c.opFX07(X)
		case 0x15:
			c.opFX15(X)
		case 0x18:
			c.opFX18(X)
		case 0x1E:
			c.opFX1E(X)
		case 0x0A:
			c.opFX0A(X)
		case 0x29:
			c.opFX29(X)
		case 0x33:
			c.opFX33(X)
		default:
			return fmt.Errorf("unknown instruction: %04X", instruction)
		}
	default:
		return fmt.Errorf("unknown instruction: %04X", instruction)
	}
	return nil
}
