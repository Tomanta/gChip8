package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	time int
}

func (g *Game) Update() error {
	g.time += 1
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	dbgMsg := fmt.Sprintf("Hello! Current counter: %d", g.time)
	ebitenutil.DebugPrint(screen, dbgMsg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{time: 0}); err != nil {
		log.Fatal(err)
	}
}
