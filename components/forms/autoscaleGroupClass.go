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
	terminationPolicies = []string{"OldestInstance", "NewestInstance", "OldestLaunchConfiguration", "ClosestToNextInstanceHour", "Default"}
	healthCheckTypes    = []string{"EC2", "ELB"}
)

type AutoscaleGroupClassForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (a AutoscaleGroupClassForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "",
		"monitoring":      false,
		"publicIpAddress": false,
		"ebsOptimized":    false,
		"step":            1,
		"propagate":       false,
	}
}

// Implements the ComponentWillMount interface
func (a AutoscaleGroupClassForm) ComponentWillMount() {
	var class map[string]interface{}

	if a.Props().Interface("class") != nil {
		classJson := a.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	a.SetState(class)
	a.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/classes/" + a.Props().String("apiType") + "/options"
		resp, err := helpers.GetAPI(endpoint)
		if !a.IsMounted() {
			return
		}
		if err != nil {
			a.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		a.SetState(gr.State{"classOptionsResp": resp, "querying": false})
	}()
}

func (a AutoscaleGroupClassForm) Render() gr.Component {

	state := a.State()
	props := a.Props()

	// Form placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)
	helpers.SuccessElem(state.String("success")).Modify(response)

	if state.Int("step") == 1 {
		if state.Bool("querying") {
			gr.Text("Loading...").Modify(response)
		} else {
			a.BuildClassForm(props.String("className"), state.Interface("classOptionsResp")).Modify(response)
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
				evt.Click(a.backButton).PreventDefault(),
				gr.CSS("btn", "btn-secondary"),
				gr.Text("Back"),
			).Modify(buttons)

			// Done
			el.Button(
				evt.Click(a.doneButton).PreventDefault(),
				gr.CSS("btn", "btn-primary"),
				gr.Text("Done"),
			).Modify(buttons)

			buttons.Modify(response)
		}

	}

	return response
}

func (a AutoscaleGroupClassForm) BuildClassForm(className string, optionsResp interface{}) *gr.Element {

	state := a.State()
	props := a.Props()

	var classOptions map[string][]string
	jsonParsed, _ := gabs.ParseJSON(optionsResp.([]byte))
	classOptionsJson := jsonParsed.S("classOptions").Bytes()
	json.Unmarshal(classOptionsJson, &classOptions)

	classEdit := el.Div(
		el.Header3(gr.Text(className)),
		el.HorizontalRule(),
	)

	classEditForm := el.Form(evt.KeyDown(DisableEnter))

	SelectOne("Launch Configuration Class", "launchConfigurationClass", classOptions["launchconfigurations"], state.Interface("launchConfigurationClass"), a.storeSelect).Modify(classEditForm)

	// NEEDED?
	Checkbox("Propagate", "propagate", state.Bool("propagate"), a.storeValue).Modify(classEditForm)
	if state.Bool("propagate") {
		SelectMultiple("Propagate Regions", "propagateRegions", classOptions["regions"], state.Interface("propagateRegions"), a.storeSelect).Modify(classEditForm)
	}

	SelectMultiple("Availability Zones", "availabilityZones", classOptions["zones"], state.Interface("availabilityZones"), a.storeSelect).Modify(classEditForm)
	NumberField("Desired Capacity", "desiredCapacity", state.Int("desiredCapacity"), a.storeValue).Modify(classEditForm)
	NumberField("Min Size", "minSize", state.Int("minSize"), a.storeValue).Modify(classEditForm)
	NumberField("Max Size", "maxSize", state.Int("maxSize"), a.storeValue).Modify(classEditForm)
	NumberField("Default Cooldown", "defaultCooldown", state.Int("defaultCooldown"), a.storeValue).Modify(classEditForm)
	SelectOne("Subnet Class", "subnetClass", classOptions["subnets"], state.Interface("subnetClass"), a.storeSelect).Modify(classEditForm)
	SelectOne("Health Check Type", "healthCheckType", healthCheckTypes, state.Interface("healthCheckType"), a.storeSelect).Modify(classEditForm)
	NumberField("Health Check Grace Period", "healthCheckGracePeriod", state.Int("healthCheckGracePeriod"), a.storeValue).Modify(classEditForm)
	SelectMultiple("Termination Policies", "terminationPolicies", terminationPolicies, state.Interface("terminationPolicies"), a.storeSelect).Modify(classEditForm)
	SelectMultiple("Load Balancer Names", "loadBalancerNames", classOptions["loadbalancers"], state.Interface("loadBalancerNames"), a.storeSelect).Modify(classEditForm)
	SelectMultiple("Alarms", "alarms", classOptions["alarms"], state.Interface("alarms"), a.storeSelect).Modify(classEditForm)

	classEditForm.Modify(classEdit)

	buttons := el.Div(
		gr.CSS("btn-toolbar"),
	)

	// Back
	el.Button(
		evt.Click(a.backButton).PreventDefault(),
		gr.CSS("btn", "btn-secondary"),
		gr.Text("Back"),
	).Modify(buttons)

	// Save
	el.Button(
		evt.Click(a.saveButton).PreventDefault(),
		gr.CSS("btn", "btn-primary"),
		gr.Text("Save"),
	).Modify(buttons)

	// Delete
	if props.Interface("hasDelete") != nil && props.Bool("hasDelete") {
		el.Button(
			evt.Click(a.deleteButton).PreventDefault(),
			gr.CSS("btn", "btn-danger", "pull-right"),
			gr.Text("Delete"),
		).Modify(buttons)
	}

	buttons.Modify(classEdit)

	return classEdit

}

func (a AutoscaleGroupClassForm) backButton(*gr.Event) {
	a.SetState(gr.State{"success": ""})
	a.Props().Call("backButton")
}

func (a AutoscaleGroupClassForm) doneButton(*gr.Event) {
	a.SetState(gr.State{"success": ""})
	a.Props().Call("hideAllModals")
}

func (a AutoscaleGroupClassForm) saveButton(*gr.Event) {
	a.SetState(gr.State{"querying": true, "step": 2})

	cfg := make(map[string]interface{})
	for key, _ := range a.State() {
		cfg[key] = a.State().Interface(key)
	}

	go func() {
		endpoint := "//localhost:8081/api/classes/" + a.Props().String("apiType") + "/name/" + a.Props().String("className")

		_, err := helpers.PutAPI(endpoint, cfg)
		if !a.IsMounted() {
			return
		}

		if err != nil {
			a.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint), "step": 1})
			return
		}

		a.SetState(gr.State{"querying": false, "success": "Class was saved", "error": ""})
	}()

}

func (a AutoscaleGroupClassForm) deleteButton(*gr.Event) {
	a.SetState(gr.State{"querying": true})

	go func() {
		endpoint := "//localhost:8081/api/classes/" + a.Props().String("apiType") + "/name/" + a.Props().String("className")

		_, err := helpers.DeleteAPI(endpoint)
		if !a.IsMounted() {
			return
		}

		if err != nil {
			a.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		a.SetState(gr.State{"querying": false, "success": "Class was deleted", "error": "", "step": 2})
	}()
}

func (a AutoscaleGroupClassForm) storeValue(event *gr.Event) {
	key := event.Target().Get("name").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		a.SetState(gr.State{key: event.Target().Get("checked").Bool()})

	case "number":
		a.SetState(gr.State{key: event.TargetValue().Int()})

	default: // text, at least
		a.SetState(gr.State{key: event.TargetValue()})

	}
}

func (a AutoscaleGroupClassForm) storeSelect(key string, val interface{}) {
	switch value := val.(type) {

	case map[string]interface{}:
		// single
		a.SetState(gr.State{key: value["value"]})

	case []interface{}:
		// multi
		var vals []string
		options := len(value)
		for i := 0; i < options; i++ {
			vals = append(vals, value[i].(map[string]interface{})["value"].(string))
		}
		a.SetState(gr.State{key: vals})

	default:
		a.SetState(gr.State{key: val})

	}
}
