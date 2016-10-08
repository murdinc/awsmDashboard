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

	tBody := el.TableBody()
	switch assetType {

	case "instances":

		for _, a := range assets {

			instanceJson := a.Bytes()

			var instance models.Instance
			json.Unmarshal(instanceJson, &instance)

			tr := el.TableRow(
				el.TableData(gr.Text(instance.Name)),
				el.TableData(gr.Text(instance.Region)),
				el.TableData(gr.Text(instance.Size)),
				el.TableData(gr.Text(instance.PrivateIp)),
				el.TableData(gr.Text(instance.PublicIp)),
				el.TableData(gr.Text(instance.State)),
			)

			tr.Modify(tBody)
		}

	default:
		println("default")
	}

	table := el.Table(
		gr.CSS("table", "table-striped"),
		gr.Style("width", "100%"),
		el.TableHead(el.TableRow(
			el.TableHeader(gr.Text("Name")),
			el.TableHeader(gr.Text("Region")),
			el.TableHeader(gr.Text("Size")),
			el.TableHeader(gr.Text("Private IP")),
			el.TableHeader(gr.Text("Public IP")),
			el.TableHeader(gr.Text("State")),
		)))

	tBody.Modify(table)

	return table
}
