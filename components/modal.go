package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
)

type Modal struct {
	*gr.This
	Id     string
	Title  string
	Body   *gr.Component
	Footer string
}

func (m Modal) Render() gr.Component {

	modal := el.Div(
		gr.CSS("modal", "fade"),
		attr.ID(m.Id),
		attr.Role("dialog"),
	)

	modalDialog := el.Div(
		gr.CSS("modal-dialog"),
	)

	// Content
	el.Div(
		gr.CSS("modal-content"),
		el.Div( //Header
			gr.CSS("modal-header"),
			el.Div(
				gr.CSS("modal-title"),
				el.Header4(gr.Text(m.Title)),
			),
		),
		el.Div( // Body
			gr.CSS("modal-body"),
			el.Div(
				gr.CSS("modal-title"),
				gr.Text("Some Text in the Modal Body"),
			),
		),
		el.Div( // Footer
			gr.CSS("modal-footer"),
			el.Div(
				gr.CSS("modal-title"),
				gr.Text(m.Footer),
			),
		),
	).Modify(modalDialog)

	modalDialog.Modify(modal)

	return modal

}
