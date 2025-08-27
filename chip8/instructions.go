package chip8

import "fmt"

// op00E0 clears the screen
func (c *Chip8) op00E0() {
	var blankDisplay [64][32]bool
	c.Display = blankDisplay
}

// op00EE sets the stack pointer to the top value on the stack (pops)
func (c *Chip8) op00EE() {
	if c.stackPointer == 0 {
		return // Do nothing if stack is empty. Not sure if this is correct behavior.
	}
	c.stackPointer -= 1
	c.PC = c.Stack[c.stackPointer]
	c.Stack[c.stackPointer] = 0
}

// op1NNN jumps to memory location NNN
func (c *Chip8) op1NNN(location uint16) {
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

// opDXYN draws an N pixel tall sprite from the value at Index
// drawing is done at coordinates XY. If any pixels are turned off
// VF is set to 1.
func (c *Chip8) opDXYN(x_register uint8, y_register uint8, N uint8) {
	x := c.Registers[x_register] % 64 // Reset this each iteration
	y := c.Registers[y_register] % 32
	c.Registers[0xF] = 0

	for i := range N {
		sprite := c.Memory[c.Index+(uint16)(i)] // Get sprite data for this row

		// For each of 8 bits
		for s := range 8 {
			var x_pos uint8 = x + (uint8)(s)
			if x_pos >= 64 {
				continue
			}
			var bit uint8 = 0x80 & sprite // Get most significant bit
			sprite = sprite << 1          // Shift 1 left

			// Any bit that is on will flip the current pixel. Anything turned off sets
			// register F to 1
			if c.Display[x_pos][y] && bit == 128 {
				c.Registers[0xF] = 1
				c.Display[x_pos][y] = false
			} else if !c.Display[x_pos][y] && bit == 128 {
				c.Display[x_pos][y] = true
			}
		}
		y += 1
		// Stop drawing if reached bottom of screen
		if y >= 32 {
			break
		}
	}

}
