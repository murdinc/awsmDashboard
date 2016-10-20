package components

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
	"github.com/murdinc/awsmDashboard/helpers"
)

type NewClass struct {
	*gr.This
}

// Implements the StateInitializer interface
func (n NewClass) GetInitialState() gr.State {
	return gr.State{"step": 1, "querying": false, "error": "", "className": ""}
}

func (n NewClass) checkClassName(className string) {
	n.SetState(gr.State{"querying": true, "error": ""})

	go func() {

		// Make sure the classname isn't empty
		if className == "" {
			n.SetState(gr.State{"error": "Class name cannot be empty", "querying": false})
			return
		}

		// Make sure this class name doesn't already exist
		if apiType := n.Props().String("apiType"); apiType != "" {
			endpoint := "//localhost:8081/api/classes/" + apiType + "/name/" + className
			resp, err := helpers.QueryAPI(endpoint)

			if err != nil {
				n.SetState(gr.State{"error": fmt.Sprintf("Error while querying endpoint: %s", endpoint), "querying": false})
				return
			}

			jsonParsed, _ := gabs.ParseJSON(resp)
			exists := jsonParsed.S("success").Data().(bool)

			if exists {
				n.SetState(gr.State{"error": "Class name " + className + " already exists!", "querying": false})
				return
			}
		} else {
			n.SetState(gr.State{"error": "No API type, unable to query API", "querying": false})
			return
		}
		n.SetState(gr.State{"error": "", "querying": false, "step": 2})
	}()
}

func (n NewClass) Render() gr.Component {

	state := n.State()
	props := n.Props()

	// Response placeholder
	response := el.Div()

	// Print any errors
	helpers.ErrorElem(state.String("error")).Modify(response)

	if state.Int("step") == 1 {

		// STEP 1

		// `Next` Button Listener
		nextListener := func(*gr.Event) {
			n.checkClassName(state.String("className"))
		}

		// Store classname on event change
		storeName := func(event *gr.Event) {
			n.SetState(gr.State{"className": event.TargetValue()})
		}

		form := el.Form(
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
		)

		// Disables the form while querying
		if state.Bool("querying") {
			attr.Disabled("").Modify(form)
		}

		form.Modify(response)

	} else if state.Int("step") == 2 {

		// STEP 2

		NewClassFormBuilder(state.String("className"), props.String("apiType")).Modify(response)

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
