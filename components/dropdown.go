package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
)

type Dropdown struct {
	*gr.This
}

func (d Dropdown) Render() gr.Component {

	elem := el.Div(
		gr.CSS("btn-group", "dropdown"),
		el.Button(
			gr.CSS("btn", "btn-primary", "btn-xs", "dropdown-toggle"),
			el.Italic(gr.CSS("fa", "fa-gear")),
			gr.Data("toggle", "dropdown"),
			gr.Aria("haspopup", "true"),
			gr.Aria("expanded", "true"),
		),
		el.UnorderedList(
			gr.CSS("dropdown-menu"),
			el.ListItem(el.Anchor(attr.HRef("localhost"), gr.Text("New Instance"))),
			el.ListItem(el.Anchor(attr.HRef("localhost"), gr.Text("Edit Classes"))),
		),
	)

	return elem
}
