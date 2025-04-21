package input

import "github.com/hajimehoshi/ebiten/v2"

func IsKeyPressed(r rune) bool {
	return ebiten.IsKeyPressed(ebiten.Key(r))
}
