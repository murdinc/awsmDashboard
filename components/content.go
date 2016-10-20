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
			gr.New(&AssetDropdownMenu{}).CreateElement(gr.Props{"type": c.Page.Type, "apiType": c.Page.ApiType}),
		),
		gr.New(&AssetTable{}).CreateElement(gr.Props{"apiType": c.Page.ApiType}),
	)
}
