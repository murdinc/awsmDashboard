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

var ()

type KeyPairClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (k KeyPairClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "", "step": 1}
}

// Implements the ComponentWillMount interface
func (k KeyPairClassForm) ComponentWillMount() {
	var class map[string]interface{}

	if k.Props().Interface("class") != nil {
		classJson := k.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	k.SetState(class)
	k.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/classes/" + k.Props().String("apiType") + "/options"
		resp, err := helpers.GetAPI(endpoint)
		if !k.IsMounted() {
			return
		}
		if err != nil {
			k.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		k.SetState(gr.State{"classOptionsResp": resp, "querying": false})
	}()
}

func (k KeyPairClassForm) Render() gr.Component {

	state := k.State()
	props := k.Props()

	// Form placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)
	helpers.SuccessElem(state.String("success")).Modify(response)

	if state.Int("step") == 1 {
		if state.Bool("querying") {
			gr.Text("Loading...").Modify(response)
		} else {
			k.BuildClassForm(props.String("className"), state.Interface("classOptionsResp")).Modify(response)
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
				evt.Click(k.backButton).PreventDefault(),
				gr.CSS("btn", "btn-secondary"),
				gr.Text("Back"),
			).Modify(buttons)

			// Done
			el.Button(
				evt.Click(k.doneButton).PreventDefault(),
				gr.CSS("btn", "btn-primary"),
				gr.Text("Done"),
			).Modify(buttons)

			buttons.Modify(response)
		}

	}

	return response
}

func (k KeyPairClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

	state := k.State()
	props := k.Props()

	var classOptions map[string][]string
	jsonParsed, _ := gabs.ParseJSON(optionsResp.([]byte))
	classOptionsJson := jsonParsed.S("classOptions").Bytes()
	json.Unmarshal(classOptionsJson, &classOptions)

	classEdit := el.Div(
		el.Header3(gr.Text(className)),
		el.HorizontalRule(),
	)

	classEditForm := el.Form()

	TextField("Description", "description", state.String("description"), k.storeValue).Modify(classEditForm)

	if props.Interface("newClass") != nil && props.Bool("newClass") {
		el.Div(
			el.Break(nil),
			gr.Text("Leave Key fields blank to generate a new key"),
			el.HorizontalRule(nil),
		).Modify(classEditForm)
	}

	// TODO make these file uploads
	TextArea("Public Key", "publicKey", state.String("publicKey"), k.storeValue).Modify(classEditForm)

	//if props.Interface("newClass") != nil && props.Bool("newClass") {
	TextArea("Private Key", "privateKey", state.String("privateKey"), k.storeValue).Modify(classEditForm)
	//}

	classEditForm.Modify(classEdit)

	buttons := el.Div(
		gr.CSS("btn-toolbar"),
	)

	// Back
	el.Button(
		evt.Click(k.backButton).PreventDefault(),
		gr.CSS("btn", "btn-secondary"),
		gr.Text("Back"),
	).Modify(buttons)

	// Save
	el.Button(
		evt.Click(k.saveButton).PreventDefault(),
		gr.CSS("btn", "btn-primary"),
		gr.Text("Save"),
	).Modify(buttons)

	// Delete
	if props.Interface("hasDelete") != nil && props.Bool("hasDelete") {
		el.Button(
			evt.Click(k.deleteButton).PreventDefault(),
			gr.CSS("btn", "btn-danger", "pull-right"),
			gr.Text("Delete"),
		).Modify(buttons)
	}

	buttons.Modify(classEdit)

	return classEdit

}

func (k KeyPairClassForm) backButton(*gr.Event) {
	k.SetState(gr.State{"success": ""})
	k.Props().Call("backButton")
}

func (k KeyPairClassForm) doneButton(*gr.Event) {
	k.SetState(gr.State{"success": ""})
	k.Props().Call("hideAllModals")
}

func (k KeyPairClassForm) saveButton(*gr.Event) {
	k.SetState(gr.State{"querying": true, "step": 2, "error": ""})

	props := k.Props()
	if props.Interface("newClass") != nil && props.Bool("newClass") {
		// Validate key input
		pubKey := k.State().String("privateKey")
		privKey := k.State().String("publicKey")

		if pubKey != "" && privKey == "" {
			k.SetState(gr.State{"querying": false, "error": "Please enter a Public Key, or leave the Private Key blank.", "step": 1})
			return
		}

		if pubKey == "" && privKey != "" {
			k.SetState(gr.State{"querying": false, "error": "Please enter a Private Key, or leave the Public Key blank.", "step": 1})
			return
		}
	}

	cfg := make(map[string]interface{})
	for key, _ := range k.State() {
		cfg[key] = k.State().Interface(key)
	}

	go func() {
		endpoint := "//localhost:8081/api/classes/" + k.Props().String("apiType") + "/name/" + k.Props().String("className")

		_, err := helpers.PutAPI(endpoint, cfg)
		if !k.IsMounted() {
			return
		}

		if err != nil {
			k.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint), "step": 1})
			return
		}

		k.SetState(gr.State{"querying": false, "success": "Class was saved", "error": ""})
	}()

}

func (k KeyPairClassForm) deleteButton(*gr.Event) {
	k.SetState(gr.State{"querying": true})

	go func() {
		endpoint := "//localhost:8081/api/classes/" + k.Props().String("apiType") + "/name/" + k.Props().String("className")

		_, err := helpers.DeleteAPI(endpoint)
		if !k.IsMounted() {
			return
		}

		if err != nil {
			k.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		k.SetState(gr.State{"querying": false, "success": "Class was deleted", "error": "", "step": 2})
	}()
}

func (k KeyPairClassForm) storeValue(event *gr.Event) {
	id := event.Target().Get("id").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		k.SetState(gr.State{id: event.Target().Get("checked").Bool()})

	case "number":
		k.SetState(gr.State{id: event.TargetValue().Int()})

	default: // text, at least
		k.SetState(gr.State{id: event.TargetValue()})

	}
}

func (k KeyPairClassForm) storeSelect(id string, val interface{}) {
	switch value := val.(type) {

	case map[string]interface{}:
		// single
		k.SetState(gr.State{id: value["value"]})

	case []interface{}:
		// multi
		var vals []string
		options := len(value)
		for i := 0; i < options; i++ {
			vals = append(vals, value[i].(map[string]interface{})["value"].(string))
		}
		k.SetState(gr.State{id: vals})

	default:
		k.SetState(gr.State{id: val})

	}
}
