// main.go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/deglebe/webseek/internal/app"
)

func main() {
	game, err := app.NewBrowser("")
	if err != nil {
		log.Fatal("failed to initialize browser:", err)
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("webseek")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
