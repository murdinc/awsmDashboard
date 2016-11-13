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

type RSSWidget struct {
	*gr.This
}

// Implements the StateInitializer interface
func (r RSSWidget) GetInitialState() gr.State {
	return gr.State{"querying": true, "error": "", "success": ""}
}

// Implements the ComponentWillMount interface
func (r RSSWidget) ComponentWillMount() {

	props := r.Props()

	r.SetState(gr.State{"querying": true})

	// Get our options for the form
	go func() {
		endpoint := "//localhost:8081/api/dashboard/widgets/feed/" + props.String("name")
		resp, err := helpers.GetAPI(endpoint)
		if !r.IsMounted() {
			return
		}
		if err != nil {
			r.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		jsonParsed, _ := gabs.ParseJSON(resp)
		success, ok := jsonParsed.S("success").Data().(bool)
		if !ok || !success {
			r.SetState(gr.State{"querying": false, "error": "Error while processing API response"})
			return
		}

		r.SetState(gr.State{"itemsList": resp, "querying": false})
	}()
}

func (r RSSWidget) Render() gr.Component {

	state := r.State()
	props := r.Props()

	// Widget placeholder
	response := el.Div(gr.CSS("panel", "widget"))

	title := props.String("title")
	if title == "" {
		title = "Unnamed RSS Widget"
	}

	el.Div(gr.CSS("panel-heading"), gr.Text(title)).Modify(response)

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
		widget.Modify(response)
		return response
	}

	jsonParsed, _ := gabs.ParseJSON(itemsList)
	itemsJson := jsonParsed.S("feed").Bytes()
	err := json.Unmarshal(itemsJson, &items)
	if err != nil {
		println(err.Error())
	}

	if len(items) < 1 {
		gr.Text("Nothing here!").Modify(widget)
		widget.Modify(response)
		return response
	}

	r.BuildItemsTable(items).Modify(widget)

	widget.Modify(response)
	return response
}

func (r RSSWidget) BuildItemsTable(items []models.FeedItem) *gr.Element {
	count := 0

	response := el.Div()

	var header []string
	rows := make([][]string, len(items))
	links := make([]map[string]string, len(items))

	for i, item := range items {
		models.ExtractAwsmTableLinks(i, item, &header, &rows, &links)
		count++
	}

	tBody := el.TableBody()

	helpers.BuildTableRowsLinks(header, rows, links, tBody)

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
