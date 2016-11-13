package components

import (
	"time"

	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
)

type WidgetDropdownMenu struct {
	*gr.This
}

func (d WidgetDropdownMenu) triggerRefresh(event *gr.Event) {
	d.SetState(gr.State{"nonce": time.Now()})
}

func (d WidgetDropdownMenu) Render() gr.Component {

	//state := d.State()
	props := d.Props()
	apiType := props.String("apiType")
	pageType := props.String("type")

	dropdown := el.Div(
		gr.CSS("btn-group", "dropdown"),
		el.Button(
			evt.Click(d.triggerRefresh), // TODO, handle this better so we aren't throwing away api hits
			gr.CSS("btn", "btn-primary", "btn-xs", "dropdown-toggle"),
			el.Italic(gr.CSS("fa", "fa-gear")),
			gr.Data("toggle", "dropdown"),
		),
	)

	dropdownMenu := el.UnorderedList(
		gr.CSS("dropdown-menu"),
	)

	el.ListItem(el.Anchor(gr.Data("toggle", "modal"), gr.Data("target", "#new-widget-modal"), gr.Text("New Widget"))).Modify(dropdownMenu)     // New Widget
	el.ListItem(el.Anchor(gr.Data("toggle", "modal"), gr.Data("target", "#edit-widgets-modal"), gr.Text("Edit Widgets"))).Modify(dropdownMenu) // Edit Widgets

	// New Widget
	gr.New(&Modal{}).CreateElement(gr.Props{"id": "new-widget-modal", "title": "New " + pageType + " Widget"},
		gr.New(&NewWidget{}).CreateElement(gr.Props{"apiType": apiType}),
	).Modify(dropdown)

	// Edit Widgets
	gr.New(&Modal{}).CreateElement(gr.Props{"id": "edit-widgets-modal", "title": "Edit " + pageType + " Widgets"},
		gr.New(&EditWidgets{}).CreateElement(gr.Props{"apiType": apiType}),
	).Modify(dropdown)

	dropdownMenu.Modify(dropdown)

	return dropdown
}
