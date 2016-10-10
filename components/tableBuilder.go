package components

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsm/models"
)

func TableBuilder(al interface{}) *gr.Element {

	assetList := al.([]byte)

	jsonParsed, _ := gabs.ParseJSON(assetList)
	assetType := jsonParsed.S("assetType").Data().(string)
	assets, _ := jsonParsed.S("assets").Children()

	var header []string
	rows := make([][]string, len(assets))

	tBody := el.TableBody()
	switch assetType {

	case "alarms":
		var aType models.Alarm
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
		println("Asset Type not found in TableBuilder switch:")
		println(assetType)
	}

	buildTableRows(rows, tBody)

	table := el.Table(
		gr.CSS("table", "table-striped"),
		gr.Style("width", "100%"),
		el.TableHead(el.TableRow(buildTableHeader(header)...)))

	tBody.Modify(table)

	return table
}

func buildTableRows(rows [][]string, tBody *gr.Element) {
	for _, row := range rows {
		var rowEls = make([]gr.Modifier, len(row))

		for ri, rowName := range row {
			rowEls[ri] = el.TableData(gr.Text(rowName))
		}

		tr := el.TableRow(rowEls...)
		tr.Modify(tBody)
	}
}

func buildTableHeader(header []string) []gr.Modifier {
	var tableHeader = make([]gr.Modifier, len(header))
	for i, head := range header {
		tableHeader[i] = el.TableHeader(gr.Text(head))
	}
	return tableHeader
}
