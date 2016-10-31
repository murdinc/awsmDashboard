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
	elbSchemes = []string{"internal", "internet-facing"}
)

type LoadBalancerClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (l LoadBalancerClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "", "step": 1}
}

// Implements the ComponentDidMount interface
func (l LoadBalancerClassForm) ComponentWillMount() {
	var class map[string]interface{}

	if l.Props().Interface("class") != nil {
		classJson := l.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	l.SetState(class)
	l.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/classes/" + l.Props().String("apiType") + "/options"
		resp, err := helpers.GetAPI(endpoint)
		if !l.IsMounted() {
			return
		}
		if err != nil {
			l.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		l.SetState(gr.State{"classOptionsResp": resp, "querying": false})
	}()
}

func (l LoadBalancerClassForm) Render() gr.Component {

	state := l.State()
	props := l.Props()

	// Form placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)
	helpers.SuccessElem(state.String("success")).Modify(response)

	if state.Int("step") == 1 {
		if state.Bool("querying") {
			gr.Text("Loading...").Modify(response)
		} else {
			l.BuildClassForm(props.String("className"), state.Interface("classOptionsResp")).Modify(response)
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
				evt.Click(l.backButton).PreventDefault(),
				gr.CSS("btn", "btn-secondary"),
				gr.Text("Back"),
			).Modify(buttons)

			// Done
			el.Button(
				evt.Click(l.doneButton).PreventDefault(),
				gr.CSS("btn", "btn-primary"),
				gr.Text("Done"),
			).Modify(buttons)

			buttons.Modify(response)
		}

	}

	return response
}

func (l LoadBalancerClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

	state := l.State()
	props := l.Props()

	var classOptions map[string][]string
	jsonParsed, _ := gabs.ParseJSON(optionsResp.([]byte))
	classOptionsJson := jsonParsed.S("classOptions").Bytes()
	json.Unmarshal(classOptionsJson, &classOptions)

	classEdit := el.Div(
		el.Header3(gr.Text(className)),
		el.HorizontalRule(),
	)

	classEditForm := el.Form()

	selectOne("Scheme", "scheme", elbSchemes, &state, l.storeSelect).Modify(classEditForm)
	selectMultiple("Security Groups", "securityGroups", classOptions["securitygroups"], &state, l.storeSelect).Modify(classEditForm)
	selectMultiple("Subnets", "subnets", classOptions["subnets"], &state, l.storeSelect).Modify(classEditForm)
	selectMultiple("Availability Zones", "availabilityZones", classOptions["zones"], &state, l.storeSelect).Modify(classEditForm)
	// TODO Listeners

	classEditForm.Modify(classEdit)

	buttons := el.Div(
		gr.CSS("btn-toolbar"),
	)

	// Back
	el.Button(
		evt.Click(l.backButton).PreventDefault(),
		gr.CSS("btn", "btn-secondary"),
		gr.Text("Back"),
	).Modify(buttons)

	// Save
	el.Button(
		evt.Click(l.saveButton).PreventDefault(),
		gr.CSS("btn", "btn-primary"),
		gr.Text("Save"),
	).Modify(buttons)

	// Delete
	if props.Interface("hasDelete") != nil && props.Bool("hasDelete") {
		el.Button(
			evt.Click(l.deleteButton).PreventDefault(),
			gr.CSS("btn", "btn-danger", "pull-right"),
			gr.Text("Delete"),
		).Modify(buttons)
	}

	buttons.Modify(classEdit)

	return classEdit

}

func (l LoadBalancerClassForm) backButton(*gr.Event) {
	l.SetState(gr.State{"success": ""})
	l.Props().Call("backButton")
}

func (l LoadBalancerClassForm) doneButton(*gr.Event) {
	l.SetState(gr.State{"success": ""})
	l.Props().Call("hideAllModals")
}

func (l LoadBalancerClassForm) saveButton(*gr.Event) {
	l.SetState(gr.State{"querying": true, "step": 2})

	cfg := make(map[string]interface{})
	for key, _ := range l.State() {
		cfg[key] = l.State().Interface(key)
	}

	go func() {
		endpoint := "//localhost:8081/api/classes/" + l.Props().String("apiType") + "/name/" + l.Props().String("className")

		_, err := helpers.PutAPI(endpoint, cfg)
		if !l.IsMounted() {
			return
		}

		if err != nil {
			l.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint), "step": 1})
			return
		}

		l.SetState(gr.State{"querying": false, "success": "Class was saved", "error": ""})
	}()

}

func (l LoadBalancerClassForm) deleteButton(*gr.Event) {
	l.SetState(gr.State{"querying": true})

	go func() {
		endpoint := "//localhost:8081/api/classes/" + l.Props().String("apiType") + "/name/" + l.Props().String("className")

		_, err := helpers.DeleteAPI(endpoint)
		if !l.IsMounted() {
			return
		}

		if err != nil {
			l.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		l.SetState(gr.State{"querying": false, "success": "Class was deleted", "error": "", "step": 2})
	}()
}

func (l LoadBalancerClassForm) storeValue(event *gr.Event) {
	id := event.Target().Get("id").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		l.SetState(gr.State{id: event.Target().Get("checked").Bool()})

	case "number":
		l.SetState(gr.State{id: event.TargetValue().Int()})

	default: // text, at least
		l.SetState(gr.State{id: event.TargetValue()})

	}
}

func (l LoadBalancerClassForm) storeSelect(id string, val interface{}) {
	switch value := val.(type) {

	case map[string]interface{}:
		// single
		l.SetState(gr.State{id: value["value"]})

	case []interface{}:
		// multi
		var vals []string
		options := len(value)
		for i := 0; i < options; i++ {
			vals = append(vals, value[i].(map[string]interface{})["value"].(string))
		}
		l.SetState(gr.State{id: vals})

	default:
		l.SetState(gr.State{id: val})

	}
}
