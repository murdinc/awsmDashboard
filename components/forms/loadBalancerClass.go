package forms

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
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
	return gr.State{"querying": true, "error": "", "success": "", "step": 1,
		"loadBalancerListeners": []interface{}{},
	}
}

// Implements the ComponentWillMount interface
func (l LoadBalancerClassForm) ComponentWillMount() {
	var class map[string]interface{}

	if l.Props().Interface("class") != nil {
		classJson := l.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	// health check
	healthCheck := class["loadBalancerHealthCheck"].(map[string]interface{})
	for k, v := range healthCheck {
		class[k] = v
	}

	// attributes
	attributes := class["loadBalancerAttributes"].(map[string]interface{})
	for k, v := range attributes {
		class[k] = v
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

	classEditForm := el.Form(evt.KeyDown(DisableEnter))

	SelectOne("Scheme", "scheme", elbSchemes, state.Interface("scheme"), l.storeSelect).Modify(classEditForm)
	SelectOne("Vpc", "vpc", classOptions["vpcs"], state.Interface("vpc"), l.storeSelect).Modify(classEditForm)
	if vpc, ok := state.Interface("vpc").(string); ok && len(vpc) > 0 {
		SelectMultiple("Subnets", "subnets", classOptions["subnets"], state.Interface("subnets"), l.storeSelect).Modify(classEditForm)
		SelectMultiple("Security Groups", "securityGroups", classOptions["securitygroups"], state.Interface("securityGroups"), l.storeSelect).Modify(classEditForm)
	}

	SelectMultiple("Availability Zones", "availabilityZones", classOptions["zones"], state.Interface("availabilityZones"), l.storeSelect).Modify(classEditForm)
	Checkbox("Cross Zone Load Balancing", "crossZoneLoadBalancingEnabled", state.Bool("crossZoneLoadBalancingEnabled"), l.storeValue).Modify(classEditForm)

	NumberField("Idle Timeout", "idleTimeout", state.Int("idleTimeout"), l.storeValue).Modify(classEditForm)

	el.Div(
		el.Break(nil),
		el.Header4(
			gr.Text("Health Check"),
		),
		el.HorizontalRule(nil),
	).Modify(classEditForm)

	TextField("Target", "healthCheckTarget", state.String("healthCheckTarget"), l.storeValue).Modify(classEditForm)
	NumberField("Timeout", "healthCheckTimeout", state.Int("healthCheckTimeout"), l.storeValue).Modify(classEditForm)
	NumberField("Interval", "healthCheckInterval", state.Int("healthCheckInterval"), l.storeValue).Modify(classEditForm)
	NumberField("Unhealthy Threshold", "healthCheckUnhealthyThreshold", state.Int("healthCheckUnhealthyThreshold"), l.storeValue).Modify(classEditForm)
	NumberField("Healthy Threshold", "healthCheckHealthyThreshold", state.Int("healthCheckHealthyThreshold"), l.storeValue).Modify(classEditForm)

	el.Div(
		el.Break(nil),
		el.Header4(
			gr.Text("Connection Draining"),
		),
		el.HorizontalRule(nil),
	).Modify(classEditForm)

	Toggle("Disabled", "Enabled", "connectionDrainingEnabled", state.Bool("connectionDrainingEnabled"), l.storeValue).Modify(classEditForm)
	if state.Bool("connectionDrainingEnabled") {
		NumberField("Draining Timeout", "connectionDrainingTimeout", state.Int("connectionDrainingTimeout"), l.storeValue).Modify(classEditForm)
	}

	el.Div(
		el.Break(nil),
		el.Header4(
			gr.Text("Access Log"),
		),
		el.HorizontalRule(nil),
	).Modify(classEditForm)

	Toggle("Disabled", "Enabled", "accessLogEnabled", state.Bool("accessLogEnabled"), l.storeValue).Modify(classEditForm)
	if state.Bool("accessLogEnabled") {
		NumberField("Emit Interval", "accessLogEmitInterval", state.Int("accessLogEmitInterval"), l.storeValue).Modify(classEditForm)
		TextField("S3 Bucket Name", "accessLogS3BucketName", state.String("accessLogS3BucketName"), l.storeValue).Modify(classEditForm)
		TextField("S3 Bucket Prefix", "accessLogS3BucketPrefix", state.String("accessLogS3BucketPrefix"), l.storeValue).Modify(classEditForm)
	}

	el.Div(
		el.Break(nil),
		el.Header4(
			gr.Text("Listeners"),
			el.Button(
				evt.Click(l.addListener).PreventDefault(),
				gr.CSS("btn", "btn-primary", "btn-sm", "pull-right"),
				gr.Text("New"),
			),
		),
		el.HorizontalRule(nil),
	).Modify(classEditForm)

	listeners := state.Interface("loadBalancerListeners").([]interface{})

	for index, listInf := range listeners {

		listener := listInf.(map[string]interface{})

		// Form placeholder
		listenerForm := el.Div()

		//TextField("Note", "note", listener["note"], l.modifyListener(index, listener)).Modify(listenerForm)

		el.Div(
			gr.CSS("row"), el.Div(gr.CSS("col-sm-6"),
				SelectOne("Protocol", "protocol", []string{"HTTP", "HTTPS", "TCP", "SSL"}, listener["protocol"], l.storeListenerSelect(index, listener)),
			),
			el.Div(gr.CSS("col-sm-6"),
				NumberField("Load Balancer Port", "loadBalancerPort", listener["loadBalancerPort"], l.modifyListener(index, listener)),
			),
		).Modify(listenerForm)

		el.Div(
			gr.CSS("row"), el.Div(gr.CSS("col-sm-6"),
				SelectOne("Instance Protocol", "instanceProtocol", []string{"HTTP", "HTTPS", "TCP", "SSL"}, listener["instanceProtocol"], l.storeListenerSelect(index, listener)),
			),
			el.Div(gr.CSS("col-sm-6"),
				NumberField("Instance Port", "instancePort", listener["instancePort"], l.modifyListener(index, listener)),
			),
		).Modify(listenerForm)

		el.Div(
			gr.CSS("btn-toolbar"),
			el.Button(
				evt.Click(l.removeListener).PreventDefault(),
				gr.CSS("btn", "btn-danger", "btn-sm", "pull-right"),
				gr.Text("Remove"),
				attr.ID(index),
			),
		).Modify(listenerForm)

		el.HorizontalRule(nil).Modify(listenerForm)

		listenerForm.Modify(classEditForm)

	}

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

	loadBalancerHealthCheck := make(map[string]interface{})
	loadBalancerHealthCheck["healthCheckTarget"] = cfg["healthCheckTarget"]
	loadBalancerHealthCheck["healthCheckTimeout"] = cfg["healthCheckTimeout"]
	loadBalancerHealthCheck["healthCheckInterval"] = cfg["healthCheckInterval"]
	loadBalancerHealthCheck["healthCheckUnhealthyThreshold"] = cfg["healthCheckUnhealthyThreshold"]
	loadBalancerHealthCheck["healthCheckHealthyThreshold"] = cfg["healthCheckHealthyThreshold"]
	cfg["loadBalancerHealthCheck"] = loadBalancerHealthCheck

	loadBalancerAttributes := make(map[string]interface{})
	loadBalancerAttributes["connectionDrainingEnabled"] = cfg["connectionDrainingEnabled"]
	loadBalancerAttributes["connectionDrainingTimeout"] = cfg["connectionDrainingTimeout"]
	loadBalancerAttributes["idleTimeout"] = cfg["idleTimeout"]
	loadBalancerAttributes["crossZoneLoadBalancingEnabled"] = cfg["crossZoneLoadBalancingEnabled"]
	loadBalancerAttributes["accessLogEnabled"] = cfg["accessLogEnabled"]
	loadBalancerAttributes["accessLogEmitInterval"] = cfg["accessLogEmitInterval"]
	loadBalancerAttributes["accessLogS3BucketName"] = cfg["accessLogS3BucketName"]
	loadBalancerAttributes["accessLogS3BucketPrefix"] = cfg["accessLogS3BucketPrefix"]
	cfg["loadBalancerAttributes"] = loadBalancerAttributes

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
	key := event.Target().Get("name").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		l.SetState(gr.State{key: event.Target().Get("checked").Bool()})

	case "number":
		l.SetState(gr.State{key: event.TargetValue().Int()})

	default: // text, at least
		l.SetState(gr.State{key: event.TargetValue()})

	}
}

func (l LoadBalancerClassForm) storeSelect(key string, val interface{}) {
	switch value := val.(type) {

	case map[string]interface{}:
		// single
		l.SetState(gr.State{key: value["value"]})

	case []interface{}:
		// multi
		var vals []string
		options := len(value)
		for i := 0; i < options; i++ {
			vals = append(vals, value[i].(map[string]interface{})["value"].(string))
		}
		l.SetState(gr.State{key: vals})

	default:
		l.SetState(gr.State{key: val})

	}
}

func (l LoadBalancerClassForm) modifyListener(index int, listener map[string]interface{}) func(*gr.Event) {
	return func(event *gr.Event) {
		key := event.Target().Get("name").String()
		valueType := event.Target().Get("type").String()

		switch valueType {
		case "text":
			listener[key] = event.TargetValue().String()
		case "number":
			listener[key] = event.TargetValue().Int()
		}
		listeners, ok := l.State().Interface("loadBalancerListeners").([]interface{})
		if ok {
			listeners[index] = listener
			l.SetState(gr.State{"loadBalancerListeners": listeners})
			return
		}
		println("modifyListener failed?")
	}
}

func (l LoadBalancerClassForm) storeListenerSelect(index int, listener map[string]interface{}) func(string, interface{}) {
	return func(key string, val interface{}) {

		switch value := val.(type) {

		case map[string]interface{}:
			// single
			listener[key] = value["value"]

		case []interface{}:
			// multi
			var vals []string
			options := len(value)
			for i := 0; i < options; i++ {
				vals = append(vals, value[i].(map[string]interface{})["value"].(string))
			}
			listener[key] = vals

		default:
			listener[key] = val
		}

		listeners, ok := l.State().Interface("loadBalancerListeners").([]interface{})
		if ok {
			listeners[index] = listener
			l.SetState(gr.State{"loadBalancerListeners": listeners})
			return
		}
	}
}

func (l LoadBalancerClassForm) addListener(*gr.Event) {
	listeners, ok := l.State().Interface("loadBalancerListeners").([]interface{})
	if ok {

		newListener := make(map[string]interface{})
		newListener["protocol"] = "tcp"
		newListener["instanceProtocol"] = "tcp"
		newListener["loadBalancerPort"] = 0
		newListener["instancePort"] = 0

		listeners = append([]interface{}{newListener}, listeners...)

		l.SetState(gr.State{"loadBalancerListeners": listeners})
		return
	}
	println("addListener failed?")
}

func (l LoadBalancerClassForm) removeListener(event *gr.Event) {
	index := event.Target().Get("id").Int()
	listeners, ok := l.State().Interface("loadBalancerListeners").([]interface{})
	if ok {
		listeners = append(listeners[:index], listeners[index+1:]...)
		l.SetState(gr.State{"loadBalancerListeners": listeners})
		return
	}
	println("removeListener failed?")
}
