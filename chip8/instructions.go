package chip8

import (
	"fmt"
	"math/rand"
	"slices"
)

// op00E0 clears the screen
func (c *Chip8) op00E0() {
	var blankDisplay [64][32]bool
	c.Display = blankDisplay
	c.DebugMsg = "Op00E0: clear screen"
}

// op00EE sets the stack pointer to the top value on the stack (pops)
func (c *Chip8) op00EE() {
	if c.stackPointer == 0 {
		return // Do nothing if stack is empty. Not sure if this is correct behavior.
	}
	c.stackPointer -= 1
	c.PC = c.Stack[c.stackPointer]
	c.Stack[c.stackPointer] = 0
	c.DebugMsg = fmt.Sprintf("Op00EE: set PC to top of stack (0x%04X)", c.PC)
}

// op1NNN jumps to memory location NNN
func (c *Chip8) op1NNN(location uint16) {
	c.PC = location
	c.DebugMsg = fmt.Sprintf("Op1NNN: set PC to NNN (0x%04X)", c.PC)
}

// op2NNN adds NNN to the stack
func (c *Chip8) op2NNN(address uint16) {
	if c.stackPointer == len(c.Stack) {
		panic("stack overflow!")
	}

	c.Stack[c.stackPointer] = address
	c.stackPointer += 1
	c.DebugMsg = fmt.Sprintf("Op2NNN: push NNN (0x%04X) to stack", address)
}

// op3XNN skips one instruction if register X is equal to NN (adds 2 to Program Counter)
func (c *Chip8) op3XNN(x uint8, nn uint8) {
	if c.Registers[x] == nn {
		c.PC = c.PC + 2
	}
	c.DebugMsg = "TODO"
}

// op4XNN skips one instruction if register X is not equal to NN (adds 2 to Program Counter)
func (c *Chip8) op4XNN(x uint8, nn uint8) {
	if c.Registers[x] != nn {
		c.PC = c.PC + 2
	}
	c.DebugMsg = "TODO"
}

// op5XY0 skips one instruction if register X is equal to register Y (adds 2 to Program Counter)
func (c *Chip8) op5XY0(x uint8, y uint8) {
	if c.Registers[x] == c.Registers[y] {
		c.PC = c.PC + 2
	}
	c.DebugMsg = "TODO"
}

// op6XNN sets register X to NN
func (c *Chip8) op6XNN(register uint8, value uint8) {
	c.Registers[register] = value
	c.DebugMsg = "TODO"
}

// op7XNN adds NN to register X. It does not set the overflow flag.
func (c *Chip8) op7XNN(register uint8, value uint8) {
	c.Registers[register] = c.Registers[register] + value
	c.DebugMsg = "TODO"
}

// op8XY0 sets VX to value of VY
func (c *Chip8) op8XY0(x uint8, y uint8) {
	c.Registers[x] = c.Registers[y]
	c.DebugMsg = "TODO"
}

// op8XY1 sets VX to BITWISE OR of VX and VY
func (c *Chip8) op8XY1(x uint8, y uint8) {
	c.Registers[x] = c.Registers[x] | c.Registers[y]
	c.DebugMsg = "TODO"
}

// op8XY2 sets VX to BITWISE AND of VX and VY
func (c *Chip8) op8XY2(x uint8, y uint8) {
	c.Registers[x] = c.Registers[x] & c.Registers[y]
	c.DebugMsg = "TODO"
}

// op8XY3 sets VX to XOR of VX and VY
func (c *Chip8) op8XY3(x uint8, y uint8) {
	c.Registers[x] = c.Registers[x] ^ c.Registers[y]
	c.DebugMsg = "TODO"
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
	c.DebugMsg = "TODO"
}

// op8XY5 sets VX to VX - VY. This does not set the carry flag.
func (c *Chip8) op8XY5(x uint8, y uint8) {
	r_x := c.Registers[x]
	r_y := c.Registers[y]
	c.Registers[x] = r_x - r_y
	c.DebugMsg = "TODO"
}

// op08XY6 shifts VY one bit to the right and stores in VX. VF is set to the bit that
// shifted out.
// NOTE: Super-CHIP 8 has different behavior that will need to be implemented; it shifts
// VX in place and ignores Y.
func (c *Chip8) op8XY6(x uint8, y uint8) {
	r_x := c.Registers[y]
	r_f := 0x01 & r_x
	c.Registers[0xF] = r_f
	c.Registers[x] = r_x >> 1
	c.DebugMsg = "TODO"
}

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
	c.DebugMsg = "TODO"
}

// op08XYE shifts VY one bit to the left and stores in VX. VF is set to the bit that
// shifted out.
// NOTE: Super-CHIP 8 has different behavior that will need to be implemented; it shifts
// VX in place and ignores Y.
func (c *Chip8) op8XYE(x uint8, y uint8) {
	r_x := c.Registers[y]
	r_f := r_x >> 7 & 0x1
	c.Registers[0xF] = r_f
	c.Registers[x] = r_x << 1
	c.DebugMsg = "TODO"
}

// op9XY0 skips one instruction if register X is not equal to register Y (adds 2 to Program Counter)
func (c *Chip8) op9XY0(x uint8, y uint8) {
	if c.Registers[x] != c.Registers[y] {
		c.PC = c.PC + 2
	}
	c.DebugMsg = "TODO"
}

// opANNN sets the Index register to value
func (c *Chip8) opANNN(value uint16) {
	c.Index = value
	c.DebugMsg = "TODO"
}

// opBNNN sets the program counter to NNN plus value in V0
// Note: Super Chip8 sets the program counter to VX plus NN instead (unclear on exact behavior)
func (c *Chip8) opBNNN(value uint16) {
	r_0 := c.Registers[0]
	c.PC = value + uint16(r_0)
	c.DebugMsg = "TODO"
}

// opCXNN generates a random number, ands it with NN, and stores in X
func (c *Chip8) opCXNN(x uint8, value uint8) {
	var r uint8 = (uint8)(rand.Intn(0xFF)) & value
	c.Registers[x] = r
	c.DebugMsg = "TODO"
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
	c.DebugMsg = "TODO"
}

// opEX9E skips one instruction if key stored in X is pressed
func (c *Chip8) opEX9E(x uint8) {
	if slices.Contains(c.keysPressed, c.Registers[x]) {
		c.PC = c.PC + 2
	}
	c.DebugMsg = "TODO"
}

// opEXA1 skips one instruction if key stored in X is not pressed
func (c *Chip8) opEXA1(x uint8) {
	if slices.Contains(c.keysPressed, c.Registers[x]) {
		c.PC = c.PC + 2
	}
	c.DebugMsg = "TODO"
}

// opFX07 sets VX to the current value of the delay timer
func (c *Chip8) opFX07(x uint8) {
	c.Registers[x] = c.delayTimer
	c.DebugMsg = fmt.Sprintf("OpFX07: set V%X to delay timer: %d", x, c.delayTimer)
}

// opFX15 sets the delay timer to the value of X
func (c *Chip8) opFX15(x uint8) {
	c.delayTimer = c.Registers[x]
	c.DebugMsg = fmt.Sprintf("OpFX15: set delayTimer to value of V%X: %d", x, c.delayTimer)
}

// opFX18 sets the sound timer to the value of X
func (c *Chip8) opFX18(x uint8) {
	c.soundTimer = c.Registers[x]
	c.DebugMsg = fmt.Sprintf("OpFX18: set soundTimer to value of V%X: %d", x, c.soundTimer)
}

// opFX1E adds the value of X to the index register. If it overflows from
// 0FFF to 1000 it should set VF to 1, this is not standard behavior but is safe
func (c *Chip8) opFX1E(x uint8) {
	new_i := c.Index + (uint16)(c.Registers[x])
	if new_i >= 0x1000 {
		c.Registers[0xF] = 1
		new_i -= 0x1000
	}
	c.DebugMsg = fmt.Sprintf("OpFX1E: add V%X to Index 0x%03X, new value: 0x%03X; Overflow 0x%X", x, c.Index, new_i, c.Registers[0xF])
	c.Index = new_i
}

// TODO: opFX0A blocks until key X is pressed (reduces program counter by 2 if it
// is NOT pressed)

// TODO: opFX29 sets the Index to the address of the hex character in VX (look up
// from font table)

// TODO: FX33

// TODO: FX55

// TODO: FX65
