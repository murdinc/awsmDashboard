package components

import (
	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/murdinc/awsmDashboard/components/forms"
)

func EditWidgetFormBuilder(widgetsBytes []byte) (*gr.ReactComponent, *gabs.Container) {

	jsonParsed, _ := gabs.ParseJSON(widgetsBytes)
	widgetType := jsonParsed.S("widgetType").Data().(string)

	switch widgetType {

	case "rss":
		return gr.New(&forms.RSSWidgetForm{}), jsonParsed

		/*
			case "events":
				return gr.New(&forms.AlarmWidgetForm{}), jsonParsed
		*/

	default:
		println("Widget Type not found in EditWidgetFormBuilder switch:")
		println(widgetType)
	}

	return nil, nil
}

func NewWidgetFormBuilder(widgetType string) *gr.ReactComponent {

	switch widgetType {

	case "rss":
		return gr.New(&forms.RSSWidgetForm{}) // TODO

	default:
		println("Widget Type not found in NewWidgetFormBuilder switch:")
		println(widgetType)
	}

	return nil
}
