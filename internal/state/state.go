package state

import (
	"github.com/deglebe/browse/pkg/dom"
	"github.com/deglebe/browse/pkg/layout"
)

type Tab struct {
	URL		string
	DOM		*dom.Node
	Ops		[]layout.RenderOp
	ContentH	int
	Scroll		int
	// layout tree, styles, scroll offsets
}

type State struct {
	Tabs		[]Tab
	CurrentTab	int
}
