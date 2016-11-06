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
	adjustmentTypes = []string{"ChangeInCapacity", "ExactCapacity", "PercentChangeInCapacity"}
)

type ScalingPolicyClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (s ScalingPolicyClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "", "step": 1}
}

// Implements the ComponentDidMount interface
func (s ScalingPolicyClassForm) ComponentWillMount() {
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

func (s ScalingPolicyClassForm) Render() gr.Component {

	state := s.State()
	props := s.Props()

	// Form placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)
	helpers.SuccessElem(state.String("success")).Modify(response)

	if state.Int("step") == 1 {
		if state.Bool("querying") {
			gr.Text("Loading...").Modify(response)
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

func (s ScalingPolicyClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

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

	numberField("Scaling Adjustment", "scalingAdjustment", state.Int("scalingAdjustment"), s.storeValue).Modify(classEditForm)
	selectOne("Adjustment Type", "adjustmentType", adjustmentTypes, state.Interface("adjustmentType"), s.storeSelect).Modify(classEditForm)
	numberField("Cooldown", "cooldown", state.Int("cooldown"), s.storeValue).Modify(classEditForm)
	selectMultiple("Alarms", "alarms", classOptions["keypairs"], state.Interface("alarms"), s.storeSelect)

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

func (s ScalingPolicyClassForm) backButton(*gr.Event) {
	s.SetState(gr.State{"success": ""})
	s.Props().Call("backButton")
}

func (s ScalingPolicyClassForm) doneButton(*gr.Event) {
	s.SetState(gr.State{"success": ""})
	s.Props().Call("hideAllModals")
}

func (s ScalingPolicyClassForm) saveButton(*gr.Event) {
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

func (s ScalingPolicyClassForm) deleteButton(*gr.Event) {
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

func (s ScalingPolicyClassForm) storeValue(event *gr.Event) {
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

func (s ScalingPolicyClassForm) storeSelect(id string, val interface{}) {
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
