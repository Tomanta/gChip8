package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tomanta/echip8/chip8"
)

type Game struct {
	emu chip8.Chip8
}

func (g *Game) getKeys() []byte {
	var keys []byte
	// 123C
	// 456D
	// 789E
	// A0BF
	if ebiten.IsKeyPressed(ebiten.Key1) {
		keys = append(keys, 0x1)
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		keys = append(keys, 0x2)
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		keys = append(keys, 0x3)
	}
	if ebiten.IsKeyPressed(ebiten.Key4) {
		keys = append(keys, 0xC)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		keys = append(keys, 0x4)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		keys = append(keys, 0x5)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		keys = append(keys, 0x6)
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		keys = append(keys, 0xD)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		keys = append(keys, 0x7)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		keys = append(keys, 0x8)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		keys = append(keys, 0x9)
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		keys = append(keys, 0xE)
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		keys = append(keys, 0xA)
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		keys = append(keys, 0x0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		keys = append(keys, 0xB)
	}
	if ebiten.IsKeyPressed(ebiten.KeyV) {
		keys = append(keys, 0xF)
	}

	return keys
}

// TODO: Figure out how to do this ~700 a second
func (g *Game) Update() error {
	g.emu.SetKeysPressed(g.getKeys())
	g.emu.Update()
	// if g.emu.DebugMsg != "TODO" {
	// 	fmt.Println(g.emu.DebugMsg)
	// }
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var scale float32 = 10
	screen.Clear()
	for x := range 64 {
		var x_pos float32 = (float32)(x) * scale
		for y := range 32 {
			var y_pos float32 = (float32)(y) * scale
			if g.emu.Display[x][y] {
				vector.DrawFilledRect(screen, x_pos, y_pos, 10, 10, color.RGBA{51, 255, 51, 0}, false)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 320
}

func openRom() []byte {
	// testRom := "./roms/ibm_logo.ch8"
	testRom := "./roms/test_opcode.ch8"
	data, err := os.ReadFile(testRom)
	if err != nil {
		panic("could not open rom")
	}
	return data
}

func main() {
	ebiten.SetWindowSize(640, 320)
	romData := openRom()
	emu, _ := chip8.NewChip8FromByte(romData)

	if err := ebiten.RunGame(&Game{emu: emu}); err != nil {
		log.Fatal(err)
	}
}
