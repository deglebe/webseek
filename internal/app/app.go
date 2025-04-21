package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/deglebe/browse/pkg/html"
	"os"
)

type Browser struct {
	// TODO: tabs, url bar text, focus state
}

func NewBrowser() *Browser {
	return &Browser{
		// initialize state
	}
}

func (b *Browser) Update() error {
	// TODO: handle keyboard/mouse, navigation
	return nil
}

func (b *Browser) Draw(screen *ebiten.Image) {
	// TODO: layout and paint using browse engine + canvas package
}

func (b *Browser) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
