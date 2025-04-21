package app

import (
	"os"
	"io"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"golang.org/x/image/font/basicfont"

	"github.com/deglebe/browse/pkg/html"
	"github.com/deglebe/browse/pkg/layout"

	"github.com/deglebe/webseek/internal/state"
)

const (
	statusBarHeight = 16
)

type Browser struct {
	st *state.State
	w, h int
}

var bgColor = color.RGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}

func NewBrowser(path string) (*Browser, error) {
	f, err := os.Open(path)
	if err != nil { return nil, err }
	defer f.Close()

	parser := html.NewParser(f)
	root, err := parser.Parse()
	if err != nil && err != io.EOF { return nil, err }

	tab := state.Tab{
		URL:	fmt.Sprintf("file://%s", path),
		DOM:	root,
		Scroll: 0,
	}
	st := &state.State{
		Tabs:		[]state.Tab{tab},
		CurrentTab:	0,
	}

	b := &Browser{st: st}
	b.relayout()
	return b, nil
}

func (b *Browser) Update() error {
	tab := &b.st.Tabs[b.st.CurrentTab]
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if tab.Scroll < tab.ContentH - b.h { tab.Scroll++ }
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		tab := &b.st.Tabs[b.st.CurrentTab]
		if tab.Scroll > 0 { tab.Scroll-- }
	}
	return nil
}

func (b *Browser) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)

	tab := b.st.Tabs[b.st.CurrentTab]

	for _, op := range tab.Ops {
		if op.Y < tab.Scroll { continue }
		if op.Y >= tab.Scroll + b.h - statusBarHeight { continue }
		ebitenutil.DebugPrintAt(screen, op.Text, op.X, op.Y - tab.Scroll)
	}

	status := fmt.Sprintf("[%d/%d] %s", tab.Scroll, tab.ContentH, tab.URL)
	ebitenutil.DebugPrintAt(screen, status, 0, b.h-16)
}


func (b *Browser) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	b.w, b.h = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}

func (b *Browser) relayout() {
	tab := &b.st.Tabs[b.st.CurrentTab]
	face := basicfont.Face7x13
	ctx := &layout.Context{
		Face:		face,
		LineHeight:	face.Metrics().Height.Round(),
		MaxWidth:	b.w,
		ListIndent:	20,
	}
	ops, totalH := layout.Render(tab.DOM, ctx)
	tab.Ops = ops
	tab.ContentH = totalH
}

