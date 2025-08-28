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

// TODO: Figure out how to do this ~700 a second
func (g *Game) Update() error {
	g.emu.Update()
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
