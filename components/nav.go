package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/grouter"
)

type Nav struct {
	*gr.This
	Brand string
}

// Implements the Renderer interface.
func (c Nav) Render() gr.Component {
	elem := el.Div(gr.CSS("nav-wrapper"),
		el.ListItem(
			gr.CSS("nav-brand"),
			el.Italic(gr.CSS("fa", "fa-cogs")),
			gr.Text(" "),
			grouter.Link("/", c.Brand),
		),
		el.UnorderedList(
			gr.CSS("nav-menu", "nav-pills", "nav-stacked"),

			c.createLinkListItem("/instances", "Instances"),
			c.createLinkListItem("/volumes", "Volumes"),
			c.createLinkListItem("/images", "Images"),
		),
	)

	return elem
}

func (c Nav) createLinkListItem(path, title string) gr.Modifier {

	return el.ListItem(
		grouter.MarkIfActive(c.Props(), path),
		grouter.Link(path, title),
	)

}
