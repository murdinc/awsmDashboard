package components

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
	"github.com/murdinc/awsmDashboard/components/forms"
	"github.com/murdinc/awsmDashboard/helpers"
)

type NewClass struct {
	*gr.This
}

// Implements the StateInitializer interface
func (n NewClass) GetInitialState() gr.State {
	return gr.State{"step": 1, "querying": false, "error": "", "className": ""}
}

func (n NewClass) Render() gr.Component {

	state := n.State()
	props := n.Props()

	// Response placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)

	if state.Int("step") == 1 {

		// STEP 1

		newClassForm := el.Form()

		forms.TextField("Name", "className", state.String("className"), n.storeValue).Modify(newClassForm)

		buttons := el.Div(
			gr.CSS("btn-toolbar"),
		)

		// Close
		el.Button(
			evt.Click(n.closeButton).PreventDefault(),
			gr.CSS("btn", "btn-secondary"),
			gr.Text("Close"),
		).Modify(buttons)

		// Next
		el.Button(
			evt.Click(n.stepOneNext).PreventDefault(),
			gr.CSS("btn", "btn-primary"),
			gr.Text("Next"),
		).Modify(buttons)

		// Disables the form while querying
		// but does it work?
		if state.Bool("querying") {
			attr.Disabled("").Modify(newClassForm)
		}

		newClassForm.Modify(response)
		buttons.Modify(response)

	} else if state.Int("step") == 2 {

		// STEP 2

		classForm := NewClassFormBuilder(props.String("apiType"))

		classForm.CreateElement(gr.Props{
			"className":     state.String("className"),
			"backButton":    n.stepTwoBack,
			"apiType":       props.String("apiType"),
			"hideAllModals": hideAllModals,
			"newClass":      true,
		}).Modify(response)

	}

	return response
}

func (n NewClass) storeValue(event *gr.Event) {
	id := event.Target().Get("id").String()
	n.SetState(gr.State{id: event.TargetValue()})
}

func (n NewClass) closeButton(*gr.Event) {
	n.SetState(gr.State{"success": ""})
	hideAllModals()
}

func (n NewClass) stepOneNext(*gr.Event) {

	n.SetState(gr.State{"querying": true, "error": ""})

	go func(className string) {

		// Make sure the classname isn't empty
		if className == "" {
			n.SetState(gr.State{"error": "Class name cannot be empty", "querying": false})
			return
		}

		// Make sure the classname doesn't include any numbers
		for _, char := range className {
			if char >= '0' && char <= '9' {
				n.SetState(gr.State{"error": "Class name cannot contain numbers", "querying": false})
				return
			}
		}

		// Make sure this class name doesn't already exist
		if apiType := n.Props().String("apiType"); apiType != "" {
			endpoint := "//localhost:8081/api/classes/" + apiType + "/name/" + className
			resp, err := helpers.GetAPI(endpoint)

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
	}(n.State().String("className"))

}

func (n NewClass) stepTwoBack(event *gr.Event) {
	n.SetState(gr.State{"step": 1})
}
