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

type SnapshotClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (s SnapshotClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "", "step": 1}
}

// Implements the ComponentDidMount interface
func (s SnapshotClassForm) ComponentWillMount() {
	var class map[string]interface{}

	if s.Props().Interface("class") != nil {
		classJson := s.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	s.SetState(class)
	s.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/classes/" + s.Props().String("apiType") + "/options"
		resp, err := helpers.GetAPI(endpoint)
		if !s.IsMounted() {
			return
		}
		if err != nil {
			s.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		s.SetState(gr.State{"classOptionsResp": resp, "querying": false})
	}()
}

func (s SnapshotClassForm) Render() gr.Component {

	state := s.State()
	props := s.Props()

	// Form placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)
	helpers.SuccessElem(state.String("success")).Modify(response)

	if state.Int("step") == 1 {
		if state.Bool("querying") {
			gr.Text("Loading.......").Modify(response)
		} else {
			s.BuildClassForm(props.String("className"), state.Interface("classOptionsResp")).Modify(response)
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
				evt.Click(s.backButton).PreventDefault(),
				gr.CSS("btn", "btn-secondary"),
				gr.Text("Back"),
			).Modify(buttons)

			// Done
			el.Button(
				evt.Click(s.doneButton).PreventDefault(),
				gr.CSS("btn", "btn-primary"),
				gr.Text("Done"),
			).Modify(buttons)

			buttons.Modify(response)
		}

	}

	return response
}

func (s SnapshotClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

	state := s.State()
	props := s.Props()

	var classOptions map[string][]string
	jsonParsed, _ := gabs.ParseJSON(optionsResp.([]byte))
	classOptionsJson := jsonParsed.S("classOptions").Bytes()
	json.Unmarshal(classOptionsJson, &classOptions)

	classEdit := el.Div(
		el.Header3(gr.Text(className)),
		el.HorizontalRule(),
	)

	classEditForm := el.Form()

	checkbox("Rotate", "rotate", &state, s.storeValue).Modify(classEditForm)
	if state.Bool("rotate") {
		textField("Retain", "retain", &state, s.storeValue).Modify(classEditForm) // number
	}
	checkbox("Propagate", "propagate", &state, s.storeValue).Modify(classEditForm)
	if state.Bool("propagate") {
		selectMultiple("Propagate Regions", "propagateRegions", classOptions["regions"], &state, s.storeValue).Modify(classEditForm)
	}

	textField("Volume ID", "volumeID", &state, s.storeValue).Modify(classEditForm) // select one?

	classEditForm.Modify(classEdit)

	buttons := el.Div(
		gr.CSS("btn-toolbar"),
	)

	// Back
	el.Button(
		evt.Click(s.backButton).PreventDefault(),
		gr.CSS("btn", "btn-secondary"),
		gr.Text("Back"),
	).Modify(buttons)

	// Save
	el.Button(
		evt.Click(s.saveButton).PreventDefault(),
		gr.CSS("btn", "btn-primary"),
		gr.Text("Save"),
	).Modify(buttons)

	// Delete
	if props.Interface("hasDelete") != nil && props.Bool("hasDelete") {
		el.Button(
			evt.Click(s.deleteButton).PreventDefault(),
			gr.CSS("btn", "btn-danger", "pull-right"),
			gr.Text("Delete"),
		).Modify(buttons)
	}

	buttons.Modify(classEdit)

	return classEdit

}

func (s SnapshotClassForm) backButton(*gr.Event) {
	s.SetState(gr.State{"success": ""})
	s.Props().Call("backButton")
}

func (s SnapshotClassForm) doneButton(*gr.Event) {
	s.SetState(gr.State{"success": ""})
	s.Props().Call("hideAllModals")
}

func (s SnapshotClassForm) saveButton(*gr.Event) {
	s.SetState(gr.State{"querying": true, "step": 2})

	cfg := make(map[string]interface{})
	for key, _ := range s.State() {
		cfg[key] = s.State().Interface(key)
	}

	go func() {
		endpoint := "//localhost:8081/api/classes/" + s.Props().String("apiType") + "/name/" + s.Props().String("className")

		_, err := helpers.PutAPI(endpoint, cfg)
		if !s.IsMounted() {
			return
		}

		if err != nil {
			s.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint), "step": 1})
			return
		}

		s.SetState(gr.State{"querying": false, "success": "Class was saved", "error": ""})
	}()

}

func (s SnapshotClassForm) deleteButton(*gr.Event) {
	s.SetState(gr.State{"querying": true})

	go func() {
		endpoint := "//localhost:8081/api/classes/" + s.Props().String("apiType") + "/name/" + s.Props().String("className")

		_, err := helpers.DeleteAPI(endpoint)
		if !s.IsMounted() {
			return
		}

		if err != nil {
			s.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		s.SetState(gr.State{"querying": false, "success": "Class was deleted", "error": "", "step": 2})
	}()
}

func (s SnapshotClassForm) storeValue(event *gr.Event) {
	id := event.Target().Get("id").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		s.SetState(gr.State{id: event.Target().Get("checked").Bool()})

	case "select-one":
		s.SetState(gr.State{id: event.TargetValue()})

	case "select-multiple":
		var vals []string
		options := event.Target().Length()

		for i := 0; i < options; i++ {
			if event.Target().Index(i).Get("selected").Bool() && event.Target().Index(i).Get("id") != nil {
				vals = append(vals, event.Target().Index(i).Get("id").String())
			}
		}
		s.SetState(gr.State{id: vals})

	default: // text, at least
		s.SetState(gr.State{id: event.TargetValue()})

	}
}
