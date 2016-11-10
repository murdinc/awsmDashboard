package components

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsmDashboard/components/widgets"
	"github.com/murdinc/awsmDashboard/helpers"
)

type Dashboard struct {
	*gr.This
}

// Implements the StateInitializer interface
func (d Dashboard) GetInitialState() gr.State {
	return gr.State{"querying": false, "error": "", "widgetList": nil}
}

func (d Dashboard) Render() gr.Component {

	// Dashboard placeholder
	response := el.Div()

	if widgetList := d.State().Interface("widgetList"); widgetList != nil {
		widgets := WidgetsBuilder(widgetList) // Build the widgets
		widgets.Modify(response)
	} else if d.State().Bool("querying") {
		gr.Text("Loading...").Modify(response)
	} else if errStr := d.State().String("error"); errStr != "" {
		gr.Text(errStr).Modify(response)
	} else {
		gr.Text("Nothing here!").Modify(response)
	}

	elem := el.Div(gr.CSS("content"),
		response,
	)

	return elem
}

// Implements the ComponentWillMount interface
func (d Dashboard) ComponentWillMount() {
	if apiType := d.Props().String("apiType"); apiType != "" {
		d.SetState(gr.State{"querying": true})
		endpoint := "//localhost:8081/api/" + apiType + "/widgets"

		resp, err := helpers.GetAPI(endpoint)
		if !d.IsMounted() {
			return
		}

		if err != nil {
			d.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		d.SetState(gr.State{"querying": false, "widgetList": resp})
	}
}

// Implements the ShouldComponentUpdate interface.
func (d Dashboard) ShouldComponentUpdate(this *gr.This, next gr.Cops) bool {
	return d.State().HasChanged(next.State, "widgetList") &&
		d.State().HasChanged(next.State, "querying") &&
		d.State().HasChanged(next.State, "error")
}

func WidgetsBuilder(wl interface{}) *gr.Element {
	widgetList := wl.([]byte)

	jsonParsed, _ := gabs.ParseJSON(widgetList)
	widgetsJson, _ := jsonParsed.S("widgets").ChildrenMap()

	if len(widgetsJson) < 1 {
		return el.Div(gr.Text("Nothing here!"))
	}

	response := el.Div()

	for widgetName, w := range widgetsJson {
		widget := w.Data().(map[string]interface{})

		if widget["enabled"].(bool) == true {

			switch widgetName {
			case "events":
				gr.New(&widgets.EventsWidget{}).CreateElement(gr.Props{}).Modify(response)

			case "securitybulletins":
				gr.New(&widgets.SecurityBulletinsWidget{}).CreateElement(gr.Props{}).Modify(response)

			case "alarms":
				el.Div(gr.Text("Alarms Widget goes here!")).Modify(response)

			default:
				println("WidgetsBuilder does not have a switch for widget:")
				println(widgetName)
			}

		}

	}

	return response
}
