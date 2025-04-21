package state

type Tab struct {
	URL	string
	DOM	interface{}
	// layout tree, styles, scroll offsets
}

type State struct {
	Tabs		[]Tab
	CurrentTab	int
}
