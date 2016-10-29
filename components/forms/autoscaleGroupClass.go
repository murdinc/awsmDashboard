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
	}
}

// Implements the ComponentDidMount interface
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
			gr.Text("Loading.......").Modify(response)
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

	classEditForm := el.Form()

	selectOne("Launch Configuration Class", "launchConfigurationClass", classOptions["launchconfigurations"], &state, a.storeValue).Modify(classEditForm)

	//checkbox("Rotate", "rotate", &state, a.storeValue).Modify(classEditForm)
	//if state.Bool("rotate") {
	textField("Retain", "retain", &state, a.storeValue).Modify(classEditForm) // number
	//}
	checkbox("Propagate", "propagate", &state, a.storeValue).Modify(classEditForm)
	if state.Bool("propagate") {
		//selectMultiple("Propagate Regions", "propagateRegions", classOptions["regions"], &state, a.storeValue).Modify(classEditForm)
	}

	selectMultiple("Availability Zones", "availabilityZones", classOptions["zones"], &state, a.storeValue).Modify(classEditForm)
	textField("Desired Capacity", "desiredCapacity", &state, a.storeValue).Modify(classEditForm) // number
	textField("Min Size", "minSize", &state, a.storeValue).Modify(classEditForm)                 // number
	textField("Max Size", "maxSize", &state, a.storeValue).Modify(classEditForm)                 // number
	textField("Default Cooldown", "defaultCooldown", &state, a.storeValue).Modify(classEditForm) // number
	selectOne("Subnet Class", "subnetClass", classOptions["subnets"], &state, a.storeValue).Modify(classEditForm)
	selectOne("Health Check Type", "healthCheckType", healthCheckTypes, &state, a.storeValue).Modify(classEditForm)
	textField("Health Check Grace Period", "healthCheckGracePeriod", &state, a.storeValue).Modify(classEditForm) // number
	selectMultiple("Termination Policies", "terminationPolicies", terminationPolicies, &state, a.storeValue).Modify(classEditForm)
	selectMultiple("Scaling Policies", "scalingPolicies", classOptions["scalingpolicies"], &state, a.storeValue).Modify(classEditForm)
	selectMultiple("Load Balancer Names", "loadBalancerNames", classOptions["loadbalancers"], &state, a.storeValue).Modify(classEditForm)
	selectMultiple("Alarms", "alarms", classOptions["alarms"], &state, a.storeValue).Modify(classEditForm)

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
	id := event.Target().Get("id").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		a.SetState(gr.State{id: event.Target().Get("checked").Bool()})

	case "select-one":
		a.SetState(gr.State{id: event.TargetValue()})

	case "select-multiple":
		var vals []string
		options := event.Target().Length()

		for i := 0; i < options; i++ {
			if event.Target().Index(i).Get("selected").Bool() && event.Target().Index(i).Get("id") != nil {
				vals = append(vals, event.Target().Index(i).Get("id").String())
			}
		}
		a.SetState(gr.State{id: vals})

	default: // text, at least
		a.SetState(gr.State{id: event.TargetValue()})

	}
}
