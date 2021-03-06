package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/gopherjs/jquery"
)

var jQuery = jquery.NewJQuery

type Modal struct {
	*gr.This
}

func hideAllModals() {
	jQuery(".modal").Call("modal", "hide")
}

func hideModal(id string) {
	jQuery("#"+id).Call("modal", "hide")
}

func (m Modal) onShow(event *gr.Event) {
	//println("onShow")
}

func (m Modal) Render() gr.Component {

	props := m.Props()

	modal := el.Div(
		gr.CSS("modal", "fade"),
		attr.ID(m.Props().String("id")),
		attr.Role("dialog"),
		attr.Ref(m.Props().String("id")),
	)

	modalDialog := el.Div(
		gr.CSS("modal-dialog"),
	)

	content := el.Div(
		gr.CSS("modal-content"),
	)

	//Header
	el.Div(
		gr.CSS("modal-header"),
		el.Div(
			gr.CSS("modal-title"),
			el.Header4(gr.Text(props.String("title"))),
		),
	).Modify(content)

	// Body
	el.Div(
		gr.CSS("modal-body"),
		m.Children().Element(),
	).Modify(content)

	/*
		// TODO make this whole thing more reuseable
			// Footer
			if m.This.Component("footer") != nil {
				el.Div(
					gr.CSS("modal-footer"),
					//
				).Modify(content)
			}
	*/

	content.Modify(modalDialog)
	modalDialog.Modify(modal)

	return modal

}
