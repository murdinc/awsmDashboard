package components

import (
	"encoding/json"
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsm/models"
	"github.com/murdinc/awsmDashboard/helpers"
)

type AssetTable struct {
	*gr.This
}

// Implements the StateInitializer interface
func (a AssetTable) GetInitialState() gr.State {
	return gr.State{"querying": false, "error": "", "assetList": nil}
}

func (a AssetTable) Render() gr.Component {

	// Table placeholder
	response := el.Div()

	elem := el.Div(gr.CSS("content"),
		response,
	)

	if assets := a.State().Interface("assetList"); assets != nil {
		table := AssetTableBuilder(assets) // Build the table
		table.Modify(response)

		el.Break().Modify(response)
		el.HorizontalRule().Modify(response)

	} else if a.State().Bool("querying") {
		gr.Text("Loading...").Modify(response)
	} else if errStr := a.State().String("error"); errStr != "" {
		gr.Text(errStr).Modify(response)
	} else {
		gr.Text("Nothing here!").Modify(response)
	}

	return elem
}

// Implements the ComponentWillMount interface
func (a AssetTable) ComponentWillMount() {
	if apiType := a.Props().String("apiType"); apiType != "" {
		a.SetState(gr.State{"querying": true})
		endpoint := "//localhost:8081/api/assets/" + apiType

		resp, err := helpers.GetAPI(endpoint)
		if !a.IsMounted() {
			return
		}
		if err != nil {
			a.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		a.SetState(gr.State{"querying": false, "assetList": resp})
	}
}

// Implements the ShouldComponentUpdate interface.
func (a AssetTable) ShouldComponentUpdate(this *gr.This, next gr.Cops) bool {
	return a.State().HasChanged(next.State, "assetList") &&
		a.State().HasChanged(next.State, "querying") &&
		a.State().HasChanged(next.State, "error")
}

func AssetTableBuilder(al interface{}) *gr.Element {
	assetList := al.([]byte)

	jsonParsed, _ := gabs.ParseJSON(assetList)
	assetType := jsonParsed.S("assetType").Data().(string)
	assets, _ := jsonParsed.S("assets").Children()

	if len(assets) < 1 {
		return el.Div(gr.Text("Nothing here!"))
	}

	var header []string
	rows := make([][]string, len(assets))

	switch assetType {

	case "alarms":
		var aType models.Alarm
		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "buckets":
		var aType models.Bucket
		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "instances":
		var aType models.Instance
		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "volumes":
		var aType models.Volume
		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "images":
		var aType models.Image

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "keypairs":
		var aType models.KeyPair
		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "snapshots":
		var aType models.Snapshot

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "vpcs":
		var aType models.Vpc

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "subnets":
		var aType models.Subnet

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "securitygroups":
		var aType models.SecurityGroup

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "addresses":
		var aType models.Address

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "launchconfigurations":
		var aType models.LaunchConfig

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "autoscalegroups":
		var aType models.AutoScaleGroup

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "loadbalancers":
		var aType models.LoadBalancer

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "scalingpolicies":
		var aType models.ScalingPolicy

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	case "simpledbdomains":
		var aType models.SimpleDBDomain

		for i, a := range assets {
			aJson := a.Bytes()
			json.Unmarshal(aJson, &aType)
			models.ExtractAwsmTable(i, aType, &header, &rows)
		}

	default:
		println("Asset Type not found in AssetTableBuilder switch:")
		println(assetType)
	}

	tBody := el.TableBody()

	helpers.BuildTableRows(rows, tBody)

	table := el.Table(
		gr.CSS("table", "table-striped"),
		gr.Style("width", "100%"),
		el.TableHead(el.TableRow(helpers.BuildTableHeader(header)...)))

	tBody.Modify(table)

	return table
}
