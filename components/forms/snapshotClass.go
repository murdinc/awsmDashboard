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
	return gr.State{"querying": true, "queryingOpts": true, "queryingVolumes": true, "error": "", "success": "",
		"rotate":    false,
		"propagate": false,
		"step":      1,
	}
}

// Implements the ComponentWillMount interface
func (s SnapshotClassForm) ComponentWillMount() {
	var class map[string]interface{}

	if s.Props().Interface("class") != nil {
		classJson := s.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	s.SetState(class)
	s.SetState(gr.State{"queryingOpts": true, "queryingVolumes": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/classes/" + s.Props().String("apiType") + "/options"
		resp, err := helpers.GetAPI(endpoint)
		if !s.IsMounted() {
			return
		}
		if err != nil {
			s.SetState(gr.State{"queryingOpts": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		s.SetState(gr.State{"classOptionsResp": resp, "queryingOpts": false})
	}()

	// Get our existing instances for the form
	go func() {
		endpoint := "//localhost:8081/api/assets/volumes"
		resp, err := helpers.GetAPI(endpoint)
		if !s.IsMounted() {
			return
		}
		if err != nil {
			s.SetState(gr.State{"queryingVolumes": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		s.SetState(gr.State{"volumeOptionsResp": resp, "queryingVolumes": false})
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
		if state.Bool("queryingOpts") || state.Bool("queryingVolumes") {
			gr.Text("Loading...").Modify(response)
		} else {
			s.BuildClassForm(props.String("className"), state.Interface("classOptionsResp"), state.Interface("volumeOptionsResp")).Modify(response)
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

func (s SnapshotClassForm) BuildClassForm(className string, optionsResp interface{}, volumeResp interface{}) *gr.Element {

	state := s.State()
	props := s.Props()

	var classOptions map[string][]string
	jsonParsed, _ := gabs.ParseJSON(optionsResp.([]byte))
	classOptionsJson := jsonParsed.S("classOptions").Bytes()
	json.Unmarshal(classOptionsJson, &classOptions)

	volumeJsonParsed, _ := gabs.ParseJSON(volumeResp.([]byte))
	volumeOptionsSlice, _ := volumeJsonParsed.S("assets").Children()

	var volumes []string
	volumesMeta := make(map[string]string)
	for _, volumeOption := range volumeOptionsSlice {
		volume := volumeOption.S("volumeID").Data().(string)
		if volume != "" {
			volumes = append(volumes, volume)
			volumesMeta[volume] = volumeOption.S("sizeHuman").Data().(string) + " " + volumeOption.S("name").Data().(string) + " " + volumeOption.S("attachment").Data().(string)
		}
	}

	classEdit := el.Div(
		el.Header3(gr.Text(className)),
		el.HorizontalRule(),
	)

	classEditForm := el.Form()

	TextField("Description", "description", state.String("description"), s.storeValue).Modify(classEditForm)

	Checkbox("Rotate", "rotate", state.Bool("rotate"), s.storeValue).Modify(classEditForm)
	if state.Bool("rotate") {
		NumberField("Retain", "retain", state.Int("retain"), s.storeValue).Modify(classEditForm)
	}
	Checkbox("Propagate", "propagate", state.Bool("propagate"), s.storeValue).Modify(classEditForm)
	if state.Bool("propagate") {
		SelectMultiple("Propagate Regions", "propagateRegions", classOptions["regions"], state.Interface("propagateRegions"), s.storeSelect).Modify(classEditForm)
	}

	SelectOneMeta("Volume", "volume", volumes, volumesMeta, state.Interface("volume"), s.storeSelect).Modify(classEditForm)

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

	case "number":
		s.SetState(gr.State{id: event.TargetValue().Int()})

	default: // text, at least
		s.SetState(gr.State{id: event.TargetValue()})

	}
}

func (s SnapshotClassForm) storeSelect(id string, val interface{}) {
	switch value := val.(type) {

	case map[string]interface{}:
		// single
		s.SetState(gr.State{id: value["value"]})

	case []interface{}:
		// multi
		var vals []string
		options := len(value)
		for i := 0; i < options; i++ {
			vals = append(vals, value[i].(map[string]interface{})["value"].(string))
		}
		s.SetState(gr.State{id: vals})

	default:
		s.SetState(gr.State{id: val})

	}
}
