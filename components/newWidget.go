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

type NewWidget struct {
	*gr.This
}

// Implements the StateInitializer interface
func (n NewWidget) GetInitialState() gr.State {
	return gr.State{"step": 1, "querying": false, "error": "", "widgetName": ""}
}

func (n NewWidget) Render() gr.Component {

	state := n.State()
	props := n.Props()

	// Response placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)

	if state.Int("step") == 1 {

		// STEP 1

		newWidgetForm := el.Form()

		forms.TextField("Name", "widgetName", state.String("widgetName"), n.storeValue).Modify(newWidgetForm)
		forms.SelectOne("Type", "widgetType", []string{"rss"}, state.Interface("widgetType"), n.storeSelect).Modify(newWidgetForm)

		buttons := el.Div(
			gr.CSS("btn-toolbar"),
		)

		// Next
		el.Button(
			evt.Click(n.stepOneNext).PreventDefault(),
			gr.CSS("btn", "btn-primary"),
			gr.Text("Next"),
		).Modify(buttons)

		// Close
		el.Button(
			evt.Click(n.closeButton).PreventDefault(),
			gr.CSS("btn", "btn-secondary"),
			gr.Text("Close"),
		).Modify(buttons)

		// Disables the form while querying
		// but does it work?
		if state.Bool("querying") {
			attr.Disabled("").Modify(newWidgetForm)
		}

		newWidgetForm.Modify(response)
		buttons.Modify(response)

	} else if state.Int("step") == 2 {

		// STEP 2

		widgetForm := NewWidgetFormBuilder(state.String("widgetType"))

		widgetForm.CreateElement(gr.Props{
			"widgetName":    state.String("widgetName"),
			"widgetType":    state.String("widgetType"),
			"backButton":    n.stepTwoBack,
			"apiType":       props.String("apiType"),
			"hideAllModals": hideAllModals,
		}).Modify(response)

	}

	return response
}

func (n NewWidget) storeValue(event *gr.Event) {
	id := event.Target().Get("id").String()
	n.SetState(gr.State{id: event.TargetValue()})
}

func (n NewWidget) storeSelect(id string, val interface{}) {
	switch value := val.(type) {

	case map[string]interface{}:
		// single
		n.SetState(gr.State{id: value["value"]})

	case []interface{}:
		// multi
		var vals []string
		options := len(value)
		for i := 0; i < options; i++ {
			vals = append(vals, value[i].(map[string]interface{})["value"].(string))
		}
		n.SetState(gr.State{id: vals})

	default:
		n.SetState(gr.State{id: val})

	}
}

func (n NewWidget) closeButton(*gr.Event) {
	n.SetState(gr.State{"success": ""})
	hideAllModals()
}

func (n NewWidget) stepOneNext(*gr.Event) {

	n.SetState(gr.State{"querying": true, "error": ""})

	go func(widgetName string, widgetType string) {

		// Make sure the widgetName isn't empty
		if widgetName == "" {
			n.SetState(gr.State{"error": "Widget name cannot be empty", "querying": false})
			return
		}

		// Make sure the widgetName doesn't include any numbers
		for _, char := range widgetName {
			if char >= '0' && char <= '9' {
				n.SetState(gr.State{"error": "Widget name cannot contain numbers", "querying": false})
				return
			}
		}

		// Make sure this widget name doesn't already exist
		if apiType := n.Props().String("apiType"); apiType != "" {
			endpoint := "//localhost:8081/api/" + apiType + "/widgets/name/" + widgetName
			resp, err := helpers.GetAPI(endpoint)

			if err != nil {
				n.SetState(gr.State{"error": fmt.Sprintf("Error while querying endpoint: %s", endpoint), "querying": false})
				return
			}

			jsonParsed, _ := gabs.ParseJSON(resp)
			exists := jsonParsed.S("success").Data().(bool)

			if exists {
				n.SetState(gr.State{"error": "A widget named " + widgetName + " already exists!", "querying": false})
				return
			}
		} else {
			n.SetState(gr.State{"error": "No API type, unable to query API", "querying": false})
			return
		}

		// Make sure the widgetType isn't empty
		if widgetType == "" {
			n.SetState(gr.State{"error": "Widget type cannot be empty", "querying": false})
			return
		}

		n.SetState(gr.State{"error": "", "querying": false, "step": 2})
	}(n.State().String("widgetName"), n.State().String("widgetType"))

}

func (n NewWidget) stepTwoBack(event *gr.Event) {
	n.SetState(gr.State{"step": 1})
}
