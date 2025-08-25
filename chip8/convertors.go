package chip8

func uint16FromTwoBytes(leftByte, rightByte byte) uint16 {
	return (uint16)(leftByte)<<8 | (uint16)(rightByte)
}
