package chip8

// clear the screen
func (c *Chip8) clearDisplay() {
	var blankDisplay [64][32]bool
	c.Display = blankDisplay
}

func (c *Chip8) jump(location uint16) {
	c.PC = location
}
func (c *Chip8) stackPush(address uint16) {
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
