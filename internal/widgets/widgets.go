package widgets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TextInput struct {
	Text	 string
	OnSubmit func(text string)
}

func NewTextInput(onSubmit func(string)) *TextInput {
	return &TextInput{OnSubmit: onSubmit}
}

func (ti *TextInput) Update() {
	for _, c := range ebiten.AppendInputChars(nil) {
		ti.Text += string(c)
	}

	// backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(ti.Text) > 0 {
		rs := []rune(ti.Text)
		ti.Text = string(rs[:len(rs)-1])
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if ti.OnSubmit != nil {
			ti.OnSubmit(ti.Text)
		}
	}
}

func (ti *TextInput) Draw(screen *ebiten.Image, x, y int) {
	ebitenutil.DebugPrintAt(screen, "URL: "+ti.Text, x, y)
}
