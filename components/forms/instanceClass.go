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
	instanceTypes = []string{
		"t2.nano", "t2.micro", "t2.small", "t2.medium", "t2.large", "m4.large", "m4.xlarge", "m4.2xlarge", "m4.4xlarge", "m4.10xlarge", "m4.16xlarge", "m3.medium",
		"m3.large", "m3.xlarge", "m3.2xlarge", "c4.large", "c4.xlarge", "c4.2xlarge", "c4.4xlarge", "c4.8xlarge", "c3.large", "c3.xlarge", "c3.2xlarge", "c3.4xlarge",
		"c3.8xlarge", "r3.large", "r3.xlarge", "r3.2xlarge", "r3.4xlarge", "r3.8xlarge", "x1.16xlarge", "x1.32xlarge", "i2.xlarge", "i2.2xlarge", "i2.4xlarge",
		"i2.8xlarge", "d2.xlarge", "d2.2xlarge", "d2.4xlarge", "d2.8xlarge", "p2.xlarge", "p2.8xlarge", "p2.16xlarge", "g2.2xlarge", "g2.8xlarge",
	}

	shutdownBehaviors = []string{"stop", "terminate"}
)

type InstanceClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (i InstanceClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "",
		"monitoring":      false,
		"publicIpAddress": false,
		"ebsOptimized":    false,
		"step":            1,
	}
}

// Implements the ComponentDidMount interface
func (i InstanceClassForm) ComponentWillMount() {
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

func (i InstanceClassForm) Render() gr.Component {

	state := i.State()
	props := i.Props()

	// Form placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)
	helpers.SuccessElem(state.String("success")).Modify(response)

	if state.Int("step") == 1 {
		if state.Bool("querying") {
			gr.Text("Loading...").Modify(response)
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

func (i InstanceClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

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

	selectOne("Instance Type", "instanceType", instanceTypes, &state, i.storeValue).Modify(classEditForm)
	selectMultiple("Security Groups", "securityGroups", classOptions["securitygroups"], &state, i.storeValue).Modify(classEditForm)
	selectMultiple("EBS Volumes", "ebsVolumes", classOptions["volumes"], &state, i.storeValue).Modify(classEditForm)
	selectOne("Vpc", "vpc", classOptions["vpcs"], &state, i.storeValue).Modify(classEditForm)
	selectOne("Subnet", "subnet", classOptions["subnets"], &state, i.storeValue).Modify(classEditForm)
	checkbox("Public IP Address", "publicIpAddress", &state, i.storeValue).Modify(classEditForm)
	selectOne("AMI", "ami", classOptions["images"], &state, i.storeValue).Modify(classEditForm)
	selectOne("Key Name", "keyName", classOptions["keypairs"], &state, i.storeValue).Modify(classEditForm)
	checkbox("EBS Optimized", "ebsOptimized", &state, i.storeValue).Modify(classEditForm)
	checkbox("Monitoring", "monitoring", &state, i.storeValue).Modify(classEditForm)
	selectOne("Shutdown Behavior", "shutdownBehavior", shutdownBehaviors, &state, i.storeValue).Modify(classEditForm)
	selectOne("IAM User", "iamUser", classOptions["iamusers"], &state, i.storeValue).Modify(classEditForm)
	textArea("User Data", "userData", &state, i.storeValue).Modify(classEditForm)

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

func (i InstanceClassForm) backButton(*gr.Event) {
	i.SetState(gr.State{"success": ""})
	i.Props().Call("backButton")
}

func (i InstanceClassForm) doneButton(*gr.Event) {
	i.SetState(gr.State{"success": ""})
	i.Props().Call("hideAllModals")
}

func (i InstanceClassForm) saveButton(*gr.Event) {
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

func (i InstanceClassForm) deleteButton(*gr.Event) {
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

func (i InstanceClassForm) storeValue(event *gr.Event) {
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
