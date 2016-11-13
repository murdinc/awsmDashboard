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

type AwsBlogWidget struct {
	*gr.This
}

// Implements the StateInitializer interface
func (a AwsBlogWidget) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": ""}
}

// Implements the ComponentWillMount interface
func (a AwsBlogWidget) ComponentWillMount() {
	var class map[string]interface{}

	if a.Props().Interface("class") != nil {
		classJson := a.Props().Interface("class").([]byte)
		json.Unmarshal(classJson, &class)
	}

	a.SetState(class)
	a.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/dashboard/widgets/awsblog"
		resp, err := helpers.GetAPI(endpoint)
		if !a.IsMounted() {
			return
		}
		if err != nil {
			a.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		jsonParsed, _ := gabs.ParseJSON(resp)
		success, ok := jsonParsed.S("success").Data().(bool)
		if !ok || !success {
			a.SetState(gr.State{"querying": false, "error": "Error while processing API response"})
			return
		}

		a.SetState(gr.State{"itemsList": resp, "querying": false})
	}()
}

func (a AwsBlogWidget) Render() gr.Component {

	state := a.State()
	//props := a.Props()

	// Widget placeholder
	response := el.Div(gr.CSS("panel", "widget"))
	el.Div(gr.CSS("panel-heading"), gr.Text("AWS Blog")).Modify(response)
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
	itemsJson := jsonParsed.S("blogPosts").Bytes()
	err := json.Unmarshal(itemsJson, &items)
	if err != nil {
		println(err.Error())
	}

	if len(items) < 1 {
		gr.Text("Nothing here!").Modify(response)
		return response
	}

	a.BuildItemsTable(items).Modify(widget)

	widget.Modify(response)
	return response
}

func (a AwsBlogWidget) BuildItemsTable(items []models.FeedItem) *gr.Element {
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
