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

type RSSWidgetForm struct {
	*gr.This
}

// Implements the StateInitializer interface
func (r RSSWidgetForm) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": "", "step": 1,
		"enabled": true, "widgetType": "rss",
	}
}

// Implements the ComponentWillMount interface
func (r RSSWidgetForm) ComponentWillMount() {
	var widget map[string]interface{}

	if r.Props().Interface("widget") != nil {
		widgetJson := r.Props().Interface("widget").([]byte)
		json.Unmarshal(widgetJson, &widget)
	}

	r.SetState(widget)
	r.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/" + r.Props().String("apiType") + "/widgets/options"
		resp, err := helpers.GetAPI(endpoint)
		if !r.IsMounted() {
			return
		}
		if err != nil {
			r.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		r.SetState(gr.State{"widgetOptionsResp": resp, "querying": false})
	}()
}

func (r RSSWidgetForm) Render() gr.Component {

	state := r.State()
	props := r.Props()

	// Form placeholder
	response := el.Div()

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(response)
	helpers.SuccessElem(state.String("success")).Modify(response)

	if state.Int("step") == 1 {
		if state.Bool("querying") {
			gr.Text("Loading...").Modify(response)
		} else {
			r.BuildWidgetForm(props.String("widgetName"), state.Interface("widgetOptionsResp")).Modify(response)
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
				evt.Click(r.backButton).PreventDefault(),
				gr.CSS("btn", "btn-secondary"),
				gr.Text("Back"),
			).Modify(buttons)

			// Done
			el.Button(
				evt.Click(r.doneButton).PreventDefault(),
				gr.CSS("btn", "btn-primary"),
				gr.Text("Done"),
			).Modify(buttons)

			buttons.Modify(response)
		}

	}

	return response
}

func (r RSSWidgetForm) BuildWidgetForm(widgetName string, optionsResp interface{}) *gr.Element {

	state := r.State()
	props := r.Props()

	var widgetOptions map[string][]string
	jsonParsed, _ := gabs.ParseJSON(optionsResp.([]byte))
	widgetOptionsJson := jsonParsed.S("widgetOptions").Bytes()
	json.Unmarshal(widgetOptionsJson, &widgetOptions)

	widgetEdit := el.Div(
		el.Header3(gr.Text(widgetName)),
		el.HorizontalRule(),
	)

	widgetEditForm := el.Form(evt.KeyDown(DisableEnter))

	TextField("Title", "title", state.String("title"), r.storeValue).Modify(widgetEditForm)
	TextField("RSS URL", "rssUrl", state.String("rssUrl"), r.storeValue).Modify(widgetEditForm)
	NumberField("Count", "count", state.Int("count"), r.storeValue).Modify(widgetEditForm)
	NumberField("Index", "index", state.Int("index"), r.storeValue).Modify(widgetEditForm)
	Checkbox("Enabled", "enabled", state.Bool("enabled"), r.storeValue).Modify(widgetEditForm)

	widgetEditForm.Modify(widgetEdit)

	buttons := el.Div(
		gr.CSS("btn-toolbar"),
	)

	// Back
	el.Button(
		evt.Click(r.backButton).PreventDefault(),
		gr.CSS("btn", "btn-secondary"),
		gr.Text("Back"),
	).Modify(buttons)

	// Save
	el.Button(
		evt.Click(r.saveButton).PreventDefault(),
		gr.CSS("btn", "btn-primary"),
		gr.Text("Save"),
	).Modify(buttons)

	// Delete
	if props.Interface("hasDelete") != nil && props.Bool("hasDelete") {
		el.Button(
			evt.Click(r.deleteButton).PreventDefault(),
			gr.CSS("btn", "btn-danger", "pull-right"),
			gr.Text("Delete"),
		).Modify(buttons)
	}

	buttons.Modify(widgetEdit)

	return widgetEdit

}

func (r RSSWidgetForm) backButton(*gr.Event) {
	r.SetState(gr.State{"success": ""})
	r.Props().Call("backButton")
}

func (r RSSWidgetForm) doneButton(*gr.Event) {
	r.SetState(gr.State{"success": ""})
	r.Props().Call("hideAllModals")
}

func (r RSSWidgetForm) saveButton(*gr.Event) {
	r.SetState(gr.State{"querying": true, "step": 2})

	cfg := make(map[string]interface{})
	for key, _ := range r.State() {
		cfg[key] = r.State().Interface(key)
	}

	go func() {
		endpoint := "//localhost:8081/api/" + r.Props().String("apiType") + "/widgets/name/" + r.Props().String("widgetName")

		_, err := helpers.PutAPI(endpoint, cfg)
		if !r.IsMounted() {
			return
		}

		if err != nil {
			r.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint), "step": 1})
			return
		}

		r.SetState(gr.State{"querying": false, "success": "Widget was saved", "error": ""})
	}()

}

func (r RSSWidgetForm) deleteButton(*gr.Event) {
	r.SetState(gr.State{"querying": true})

	go func() {
		endpoint := "//localhost:8081/api/" + r.Props().String("apiType") + "/widgets/name/" + r.Props().String("widgetName")

		_, err := helpers.DeleteAPI(endpoint)
		if !r.IsMounted() {
			return
		}

		if err != nil {
			r.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		r.SetState(gr.State{"querying": false, "success": "Class was deleted", "error": "", "step": 2})
	}()
}

func (r RSSWidgetForm) storeValue(event *gr.Event) {
	id := event.Target().Get("id").String()
	inputType := event.Target().Get("type").String()

	switch inputType {

	case "checkbox":
		r.SetState(gr.State{id: event.Target().Get("checked").Bool()})

	case "number":
		r.SetState(gr.State{id: event.TargetValue().Int()})

	default: // text, at least
		r.SetState(gr.State{id: event.TargetValue()})

	}
}

func (r RSSWidgetForm) storeSelect(id string, val interface{}) {
	switch value := val.(type) {

	case map[string]interface{}:
		// single
		r.SetState(gr.State{id: value["value"]})

	case []interface{}:
		// multi
		var vals []string
		options := len(value)
		for i := 0; i < options; i++ {
			vals = append(vals, value[i].(map[string]interface{})["value"].(string))
		}
		r.SetState(gr.State{id: vals})

	default:
		r.SetState(gr.State{id: val})

	}
}
