package widgets

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsm/models"
	"github.com/murdinc/awsmDashboard/helpers"
)

type EventsWidget struct {
	*gr.This
}

// Implements the StateInitializer interface
func (e EventsWidget) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": ""}
}

// Implements the ComponentWillMount interface
func (e EventsWidget) ComponentWillMount() {
	var class map[string]interface{}

	if e.Props().Interface("class") != nil {
		classJson := e.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	e.SetState(class)
	e.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/dashboard/widgets/events"
		resp, err := helpers.GetAPI(endpoint)

		if !e.IsMounted() {
			return
		}

		if err != nil {
			e.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		jsonParsed, _ := gabs.ParseJSON(resp)
		success, ok := jsonParsed.S("success").Data().(bool)
		if !ok || !success {
			e.SetState(gr.State{"querying": false, "error": "Error while processing API response"})
			return
		}

		e.SetState(gr.State{"eventsList": resp, "querying": false})
	}()
}

func (e EventsWidget) Render() gr.Component {

	state := e.State()
	//props := e.Props()

	// Widget placeholder
	response := el.Div(gr.CSS("panel", "widget"))
	el.Div(gr.CSS("panel-heading"), gr.Text("AWS Events")).Modify(response)
	widget := el.Div(gr.CSS("panel-body"))

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(widget)

	if state.Bool("querying") {
		el.Div(gr.CSS("panel-body"), gr.Text("Loading...")).Modify(response)
		return response
	}

	var events []models.Event
	eventsList, ok := state.Interface("eventsList").([]byte)
	if !ok {
		return response
	}

	jsonParsed, _ := gabs.ParseJSON(eventsList)

	eventsJson := jsonParsed.S("events").Bytes()
	json.Unmarshal(eventsJson, &events)

	if len(events) < 1 {
		gr.Text("Nothing here!").Modify(response)
		return response
	}
	e.BuildEventsTable(events).Modify(widget)

	widget.Modify(response)
	return response
}

func (e EventsWidget) BuildEventsTable(events []models.Event) *gr.Element {
	count := 0
	maxArchive := 10

	response := el.Div()

	var header []string
	rows := make([][]string, len(events))

Loop:
	for i, event := range events {
		if count >= maxArchive {
			break Loop
		}

		models.ExtractAwsmTable(i, event, &header, &rows)
		if event.Archive {
			count++
		}
	}

	tBody := el.TableBody()

	helpers.BuildTableRows(rows, tBody)

	table := el.Table(
		gr.CSS("table", "table-striped"),
		gr.Style("width", "100%"),
		el.TableHead(el.TableRow(helpers.BuildTableHeader(header)...)))

	tBody.Modify(table)

	table.Modify(response)

	return response

}
