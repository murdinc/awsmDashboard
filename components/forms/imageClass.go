package forms

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
	"github.com/murdinc/awsmDashboard/helpers"
)

type ImageClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (i ImageClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "", "step": 1}
}

// Implements the ComponentDidMount interface
func (i ImageClassForm) ComponentWillMount() {
	var class map[string]interface{}

	if i.Props().Interface("class") != nil {
		classJson := i.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	i.SetState(class)
	i.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/classes/" + i.Props().String("apiType") + "/options"
		resp, err := helpers.GetAPI(endpoint)
		if !i.IsMounted() {
			return
		}
		if err != nil {
			i.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		i.SetState(gr.State{"classOptionsResp": resp, "querying": false})
	}()
}

func (i ImageClassForm) Render() gr.Component {

	state := i.State()
	props := i.Props()

	// Form placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)
	helpers.SuccessElem(state.String("success")).Modify(response)

	if state.Int("step") == 1 {
		if state.Bool("querying") {
			gr.Text("Loading.......").Modify(response)
		} else {
			i.BuildClassForm(props.String("className"), state.Interface("classOptionsResp")).Modify(response)
		}

	} else if state.Int("step") == 2 {

		if state.Bool("querying") {
			gr.Text("Saving...").Modify(response)
		} else {

			buttons := el.Div(
				gr.CSS("btn-toolbar"),
			)

			// Back
			el.Button(
				evt.Click(i.backButton).PreventDefault(),
				gr.CSS("btn", "btn-secondary"),
				gr.Text("Back"),
			).Modify(buttons)

			// Done
			el.Button(
				evt.Click(i.doneButton).PreventDefault(),
				gr.CSS("btn", "btn-primary"),
				gr.Text("Done"),
			).Modify(buttons)

			buttons.Modify(response)
		}

	}

	return response
}

func (i ImageClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

	state := i.State()
	props := i.Props()

	var classOptions map[string][]string
	jsonParsed, _ := gabs.ParseJSON(optionsResp.([]byte))
	classOptionsJson := jsonParsed.S("classOptions").Bytes()
	json.Unmarshal(classOptionsJson, &classOptions)

	classEdit := el.Div(
		el.Header3(gr.Text(className)),
		el.HorizontalRule(),
	)

	classEditForm := el.Form()

	checkbox("Propagate", "propagate", &state, i.storeValue).Modify(classEditForm)
	if state.Bool("propagate") {
		selectMultiple("Propagate Regions", "propagateRegions", classOptions["regions"], &state, i.storeValue).Modify(classEditForm)
	}
	checkbox("Rotate", "rotate", &state, i.storeValue).Modify(classEditForm)
	if state.Bool("rotate") {
		textField("Retain", "retain", &state, i.storeValue).Modify(classEditForm) // number
	}

	textField("Instance ID", "instanceID", &state, i.storeValue).Modify(classEditForm) // select one?

	classEditForm.Modify(classEdit)

	buttons := el.Div(
		gr.CSS("btn-toolbar"),
	)

	// Back
	el.Button(
		evt.Click(i.backButton).PreventDefault(),
		gr.CSS("btn", "btn-secondary"),
		gr.Text("Back"),
	).Modify(buttons)

	// Save
	el.Button(
		evt.Click(i.saveButton).PreventDefault(),
		gr.CSS("btn", "btn-primary"),
		gr.Text("Save"),
	).Modify(buttons)

	// Delete
	if props.Interface("hasDelete") != nil && props.Bool("hasDelete") {
		el.Button(
			evt.Click(i.deleteButton).PreventDefault(),
			gr.CSS("btn", "btn-danger", "pull-right"),
			gr.Text("Delete"),
		).Modify(buttons)
	}

	buttons.Modify(classEdit)

	return classEdit

}

func (i ImageClassForm) backButton(*gr.Event) {
	i.SetState(gr.State{"success": ""})
	i.Props().Call("backButton")
}

func (i ImageClassForm) doneButton(*gr.Event) {
	i.SetState(gr.State{"success": ""})
	i.Props().Call("hideAllModals")
}

func (i ImageClassForm) saveButton(*gr.Event) {
	i.SetState(gr.State{"querying": true, "step": 2})

	cfg := make(map[string]interface{})
	for key, _ := range i.State() {
		cfg[key] = i.State().Interface(key)
	}

	go func() {
		endpoint := "//localhost:8081/api/classes/" + i.Props().String("apiType") + "/name/" + i.Props().String("className")

		_, err := helpers.PutAPI(endpoint, cfg)
		if !i.IsMounted() {
			return
		}

		if err != nil {
			i.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint), "step": 1})
			return
		}

		i.SetState(gr.State{"querying": false, "success": "Class was saved", "error": ""})
	}()

}

func (i ImageClassForm) deleteButton(*gr.Event) {
	i.SetState(gr.State{"querying": true})

	go func() {
		endpoint := "//localhost:8081/api/classes/" + i.Props().String("apiType") + "/name/" + i.Props().String("className")

		_, err := helpers.DeleteAPI(endpoint)
		if !i.IsMounted() {
			return
		}

		if err != nil {
			i.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		i.SetState(gr.State{"querying": false, "success": "Class was deleted", "error": "", "step": 2})
	}()
}

func (i ImageClassForm) storeValue(event *gr.Event) {
	id := event.Target().Get("id").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		i.SetState(gr.State{id: event.Target().Get("checked").Bool()})

	case "select-one":
		i.SetState(gr.State{id: event.TargetValue()})

	case "select-multiple":
		var vals []string
		options := event.Target().Length()

		for i := 0; i < options; i++ {
			if event.Target().Index(i).Get("selected").Bool() && event.Target().Index(i).Get("id") != nil {
				vals = append(vals, event.Target().Index(i).Get("id").String())
			}
		}
		i.SetState(gr.State{id: vals})

	default: // text, at least
		i.SetState(gr.State{id: event.TargetValue()})

	}
}
