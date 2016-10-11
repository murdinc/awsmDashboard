package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
)

type Layout struct {
	*gr.This
	Pages
	Brand       string
	ActivePage  string
	ApiEndpoint string
	Content     *gr.ReactComponent
}

// Implements the Renderer interface.
func (l Layout) Render() gr.Component {
	return el.Div(
		gr.CSS("main-wrapper"),

		// Nav
		gr.New(&Nav{Brand: l.Brand, Pages: l.Pages}).CreateElement(l.This.Props()),

		//Content
		gr.New(&Content{}).CreateElement(gr.Props{"Title": l.ActivePage, "ApiEndpoint": l.Pages[l.ActivePage].ApiEndpoint}),
	)
}
