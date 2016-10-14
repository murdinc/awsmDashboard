package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
)

type Content struct {
	*gr.This
	Page Page
}

func (c Content) Render() gr.Component {

	return el.Div(gr.CSS("content-wrapper"),
		el.Div(gr.CSS("content-header"),
			el.Header1(
				gr.Text(c.Props().String("activePage")+" "),
			),
			gr.New(&Dropdown{DropdownOptions: c.Page.DropdownOptions, ClassEndpoint: c.Page.ClassEndpoint}).CreateElement(gr.Props{}),
		),
		gr.New(&AssetTable{AssetEndpoint: c.Page.AssetEndpoint}).CreateElement(gr.Props{}),
	)
}
