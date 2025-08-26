package chip8

import "fmt"

// clear the screen
func (c *Chip8) clearDisplay() {
	var blankDisplay [64][32]bool
	c.Display = blankDisplay
}

func (c *Chip8) jump(location uint16) {
	c.PC = location
}

// op2NNN adds NNN to the stack
func (c *Chip8) op2NNN(address uint16) {
	if c.stackPointer == len(c.Stack) {
		panic("stack overflow!")
	}

	c.Stack[c.stackPointer] = address
	c.stackPointer += 1
}

func (c *Chip8) stackPop() {
	if c.stackPointer == 0 {
		return // Do nothing if stack is empty. Not sure if this is correct behavior.
	}
	c.stackPointer -= 1
	c.PC = c.Stack[c.stackPointer]
	c.Stack[c.stackPointer] = 0
}

// op6XNN sets register X to NN
func (c *Chip8) op6XNN(register int, value uint8) error {
	if register < 0 || register > 15 {
		return fmt.Errorf("invalid register: %d", register)
	}

	c.Registers[register] = value
	return nil
}

// op7XNN adds NN to register X. It does not set the overflow flag.
func (c *Chip8) op7XNN(register int, value uint8) error {
	if register < 0 || register > 15 {
		return fmt.Errorf("invalid register: %d", register)
	}

	c.Registers[register] = c.Registers[register] + value
	return nil
}

// opANNN sets the Index register to value
func (c *Chip8) opANNN(value uint16) {
	c.Index = value
}
