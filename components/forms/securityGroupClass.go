package forms

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
	"github.com/murdinc/awsm/config"
	"github.com/murdinc/awsmDashboard/helpers"
)

var ()

type SecurityGroupClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (s SecurityGroupClassForm) GetInitialState() gr.State {
	return gr.State{"querying": false, "error": "", "success": "", "step": 1}
}

// Implements the ComponentWillMount interface
func (s SecurityGroupClassForm) ComponentWillMount() {
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

func (s SecurityGroupClassForm) addGrant(*gr.Event) {
	grants, ok := s.State().Interface("securityGroupGrants").([]interface{})
	if ok {
		grants = append([]interface{}{
			&config.SecurityGroupGrant{
				Note:       "New Grant",
				Type:       "ingress",
				FromPort:   0,
				ToPort:     0,
				IPProtocol: "tcp",
				CidrIP:     []string{},
			},
		}, grants...)
		s.SetState(gr.State{"securityGroupGrants": grants})
		return
	}
	println("addGrant failed?")
}

func (s SecurityGroupClassForm) removeGrant(event *gr.Event) {
	index := event.Target().Get("id").Int()
	println("Deleting: ")
	println(index)
	grants, ok := s.State().Interface("securityGroupGrants").([]interface{})
	if ok {
		grants = append(grants[:index], grants[index+1:]...)
		s.SetState(gr.State{"securityGroupGrants": grants})
		return
	}
	println("removeGrant failed?")
}

func (s SecurityGroupClassForm) Render() gr.Component {

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

func (s SecurityGroupClassForm) modifyGrant(index int, grant map[string]interface{}) func(*gr.Event) {
	return func(event *gr.Event) {
		key := event.Target().Get("id").String()
		valueType := event.Target().Get("type").String()

		switch valueType {
		case "text":
			grant[key] = event.TargetValue().String()
		case "number":
			grant[key] = event.TargetValue().Int()
		default:
			println("modifyGrant does not have a switch for type:")
			println(valueType)
		}

		grants, ok := s.State().Interface("securityGroupGrants").([]interface{})
		if ok {
			grants[index] = grant
			s.SetState(gr.State{"securityGroupGrants": grants})
			return
		}
		println("modifyGrant failed?")
	}
}

func (s SecurityGroupClassForm) storeSelect(index int, grant map[string]interface{}) func(string, interface{}) {
	return func(id string, val interface{}) {

		switch value := val.(type) {

		case map[string]interface{}:
			// single
			grant[id] = value["value"]

		case []interface{}:
			// multi
			var vals []string
			options := len(value)
			for i := 0; i < options; i++ {
				vals = append(vals, value[i].(map[string]interface{})["value"].(string))
			}
			grant[id] = vals

		default:
			grant[id] = val
		}

		grants, ok := s.State().Interface("securityGroupGrants").([]interface{})
		if ok {
			grants[index] = grant
			s.SetState(gr.State{"securityGroupGrants": grants})
			return
		}
	}
}

func (s SecurityGroupClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

	state := s.State()
	props := s.Props()

	classEdit := el.Div(
		el.Header3(gr.Text(className)),
		el.HorizontalRule(),
	)

	classEditForm := el.Form()

	TextField("Description", "description", state.String("description"), s.storeValue).Modify(classEditForm)

	el.Div(
		el.Break(nil),
		el.Header4(
			gr.Text("Grants"),
			el.Button(
				evt.Click(s.addGrant).PreventDefault(),
				gr.CSS("btn", "btn-primary", "btn-sm", "pull-right"),
				gr.Text("New"),
			),
		),
		el.HorizontalRule(nil),
	).Modify(classEditForm)

	optionsRespByte, ok := optionsResp.([]byte)
	if ok {

		var classOptions map[string][]string

		jsonParsed, _ := gabs.ParseJSON(optionsRespByte)
		classOptionsJson := jsonParsed.S("classOptions").Bytes()
		json.Unmarshal(classOptionsJson, &classOptions)

		grants := state.Interface("securityGroupGrants").([]interface{})

		for index, g := range grants {

			grant := g.(map[string]interface{})

			// Form placeholder
			grantForm := el.Div()

			TextField("Note", "note", grant["note"], s.modifyGrant(index, grant)).Modify(grantForm)

			el.Div(
				gr.CSS("row"), el.Div(gr.CSS("col-sm-3"),
					SelectOne("Type", "type", []string{"ingress", "egress"}, grant["type"], s.storeSelect(index, grant)),
				),
				el.Div(gr.CSS("col-sm-3"),
					NumberField("From Port", "fromPort", grant["fromPort"], s.modifyGrant(index, grant)),
				),
				el.Div(gr.CSS("col-sm-3"),
					NumberField("To Port", "toPort", grant["toPort"], s.modifyGrant(index, grant)),
				),
				el.Div(gr.CSS("col-sm-3"),
					SelectOne("IP Protocol", "ipProtocol", []string{"tcp", "udp", "icmp"}, grant["ipProtocol"], s.storeSelect(index, grant)),
				),
			).Modify(grantForm)

			el.Div(
				gr.CSS("row"),
				el.Div(gr.CSS("col-sm-12"),
					CreateableSelectMultiple("CIDR IP's", "cidrIP", []string{ /* TODO */ }, grant["cidrIP"], s.storeSelect(index, grant)),
				),
			).Modify(grantForm)

			el.Div(
				gr.CSS("btn-toolbar"),
				el.Button(
					evt.Click(s.removeGrant).PreventDefault(),
					gr.CSS("btn", "btn-danger", "btn-sm", "pull-right"),
					gr.Text("Remove"),
					attr.ID(index),
				),
			).Modify(grantForm)

			el.HorizontalRule(nil).Modify(grantForm)

			grantForm.Modify(classEditForm)

		}

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

	}

	return classEdit

}

func (s SecurityGroupClassForm) backButton(*gr.Event) {
	s.SetState(gr.State{"success": ""})
	s.Props().Call("backButton")
}

func (s SecurityGroupClassForm) doneButton(*gr.Event) {
	s.SetState(gr.State{"success": ""})
	s.Props().Call("hideAllModals")
}

func (s SecurityGroupClassForm) saveButton(*gr.Event) {
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

func (s SecurityGroupClassForm) deleteButton(*gr.Event) {
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

func (s SecurityGroupClassForm) storeValue(event *gr.Event) {
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
