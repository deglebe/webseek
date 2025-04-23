package app

import (
	"os"
	"io"
	"fmt"
	"image"
	"strings"
	"net/url"
	"net/http"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/deglebe/browse/pkg/html"
	"github.com/deglebe/browse/pkg/layout"
	"github.com/deglebe/browse/pkg/dom"
	"github.com/deglebe/webseek/internal/state"
	"github.com/deglebe/webseek/internal/widgets"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const statusBarHeight = 16

type Browser struct {
	st	*state.State
	w, h	int
	urlBar	*widgets.TextInput
	face	font.Face
}

func NewBrowser(path string) (*Browser, error) {
	var root *dom.Node
	var url string

	if path != "" {
		f, err := os.Open(path)
		if err != nil { return nil, err }
		defer f.Close()

		parser := html.NewParser(f)
		r, err := parser.Parse()
		if err != nil && err != io.EOF { return nil, err }
		root = r
		url = fmt.Sprintf("file://%s", path)
	} else {
		root = &dom.Node{
			Type:	 dom.ElementNode,
			Data:	 "root",
			Attrs:	map[string]string{},
			Children: []*dom.Node{},
		}
		url = ""
	}

	tab := state.Tab{
		URL:	url,
		DOM:	root,
		Scroll: 0,
	}
	st := &state.State{
		Tabs:	   []state.Tab{tab},
		CurrentTab: 0,
	}

	b := &Browser{
		st:	 st,
		urlBar: widgets.NewTextInput(nil),
		face:   basicfont.Face7x13,
	}

	b.urlBar.OnSubmit = func(text string) {
		b.loadURL(text)
	}

	return b, nil
}

func (b *Browser) loadURL(raw string) {
	var reader io.Reader
	if u, err := url.Parse(raw); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		resp, err := http.Get(raw)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch error: %v\n", err)
			return
		}
		reader = resp.Body
	} else {
		// file fallback
		fp := raw
		if strings.HasPrefix(raw, "file://") {
			fp = raw[len("file://"):]
		}
		f, err := os.Open(fp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open file error: %v\n", err)
			return
		}
		reader = f
	}

	parser := html.NewParser(reader)
	root, err := parser.Parse()
	if err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		return
	}

	tab := &b.st.Tabs[b.st.CurrentTab]
	tab.URL = raw
	tab.DOM = root
	tab.Scroll = 0
	b.relayout()
}

func (b *Browser) Update() error {
	b.urlBar.Update()

	tab := &b.st.Tabs[b.st.CurrentTab]

	// scroll down
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if tab.Scroll < tab.ContentH-b.h {
			tab.Scroll++
		}
	}

	// scroll up
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if tab.Scroll > 0 {
			tab.Scroll--
		}
	}
	return nil
}

func (b *Browser) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x20, 0x20, 0x20, 0xFF})

	urlText := "URL: " + b.urlBar.Text
	dr := &font.Drawer{
		Dst:  screen,
		Src:  image.NewUniform(color.White),
		Face: b.face,
		Dot: fixed.Point26_6{
			X: fixed.I(0),
			Y: fixed.I(b.face.Metrics().Ascent.Round()),
		},
	}
	dr.DrawString(urlText)

	tab := &b.st.Tabs[b.st.CurrentTab]
	for _, op := range tab.Ops {
		// cull above
		if op.Y < tab.Scroll {
			continue
		}
		// cull below
		if op.Y >= tab.Scroll+b.h-statusBarHeight {
			continue
		}

		dr := &font.Drawer{
			Dst:  screen,
			Src:  image.NewUniform(color.White),
			Face: op.Face,
			Dot: fixed.Point26_6{
				X: fixed.I(op.X),
				Y: fixed.I(op.Y - tab.Scroll + op.Face.Metrics().Ascent.Round() + statusBarHeight),
			},
		}
		dr.DrawString(op.Text)
	}

	status := fmt.Sprintf("Scroll %d/%d  %s", tab.Scroll, tab.ContentH, tab.URL)
	dr = &font.Drawer{
		Dst:  screen,
		Src:  image.NewUniform(color.White),
		Face: b.face,
		Dot: fixed.Point26_6{
			X: fixed.I(0),
			Y: fixed.I(b.h - statusBarHeight + b.face.Metrics().Ascent.Round()),
		},
	}
	dr.DrawString(status)
}

func (b *Browser) Layout(width, height int) (int, int) {
	if width != b.w || height != b.h {
		b.w, b.h = width, height
		b.relayout()
	}
	return width, height
}

func (b *Browser) relayout() {
	tab := &b.st.Tabs[b.st.CurrentTab]
	ctx, err := layout.NewContext(b.w)
	if err != nil {
		panic(err)
	}
	ops, totalH := layout.Render(tab.DOM, ctx)
	tab.Ops = ops
	tab.ContentH = totalH
}
