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

	resp := el.Div(
		gr.CSS("content-wrapper"),
	)

	// Header
	header := el.Div(gr.CSS("content-header"),
		el.Header1(
			gr.Text(c.Props().String("activePage")+" "),
		),
	)
	if c.Page.HasClasses {
		gr.New(&AssetDropdownMenu{}).CreateElement(gr.Props{"type": c.Page.Type, "apiType": c.Page.ApiType}).Modify(header)
	}
	header.Modify(resp)

	// Asset Table
	gr.New(&AssetTable{}).CreateElement(gr.Props{"apiType": c.Page.ApiType}).Modify(resp)

	return resp
}
