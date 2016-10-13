package components

import (
	"fmt"

	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsmDashboard/helpers"
)

type Modal struct {
	*gr.This
	Id            string
	Title         string
	Body          *gr.Component
	Footer        string
	ClassEndpoint string
}

// Implements the StateInitializer interface.
func (m Modal) GetInitialState() gr.State {
	return gr.State{"querying": false, "error": nil, "assetList": nil}
}

func (m Modal) Render() gr.Component {

	// Asset List placeholder
	response := el.Div()

	modal := el.Div(
		gr.CSS("modal", "fade"),
		attr.ID(m.Id),
		attr.Role("dialog"),
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
			el.Header4(gr.Text(m.Title)),
		),
	).Modify(content)

	// Body
	el.Div(
		gr.CSS("modal-body"),
		response,
	).Modify(content)

	// Footer
	if m.Footer != "" {
		el.Div(
			gr.CSS("modal-footer"),
			gr.Text(m.Footer),
		).Modify(content)
	}

	if classes := m.State().Interface("classList"); classes != nil {
		classList := ClassListBuilder(classes) // Build the class list
		classList.Modify(response)
	} else if m.State().Bool("querying") {
		gr.Text("Loading...").Modify(response)
	} else if errStr := m.State().Interface("error"); errStr != nil {
		gr.Text(errStr).Modify(response)
	} else {
		gr.Text("Nothing here!").Modify(response)
	}

	content.Modify(modalDialog)

	modalDialog.Modify(modal)

	return modal

}

// Implements the ComponentDidMount interface
func (m Modal) ComponentDidMount() {

	if endpoint := m.ClassEndpoint; endpoint != "" {

		m.SetState(gr.State{"querying": true})

		resp, err := helpers.QueryAPI("//localhost:8081/api" + endpoint)
		if !m.IsMounted() {
			return
		}
		if err != nil {
			m.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		m.SetState(gr.State{"querying": false, "classList": resp})

	}
}

// Implements the ShouldComponentUpdate interface.
func (m Modal) ShouldComponentUpdate(this *gr.This, next gr.Cops) bool {
	return m.State().HasChanged(next.State, "classList") &&
		m.State().HasChanged(next.State, "querying") &&
		m.State().HasChanged(next.State, "error")
}
