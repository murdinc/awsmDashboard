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

type SecurityBulletinsWidget struct {
	*gr.This
}

// Implements the StateInitializer interface
func (s SecurityBulletinsWidget) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": ""}
}

// Implements the ComponentWillMount interface
func (s SecurityBulletinsWidget) ComponentWillMount() {
	var class map[string]interface{}

	if s.Props().Interface("class") != nil {
		classJson := s.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	s.SetState(class)
	s.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/dashboard/widgets/securitybulletins"
		resp, err := helpers.GetAPI(endpoint)
		if !s.IsMounted() {
			return
		}
		if err != nil {
			s.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		jsonParsed, _ := gabs.ParseJSON(resp)
		success, ok := jsonParsed.S("success").Data().(bool)
		if !ok || !success {
			s.SetState(gr.State{"querying": false, "error": "Error while processing API response"})
			return
		}

		s.SetState(gr.State{"itemsList": resp, "querying": false})
	}()
}

func (s SecurityBulletinsWidget) Render() gr.Component {

	state := s.State()
	//props := s.Props()

	// Widget placeholder
	response := el.Div(gr.CSS("panel", "widget"))
	el.Div(gr.CSS("panel-heading"), gr.Text("AWS Security Bulletins")).Modify(response)
	widget := el.Div(gr.CSS("panel-body"))

	// Print any alerts
	helpers.ErrorElem(state.String("error")).Modify(widget)

	if state.Bool("querying") {
		el.Div(gr.CSS("panel-body"), gr.Text("Loading...")).Modify(response)
		return response
	}

	var items []models.FeedItem
	itemsList, ok := state.Interface("itemsList").([]byte)
	if !ok {
		return response
	}

	jsonParsed, _ := gabs.ParseJSON(itemsList)
	itemsJson := jsonParsed.S("securityBulletins").Bytes()
	err := json.Unmarshal(itemsJson, &items)
	if err != nil {
		println(err.Error())
	}

	if len(items) < 1 {
		gr.Text("Nothing here!").Modify(response)
		return response
	}

	s.BuildItemsTable(items).Modify(widget)

	widget.Modify(response)
	return response
}

func (s SecurityBulletinsWidget) BuildItemsTable(items []models.FeedItem) *gr.Element {
	count := 0
	maxItems := 10

	response := el.Div()

	var header []string
	rows := make([][]string, len(items))

	for i, item := range items {
		if count < maxItems {
			models.ExtractAwsmTable(i, item, &header, &rows)
			count++
		}
	}

	tBody := el.TableBody()

	helpers.BuildTableRows(rows, tBody)

	table := el.Table(
		gr.CSS("table", "table-striped"),
		gr.Style("width", "100%"),
		el.TableHead(el.TableRow(helpers.BuildTableHeader(header)...)))

	if count < 1 {
		gr.Text("Nothing here!").Modify(response)
		return response
	}

	tBody.Modify(table)

	table.Modify(response)

	return response

}
