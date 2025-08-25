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

func TestHasClearScreen(t *testing.T) {
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

func TestHasJumpInstruction(t *testing.T) {
	rom := []byte{0x12, 0x34}
	emu, _ := NewChip8FromByte(rom)
	var want uint16 = 0x0234

	err := emu.Update()
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}
	got := emu.PC

	if emu.PC != want {
		t.Errorf("1NNN instruction did not advance program counter, wanted %X, got %X", want, got)
	}
}

func TestErrorsOnBadInstruction(t *testing.T) {
	rom := []byte{0xF0, 0xF0}
	emu, _ := NewChip8FromByte(rom)
	err := emu.Update()
	if err == nil {
		t.Errorf("expected error, recieved nil")
	}
}

func TestSubroutine(t *testing.T) {
	rom := []byte{0x23, 0x4F}
	t.Run("pushes correct address onto stack", func(t *testing.T) {
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
	})

	t.Run("popping stack returns to correct pointer", func(t *testing.T) {
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

	})
}
