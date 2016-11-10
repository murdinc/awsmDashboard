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
	tenancy = []string{"default", "dedicated"}
)

type VpcClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (v VpcClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "", "step": 1,
		"tenancy": "default",
	}
}

// Implements the ComponentWillMount interface
func (v VpcClassForm) ComponentWillMount() {
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

func (v VpcClassForm) Render() gr.Component {

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

func (v VpcClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

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

	textField("CIDR", "cidr", state.String("cidr"), v.storeValue).Modify(classEditForm)
	selectOne("Tenancy", "tenancy", tenancy, state.Interface("tenancy"), v.storeSelect).Modify(classEditForm)

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

func (v VpcClassForm) backButton(*gr.Event) {
	v.SetState(gr.State{"success": ""})
	v.Props().Call("backButton")
}

func (v VpcClassForm) doneButton(*gr.Event) {
	v.SetState(gr.State{"success": ""})
	v.Props().Call("hideAllModals")
}

func (v VpcClassForm) saveButton(*gr.Event) {
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

func (v VpcClassForm) deleteButton(*gr.Event) {
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

func (v VpcClassForm) storeValue(event *gr.Event) {
	id := event.Target().Get("id").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		v.SetState(gr.State{id: event.Target().Get("checked").Bool()})

	case "number":
		v.SetState(gr.State{id: event.TargetValue().Int()})

	default: // text, at least
		v.SetState(gr.State{id: event.TargetValue()})

	}
}

func (v VpcClassForm) storeSelect(id string, val interface{}) {
	switch value := val.(type) {

	case map[string]interface{}:
		// single
		v.SetState(gr.State{id: value["value"]})

	case []interface{}:
		// multi
		var vals []string
		options := len(value)
		for i := 0; i < options; i++ {
			vals = append(vals, value[i].(map[string]interface{})["value"].(string))
		}
		v.SetState(gr.State{id: vals})

	default:
		v.SetState(gr.State{id: val})

	}
}
