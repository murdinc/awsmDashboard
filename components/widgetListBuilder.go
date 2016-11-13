package components

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
	"github.com/murdinc/awsm/config"
)

func WidgetListBuilder(wl interface{}, onClick func(string)) *gr.Element {

	widgetList := wl.([]byte)
	jsonParsed, _ := gabs.ParseJSON(widgetList)
	widgets, _ := jsonParsed.S("widgets").ChildrenMap()

	widgetListGroup := el.Div(
		el.Div(
			gr.CSS("list-group"),
		),
	)

	for widgetName, widget := range widgets {

		widgetType := widget.S("widgetType").Data().(string)

		switch widgetType {
		case "rss":
			var wType config.Widget

			wJson := widget.Bytes()
			json.Unmarshal(wJson, &wType)
			keys, values := config.ExtractAwsmWidget(wType)
			buildWidgetButton(widgetName, keys, values, widgetListGroup, onClick)

		case "events":
			/*
				var wType config.Widget

				wJson := widget.Bytes()
				json.Unmarshal(wJson, &wType)
				keys, values := config.ExtractAwsmWidget(wType)
				buildWidgetButton(widgetName, keys, values, widgetListGroup, onClick)
			*/
		default:
			println("Widget Type not found in WidgetListBuilder switch:")
			println(widgetType)

		}

	}

	return widgetListGroup
}

func buildWidgetButton(name string, keys []string, values []string, widgetListGroup *gr.Element, onClick func(string)) {

	clickListener := func(event *gr.Event) {
		onClick(name)
	}

	button := el.Button(
		attr.Type("button"),
		gr.CSS("list-group-item"),
		attr.ID(name),
		evt.Click(clickListener),
		el.Header5(
			gr.CSS("list-group-item-heading"),
			gr.Text(name),
		),
	)

	description := el.Div(
		gr.CSS("list-group-item-text"),
	)

	left := el.DescriptionList(gr.CSS("dl-horizontal"))
	right := el.DescriptionList(gr.CSS("dl-horizontal"))

	keyLen := len(keys)

	for i, key := range keys {
		if i < (keyLen / 2) {
			el.DefinitionTerm(gr.Text(key)).Modify(left)
			el.Description(gr.Text(values[i])).Modify(left)
		} else {
			el.DefinitionTerm(gr.Text(key)).Modify(right)
			el.Description(gr.Text(values[i])).Modify(right)
		}
	}

	el.Div(
		gr.CSS("row"),
		el.Div(gr.CSS("col-sm-6"),
			left),
		el.Div(gr.CSS("col-sm-6"),
			right),
	).Modify(description)

	description.Modify(button)
	button.Modify(widgetListGroup)

}
