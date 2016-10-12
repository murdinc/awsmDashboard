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

// Implements the Renderer interface.
func (l Layout) Render() gr.Component {

	return el.Div(
		gr.CSS("main-wrapper"),

		// Nav
		gr.New(&Nav{Brand: l.Brand, Pages: l.Pages}).CreateElement(l.Props()), // layout passes the router to the nav

		//Content
		gr.New(&Content{}).CreateElement(gr.Props{"activePage": l.ActivePage, "apiEndpoint": l.Pages[l.ActivePage].ApiEndpoint},
			gr.New(&Dropdown{DropdownOptions: l.Pages[l.ActivePage].DropdownOptions}).CreateElement(gr.Props{}),
			//gr.New(&Dropdown{DropdownOptions: l.Pages[l.ActivePage].DropdownOptions}).CreateElement(gr.Props{}),
		),
	)
}

// Implements the ComponentWillReceiveProps interface
func (l Layout) ComponentWillReceiveProps(data gr.Cops) {
	println("ComponentWillReceiveProps")
}
