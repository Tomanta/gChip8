package chip8

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

// op3XNN skips one instruction if register X is equal to NN (adds 2 to Program Counter)
func (c *Chip8) op3XNN(x uint8, nn uint8) {
	if c.Registers[x] == nn {
		c.PC = c.PC + 2
	}
}

// op4XNN skips one instruction if register X is not equal to NN (adds 2 to Program Counter)
func (c *Chip8) op4XNN(x uint8, nn uint8) {
	if c.Registers[x] != nn {
		c.PC = c.PC + 2
	}
}

// op5XY0 skips one instruction if register X is equal to register Y (adds 2 to Program Counter)
func (c *Chip8) op5XY0(x uint8, y uint8) {
	if c.Registers[x] == c.Registers[y] {
		c.PC = c.PC + 2
	}
}

// op6XNN sets register X to NN
func (c *Chip8) op6XNN(register uint8, value uint8) {
	c.Registers[register] = value
}

// op7XNN adds NN to register X. It does not set the overflow flag.
func (c *Chip8) op7XNN(register uint8, value uint8) {
	c.Registers[register] = c.Registers[register] + value
}

// op8XY0 sets VX to value of VY
func (c *Chip8) op8XY0(x uint8, y uint8) {
	c.Registers[x] = c.Registers[y]
}

// op8XY1 sets VX to BITWISE OR of VX and VY
func (c *Chip8) op8XY1(x uint8, y uint8) {
	c.Registers[x] = c.Registers[x] | c.Registers[y]
}

// op8XY2 sets VX to BITWISE AND of VX and VY
func (c *Chip8) op8XY2(x uint8, y uint8) {
	c.Registers[x] = c.Registers[x] & c.Registers[y]
}

// op8XY3 sets VX to XOR of VX and VY
func (c *Chip8) op8XY3(x uint8, y uint8) {
	c.Registers[x] = c.Registers[x] ^ c.Registers[y]
}

// op8XY4 sets VX to VX plus VY. Will set carry flag.
func (c *Chip8) op8XY4(x uint8, y uint8) {
	r_x := c.Registers[x]
	r_y := c.Registers[y]
	c.Registers[x] = r_x + r_y
	if (uint16)(r_x)+(uint16)(r_y) > 0xff {
		c.Registers[0xF] = 1
	} else {
		c.Registers[0xF] = 0
	}
}

// op8XY5 sets VX to VX - VY. This does not set the carry flag.
func (c *Chip8) op8XY5(x uint8, y uint8) {
	r_x := c.Registers[x]
	r_y := c.Registers[y]
	c.Registers[x] = r_x - r_y
}

// TODO: 8XY6

// op8XY7 sets VX to VY - VX. If X is larger than Y, VF is set to 1.
// Unsure of behavior if both are = but assume it is NOT set since it does not underflow.
func (c *Chip8) op8XY7(x uint8, y uint8) {
	r_x := c.Registers[x]
	r_y := c.Registers[y]
	c.Registers[x] = r_x - r_y
	if (uint16)(r_x) >= (uint16)(r_y) {
		c.Registers[0xF] = 0
	} else {
		c.Registers[0xF] = 1
	}
}

// TODO: 8XYE

// op9XY0 skips one instruction if register X is not equal to register Y (adds 2 to Program Counter)
func (c *Chip8) op9XY0(x uint8, y uint8) {
	if c.Registers[x] != c.Registers[y] {
		c.PC = c.PC + 2
	}
}

// opANNN sets the Index register to value
func (c *Chip8) opANNN(value uint16) {
	c.Index = value
}

// TODO: BNNN

// TODO: CXNN

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

// TODO: EX9E

// TODO: EX9A1

// TODO: FX01

// TODO: FX15

// TODO: FX18

// TODO: FX1E

// TODO: FX0A

// TODO: FX29

// TODO: FX33

// TODO: FX55

// TODO: FX65
