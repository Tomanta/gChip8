package chip8

import (
	"os"
	"reflect"
	"testing"
)

const (
	StartProgramCounter    uint16 = 0x200
	StartRomMemoryLocation uint16 = 0x200
)

func openTestRom(t testing.TB) []byte {
	t.Helper()
	testRom := "../roms/ibm_logo.ch8"
	data, err := os.ReadFile(testRom)
	if err != nil {
		t.Fatalf("could not open test rom '%s', received error: %v", testRom, err)
	}
	return data
}

func getIBMEmulator(t testing.TB) Chip8 {
	t.Helper()
	romData := openTestRom(t)

	emu, err := NewChip8FromByte(romData)
	if err != nil {
		t.Fatalf("could not get emulator from rom file, received: %v", err)
	}
	return emu
}

func TestCanCreateFromBytes(t *testing.T) {
	t.Run("can load a rom file to memory", func(t *testing.T) {
		romData := openTestRom(t)
		romDataLength := uint16(len(romData))
		got, err := NewChip8FromByte(romData)

		if err != nil {
			t.Fatalf("Received error creating file: %v", err)
		}

		gotr := got.Memory[StartRomMemoryLocation : StartRomMemoryLocation+romDataLength]

		if !reflect.DeepEqual(gotr, romData) {
			t.Errorf("ROM not loaded to memory. want %X, got %X", romData, gotr)
		}
	})

	t.Run("empty rom returns error", func(t *testing.T) {
		romData := []byte{}
		_, err := NewChip8FromByte(romData)
		if err == nil {
			t.Fatalf("expected error, did not receive one")
		}
	})

	t.Run("program counter set to 0x200", func(t *testing.T) {
		want := 0x200
		chip8 := getIBMEmulator(t)
		if chip8.PC != 0x200 {
			t.Errorf("program counter not set. want %X, got %X", want, chip8.PC)
		}
	})

	t.Run("fonts are loaded into memory starting at 0x50", func(t *testing.T) {
		chip8 := getIBMEmulator(t)
		var want1 byte = 0xF0
		var want1_loc uint16 = 0x50
		var got1 byte = chip8.Memory[want1_loc]
		if got1 != want1 {
			t.Errorf("Expected Memory 0x%02X to be 0x%02X, got 0x%02X", want1_loc, want1, got1)
		}

		var want2 byte = 0x80
		var want2_loc uint16 = 0x129
		var got2 byte = chip8.Memory[want2_loc]
		if got1 != want1 {
			t.Errorf("Expected Memory 0x%02X to be 0x%02X, got 0x%02X", want2_loc, want2, got2)
		}

	})

	t.Run("initial display is blank", func(t *testing.T) {
		want := [64][32]bool{}
		got := getIBMEmulator(t).Display

		if !reflect.DeepEqual(got, want) {
			t.Errorf("display not blank")
		}
	})
}

func TestCanFetchInstructions(t *testing.T) {
	emu := getIBMEmulator(t)
	emu.Update()
	var wantPC uint16 = 0x0202
	var wantInstr uint8 = 0xA2

	if emu.PC != wantPC {
		t.Errorf("program counter not advanced, wanted %X, got %X", wantPC, emu.PC)
	}

	if emu.Memory[emu.PC] != wantInstr {
		t.Errorf("instruction at next program counter incorrect, wanted %X, got %X", wantInstr, emu.Memory[emu.PC])
	}
}

func TestErrorsOnBadInstruction(t *testing.T) {
	rom := []byte{0xFF, 0xFF}
	emu, _ := NewChip8FromByte(rom)
	err := emu.Update()
	if err == nil {
		t.Errorf("expected error, recieved nil")
	}
}

func TestOp00E0(t *testing.T) {
	var want [64][32]bool
	var dirtyDisplay [64][32]bool
	dirtyDisplay[5][1] = true

	emu := getIBMEmulator(t)
	emu.Display = dirtyDisplay
	emu.Update()
	got := emu.Display
	if !reflect.DeepEqual(got, want) {
		t.Error("0x00E0 instruction did not clear display")
	}
}

func TestOp00EE(t *testing.T) {
	rom := []byte{0x23, 0x4F}
	emu, _ := NewChip8FromByte(rom)
	// rig the stack
	emu.Memory[0x034F] = 0x00
	emu.Memory[0x034F+1] = 0xEE
	emu.Update()
	emu.Update()
	var pc_want uint16 = 0x0202
	var stack_want uint16 = 0x000
	pc_got := emu.PC
	stack_got := emu.Stack[0]

	if stack_got != stack_want {
		t.Errorf("expected stack[0] to have 0x%03X, got 0x%03X", stack_want, stack_got)
	}

	if pc_got != pc_want {
		t.Errorf("expected program counter to have 0x%03X, got 0x%03X", pc_want, pc_got)
	}

}

var cases = []struct {
	name        string
	rom         []byte
	num_updates int
	want        uint16
	got         func(emu Chip8) uint16
}{
	{name: "op1NNN jumps to memory location NNN", rom: []byte{0x12, 0x34}, num_updates: 1, want: 0x0234, got: func(emu Chip8) uint16 { return emu.PC }},
	{name: "op8XY0 sets X to Y", rom: []byte{0x61, 0x82, 0x62, 0x85, 0x81, 0x20}, num_updates: 3, want: 0x85, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[1]) }},
	{name: "op8XY1 binary OR X and Y", rom: []byte{0x61, 0x45, 0x62, 0x32, 0x81, 0x21}, num_updates: 3, want: 0x45 | 0x32, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[1]) }},
	{name: "op8XY2 binary AND X and Y", rom: []byte{0x61, 0x45, 0x62, 0x42, 0x81, 0x22}, num_updates: 3, want: 0x45 & 0x42, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[1]) }},
	{name: "op8XY3 binary XOR X and Y", rom: []byte{0x61, 0x45, 0x62, 0x42, 0x81, 0x23}, num_updates: 3, want: 0x45 ^ 0x42, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[1]) }},
	{name: "op8XY4 adds X and Y into X", rom: []byte{0x61, 0x45, 0x62, 0x42, 0x81, 0x24}, num_updates: 3, want: 0x45 + 0x42, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[1]) }},
	{name: "op8XY4 sets overflow flag = 1 if overflow", rom: []byte{0x61, 0xBB, 0x62, 0x88, 0x81, 0x24}, num_updates: 3, want: 0x01, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[0xF]) }},
	{name: "op8XY5 subtracts Y from X and stores into X", rom: []byte{0x61, 0x88, 0x62, 0x42, 0x81, 0x25}, num_updates: 3, want: 0x88 - 0x42, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[1]) }},
	{name: "op8XY5 does not set overflow flag", rom: []byte{0x61, 0xBB, 0x62, 0x88, 0x81, 0x25}, num_updates: 3, want: 0x00, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[0xF]) }},
	{name: "op8XY6 shifts Y one bit to right, stores in X", rom: []byte{0x62, 0x10, 0x81, 0x26}, num_updates: 2, want: 0x10 >> 1, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[1]) }},
	{name: "op8XY6 sets flag to 0 if right bit is 0", rom: []byte{0x6F, 0x01, 0x62, 0x10, 0x81, 0x26}, num_updates: 3, want: 0, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[0xf]) }},
	{name: "op8XY6 sets flag to 1 if right bit is 1", rom: []byte{0x62, 0x11, 0x81, 0x26}, num_updates: 2, want: 1, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[0xf]) }},
	{name: "op8XY7 subtracts Y from X and stores into X", rom: []byte{0x61, 0x88, 0x62, 0x42, 0x81, 0x27}, num_updates: 3, want: 0x88 - 0x42, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[1]) }},
	{name: "op8XY7 does not set underflow if X > Y flag", rom: []byte{0x6F, 0x01, 0x61, 0xBB, 0x62, 0x88, 0x81, 0x27}, num_updates: 4, want: 0x00, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[0xF]) }},
	{name: "op8XY7 does set underflow flag if X < Y flag", rom: []byte{0x61, 0x88, 0x62, 0xBB, 0x81, 0x27}, num_updates: 3, want: 0x01, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[0xF]) }},
	{name: "op8XYE shifts Y one bit to left, stores in X", rom: []byte{0x62, 0xAA, 0x81, 0x2E}, num_updates: 2, want: 0x54, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[1]) }},
	{name: "op8XYE sets flag to 0 if left bit is 0", rom: []byte{0x6F, 0x01, 0x62, 0x10, 0x81, 0x2E}, num_updates: 3, want: 0, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[0xf]) }},
	{name: "op8XYE sets flag to 1 if left bit is 1", rom: []byte{0x62, 0x80, 0x81, 0x2E}, num_updates: 2, want: 1, got: func(emu Chip8) uint16 { return (uint16)(emu.Registers[0xf]) }},
}

func TestBasicInstructions(t *testing.T) {
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			emu, _ := NewChip8FromByte(test.rom)
			for range test.num_updates {
				emu.Update()
			}
			got := test.got(emu)
			if got != test.want {
				t.Errorf("Expected 0x%04X, got 0x%04X", test.want, got)
			}
		})
	}
}

func TestOp2NNN(t *testing.T) {
	rom := []byte{0x23, 0x4F}
	emu, _ := NewChip8FromByte(rom)
	emu.Update()
	var stack_want uint16 = 0x0202
	var pc_want uint16 = 0x034F
	stack_got := emu.Stack[0]
	pc_got := emu.PC

	if stack_got != stack_want {
		t.Errorf("expected stack[0] to have 0x%03X, got 0x%03X", stack_want, stack_got)
	}

	if pc_want != pc_got {
		t.Errorf("expected program counter 0x%03X, got 0x%03X", pc_want, pc_got)
	}
}

func TestOp3XNN(t *testing.T) {
	rom := []byte{0x61, 0x82, 0x31, 0x82, 0xFF, 0xFF, 0x82, 0xEE}
	emu, _ := NewChip8FromByte(rom)
	emu.Update() // Set register
	emu.Update() // Skip next instruction
	var want_byte byte = 0x82
	got_byte := emu.Memory[emu.PC]
	if got_byte != want_byte {
		t.Errorf("Expected next byte to be 0x%02X, got 0x%02X", want_byte, got_byte)
	}
}

func TestOp4XNN(t *testing.T) {
	rom := []byte{0x61, 0x82, 0x41, 0x85, 0xFF, 0xFF, 0x82, 0xEE}
	emu, _ := NewChip8FromByte(rom)
	emu.Update() // Set register
	emu.Update() // Skip next instruction
	var want_byte byte = 0x82
	got_byte := emu.Memory[emu.PC]
	if got_byte != want_byte {
		t.Errorf("Expected next byte to be 0x%02X, got 0x%02X", want_byte, got_byte)
	}
}

func TestOp5XY0(t *testing.T) {
	rom := []byte{0x61, 0x82, 0x62, 0x82, 0x51, 0x20, 0xFF, 0xFF, 0x88, 0x92}
	emu, _ := NewChip8FromByte(rom)
	emu.Update() // Set register X
	emu.Update() // Set register Y
	emu.Update() // Skip next instruction
	var want_byte byte = 0x88
	got_byte := emu.Memory[emu.PC]
	if got_byte != want_byte {
		t.Errorf("Expected next byte to be 0x%02X, got 0x%02X", want_byte, got_byte)
	}
}

func TestOp6XNN(t *testing.T) {
	rom := []byte{0x61, 0x82}
	emu, _ := NewChip8FromByte(rom)
	emu.Update()

	var want uint8 = 0x82
	got := emu.Registers[1]

	if got != want {
		t.Errorf("expected register 1 to contain 0x%02X, got 0x%02X", want, got)
	}
}

func TestOp7XNN(t *testing.T) {
	t.Run("can add basic register", func(t *testing.T) {
		rom := []byte{0x61, 0x82, 0x71, 0x11}
		emu, _ := NewChip8FromByte(rom)
		emu.Update()
		emu.Update()

		var want uint8 = 0x82 + 0x11
		got := emu.Registers[1]

		if got != want {
			t.Errorf("expected register 1 to contain 0x%02X, got 0x%02X", want, got)
		}
	})

	t.Run("overflow does not set overflow flag", func(t *testing.T) {
		rom := []byte{0x61, 0xFF, 0x71, 0x01}
		emu, _ := NewChip8FromByte(rom)
		emu.Update()
		emu.Update()

		var want uint8 = 0x00
		got := emu.Registers[1]

		if got != want {
			t.Errorf("expected register 1 to contain 0x%02X, got 0x%02X", want, got)
		}
	})
}

func TestOp9XY0(t *testing.T) {
	rom := []byte{0x61, 0x82, 0x62, 0x85, 0x91, 0x20, 0xFF, 0xFF, 0x88, 0x92}
	emu, _ := NewChip8FromByte(rom)
	emu.Update() // Set register X
	emu.Update() // Set register Y
	emu.Update() // Skip next instruction
	var want_byte byte = 0x88
	got_byte := emu.Memory[emu.PC]
	if got_byte != want_byte {
		t.Errorf("Expected next byte to be 0x%02X, got 0x%02X", want_byte, got_byte)
	}
}

func TestOpANNN(t *testing.T) {
	rom := []byte{0xA1, 0x22}
	emu, _ := NewChip8FromByte(rom)
	emu.Update()

	var want uint16 = 0x0122
	got := emu.Index
	if got != want {
		t.Errorf("expected index register 0x%03X, got 0x%03X", want, got)
	}
}
