package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
)

type Layout struct {
	*gr.This
	Pages
	Brand      string
	ActivePage string
}

type Pages map[string]Page

type Page struct {
	Route           string
	DropdownOptions []DropdownOption
	ClassEndpoint   string
	AssetEndpoint   string
}

// Implements the Renderer interface.
func (l Layout) Render() gr.Component {

	return el.Div(
		gr.CSS("main-wrapper"),
		// Nav
		gr.New(&Nav{Brand: l.Brand, Pages: l.Pages}).CreateElement(l.Props()), // layout passes the router to the nav

		//Content
		gr.New(&Content{Page: l.Pages[l.ActivePage]}).CreateElement(gr.Props{"activePage": l.ActivePage}),
	)
}
