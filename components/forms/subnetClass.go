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

type SubnetClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (s SubnetClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "", "step": 1}
}

// Implements the ComponentWillMount interface
func (s SubnetClassForm) ComponentWillMount() {
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

func (s SubnetClassForm) Render() gr.Component {

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

func (s SubnetClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

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

	classEditForm := el.Form(evt.KeyDown(DisableEnter))

	TextField("CIDR", "cidr", state.String("cidr"), s.storeValue).Modify(classEditForm)

	el.Div(
		el.Break(nil),
		el.Header4(
			gr.Text("Internet Gateway"),
		),
		el.HorizontalRule(nil),
	).Modify(classEditForm)

	el.Label(gr.Text("Create Internet Gateway?")).Modify(classEditForm)
	Toggle("No", "Yes", "createInternetGateway", state.Bool("createInternetGateway"), s.storeValue).Modify(classEditForm)

	if state.Bool("createInternetGateway") {
		el.Label(gr.Text("Add Internet Gateway to Main Route Table?")).Modify(classEditForm)
		Toggle("No", "Yes", "addInternetGatewayToMainRouteTable", state.Bool("addInternetGatewayToMainRouteTable"), s.storeValue).Modify(classEditForm)

		el.Label(gr.Text("Add Internet Gateway to New Route Table?")).Modify(classEditForm)
		Toggle("No", "Yes", "addInternetGatewayToNewRouteTable", state.Bool("addInternetGatewayToNewRouteTable"), s.storeValue).Modify(classEditForm)
	}

	el.Div(
		el.Break(nil),
		el.Header4(
			gr.Text("NAT Gateway"),
		),
		el.HorizontalRule(nil),
	).Modify(classEditForm)

	el.Label(gr.Text("Create NAT Gateway?")).Modify(classEditForm)
	Toggle("No", "Yes", "createNatGateway", state.Bool("createNatGateway"), s.storeValue).Modify(classEditForm)

	if state.Bool("createNatGateway") {

		el.Label(gr.Text("Add NAT Gateway to Main Route Table?")).Modify(classEditForm)
		Toggle("No", "Yes", "addNatGatewayToMainRouteTable", state.Bool("addNatGatewayToMainRouteTable"), s.storeValue).Modify(classEditForm)

		el.Label(gr.Text("Add NAT Gateway to New Route Table?")).Modify(classEditForm)
		Toggle("No", "Yes", "addNatGatewayToNewRouteTable", state.Bool("addNatGatewayToNewRouteTable"), s.storeValue).Modify(classEditForm)

	}
	el.HorizontalRule(nil).Modify(classEditForm)

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

func (s SubnetClassForm) backButton(*gr.Event) {
	s.SetState(gr.State{"success": ""})
	s.Props().Call("backButton")
}

func (s SubnetClassForm) doneButton(*gr.Event) {
	s.SetState(gr.State{"success": ""})
	s.Props().Call("hideAllModals")
}

func (s SubnetClassForm) saveButton(*gr.Event) {
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

func (s SubnetClassForm) deleteButton(*gr.Event) {
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

func (s SubnetClassForm) storeValue(event *gr.Event) {
	key := event.Target().Get("name").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		s.SetState(gr.State{key: event.Target().Get("checked").Bool()})

	case "number":
		s.SetState(gr.State{key: event.TargetValue().Int()})

	default: // text, at least
		s.SetState(gr.State{key: event.TargetValue()})

	}
}

func (s SubnetClassForm) storeSelect(key string, val interface{}) {
	switch value := val.(type) {

	case map[string]interface{}:
		// single
		s.SetState(gr.State{key: value["value"]})

	case []interface{}:
		// multi
		var vals []string
		options := len(value)
		for i := 0; i < options; i++ {
			vals = append(vals, value[i].(map[string]interface{})["value"].(string))
		}
		s.SetState(gr.State{key: vals})

	default:
		s.SetState(gr.State{key: val})

	}
}
