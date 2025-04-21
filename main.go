package main

import (
	"os"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/deglebe/webseek/internal/app"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <file>", os.Args[0])
	}
	filename := os.Args[1]

	game, err := app.NewBrowser(filename)
	if err != nil {
		log.Fatal("Failed to load page:", err)
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("webseek")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
