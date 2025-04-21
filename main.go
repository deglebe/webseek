package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/deglebe/webseek/internal/app"
)

func main() {
	game := app.NewBrowser()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("webseek")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
