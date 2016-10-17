package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
)

type NewClass struct {
	*gr.This
}

// Implements the StateInitializer interfacn.
func (n NewClass) GetInitialState() gr.State {
	return gr.State{"step": 1, "querying": false, "error": "", "className": ""}
}

func (n NewClass) Render() gr.Component {

	state := n.State()
	props := n.Props()

	// Response placeholder
	response := el.Div()

	if state.Int("step") == 1 {

		// STEP 1

		nextListener := func(event *gr.Event) {
			n.SetState(gr.State{"step": 2})
		}

		storeName := func(event *gr.Event) {
			n.SetState(gr.State{"className": event.TargetValue()})
		}

		el.Form(
			el.Div(
				gr.CSS("form-group"),
				el.Label(
					gr.Text("Name"),
				),
				el.Input(
					attr.Type("name"),
					attr.ClassName("form-control"),
					attr.ID("name"),
					attr.Placeholder("Class Name"),
					attr.Value(state.String("className")),
					evt.Change(storeName),
				),
			),
			el.Button(
				//attr.Type("submit"),
				evt.Click(nextListener).PreventDefault(),
				gr.CSS("btn", "btn-primary"),
				gr.Text("Next"),
			),
		).Modify(response)

	} else if state.Int("step") == 2 {

		// STEP 2

		NewClassFormBuilder(state.String("className"), props.String("classType")).Modify(response)

		backListener := func(event *gr.Event) {
			n.SetState(gr.State{"step": 1})
		}

		el.Div(
			gr.CSS("btn-toolbar"),

			// Back Button
			el.Button(
				attr.Type("button"),
				evt.Click(backListener).PreventDefault(),
				gr.CSS("btn", "btn-secondary"),
				gr.Text("Back"),
			),

			// Save Button
			el.Button(
				attr.Type("button"),
				evt.Click(backListener).PreventDefault(),
				gr.CSS("btn", "btn-primary"),
				gr.Text("Save"),
			),
		).Modify(response)

	}

	return response
}
