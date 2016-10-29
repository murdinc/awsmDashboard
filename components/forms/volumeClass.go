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

var (
	volumeTypes = []string{"standard", "io1", "gp2", "sc1", "st1"}
)

type VolumeClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (v VolumeClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "", "step": 1}
}

// Implements the ComponentDidMount interface
func (v VolumeClassForm) ComponentWillMount() {
	var class map[string]interface{}

	if v.Props().Interface("class") != nil {
		classJson := v.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	v.SetState(class)
	v.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/classes/" + v.Props().String("apiType") + "/options"
		resp, err := helpers.GetAPI(endpoint)
		if !v.IsMounted() {
			return
		}
		if err != nil {
			v.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		v.SetState(gr.State{"classOptionsResp": resp, "querying": false})
	}()
}

func (v VolumeClassForm) Render() gr.Component {

	state := v.State()
	props := v.Props()

	// Form placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)
	helpers.SuccessElem(state.String("success")).Modify(response)

	if state.Int("step") == 1 {
		if state.Bool("querying") {
			gr.Text("Loading...").Modify(response)
		} else {
			v.BuildClassForm(props.String("className"), state.Interface("classOptionsResp")).Modify(response)
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
				evt.Click(v.backButton).PreventDefault(),
				gr.CSS("btn", "btn-secondary"),
				gr.Text("Back"),
			).Modify(buttons)

			// Done
			el.Button(
				evt.Click(v.doneButton).PreventDefault(),
				gr.CSS("btn", "btn-primary"),
				gr.Text("Done"),
			).Modify(buttons)

			buttons.Modify(response)
		}

	}

	return response
}

func (v VolumeClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

	state := v.State()
	props := v.Props()

	var classOptions map[string][]string
	jsonParsed, _ := gabs.ParseJSON(optionsResp.([]byte))
	classOptionsJson := jsonParsed.S("classOptions").Bytes()
	json.Unmarshal(classOptionsJson, &classOptions)

	classEdit := el.Div(
		el.Header3(gr.Text(className)),
		el.HorizontalRule(),
	)

	classEditForm := el.Form()

	textField("Device Name", "deviceName", &state, v.storeValue).Modify(classEditForm)
	textField("Volume Size", "volumeSize", &state, v.storeValue).Modify(classEditForm) // number
	checkbox("Delete On Termination", "deleteOnTermination", &state, v.storeValue).Modify(classEditForm)
	textField("Mount Mount", "mountPoint", &state, v.storeValue).Modify(classEditForm)
	selectOne("Snapshot", "snapshot", classOptions["snapshots"], &state, v.storeValue).Modify(classEditForm)
	selectOne("Volume Type", "volumeType", volumeTypes, &state, v.storeValue).Modify(classEditForm)
	if state.String("volumeType") == "io1" {
		textField("IOPS", "iops", &state, v.storeValue).Modify(classEditForm) // number
	}
	checkbox("Encrypted", "encrypted", &state, v.storeValue).Modify(classEditForm)

	classEditForm.Modify(classEdit)

	buttons := el.Div(
		gr.CSS("btn-toolbar"),
	)

	// Back
	el.Button(
		evt.Click(v.backButton).PreventDefault(),
		gr.CSS("btn", "btn-secondary"),
		gr.Text("Back"),
	).Modify(buttons)

	// Save
	el.Button(
		evt.Click(v.saveButton).PreventDefault(),
		gr.CSS("btn", "btn-primary"),
		gr.Text("Save"),
	).Modify(buttons)

	// Delete
	if props.Interface("hasDelete") != nil && props.Bool("hasDelete") {
		el.Button(
			evt.Click(v.deleteButton).PreventDefault(),
			gr.CSS("btn", "btn-danger", "pull-right"),
			gr.Text("Delete"),
		).Modify(buttons)
	}

	buttons.Modify(classEdit)

	return classEdit

}

func (v VolumeClassForm) backButton(*gr.Event) {
	v.SetState(gr.State{"success": ""})
	v.Props().Call("backButton")
}

func (v VolumeClassForm) doneButton(*gr.Event) {
	v.SetState(gr.State{"success": ""})
	v.Props().Call("hideAllModals")
}

func (v VolumeClassForm) saveButton(*gr.Event) {
	v.SetState(gr.State{"querying": true, "step": 2})

	cfg := make(map[string]interface{})
	for key, _ := range v.State() {
		cfg[key] = v.State().Interface(key)
	}

	go func() {
		endpoint := "//localhost:8081/api/classes/" + v.Props().String("apiType") + "/name/" + v.Props().String("className")

		_, err := helpers.PutAPI(endpoint, cfg)
		if !v.IsMounted() {
			return
		}

		if err != nil {
			v.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint), "step": 1})
			return
		}

		v.SetState(gr.State{"querying": false, "success": "Class was saved", "error": ""})
	}()

}

func (v VolumeClassForm) deleteButton(*gr.Event) {
	v.SetState(gr.State{"querying": true})

	go func() {
		endpoint := "//localhost:8081/api/classes/" + v.Props().String("apiType") + "/name/" + v.Props().String("className")

		_, err := helpers.DeleteAPI(endpoint)
		if !v.IsMounted() {
			return
		}

		if err != nil {
			v.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		v.SetState(gr.State{"querying": false, "success": "Class was deleted", "error": "", "step": 2})
	}()
}

func (v VolumeClassForm) storeValue(event *gr.Event) {
	id := event.Target().Get("id").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		v.SetState(gr.State{id: event.Target().Get("checked").Bool()})

	case "select-one":
		v.SetState(gr.State{id: event.TargetValue()})

	case "select-multiple":
		var vals []string
		options := event.Target().Length()

		for i := 0; i < options; i++ {
			if event.Target().Index(i).Get("selected").Bool() && event.Target().Index(i).Get("id") != nil {
				vals = append(vals, event.Target().Index(i).Get("id").String())
			}
		}
		v.SetState(gr.State{id: vals})

	default: // text, at least
		v.SetState(gr.State{id: event.TargetValue()})

	}
}
