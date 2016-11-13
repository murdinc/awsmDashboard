package components

import (
	"time"

	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
)

type ClassDropdownMenu struct {
	*gr.This
}

func (d ClassDropdownMenu) triggerRefresh(event *gr.Event) {
	d.SetState(gr.State{"nonce": time.Now()})
}

func (d ClassDropdownMenu) Render() gr.Component {

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

	//el.ListItem(el.Anchor(gr.Data("toggle", "modal"), gr.Data("target", "#new-asset-modal"), gr.Text("New "+pageType))).Modify(dropdownMenu) // New Asset
	el.ListItem(el.Anchor(gr.Data("toggle", "modal"), gr.Data("target", "#new-class-modal"), gr.Text("New Class"))).Modify(dropdownMenu)   // New Class
	el.ListItem(el.Anchor(gr.Data("toggle", "modal"), gr.Data("target", "#edit-class-modal"), gr.Text("Edit Class"))).Modify(dropdownMenu) // Edit Classes

	// New Asset
	// TODO
	/*
	   gr.New(&Modal{}).CreateElement(gr.Props{"id": "new-asset-modal", "title": "New " + pageType},
	   		gr.New(&NewAsset{}).CreateElement(gr.Props{"classList": classList, "apiType": apiType}),
	   	).Modify(dropdown)
	*/

	// New Class
	gr.New(&Modal{}).CreateElement(gr.Props{"id": "new-class-modal", "title": "New " + pageType + " Class"},
		gr.New(&NewClass{}).CreateElement(gr.Props{"apiType": apiType}),
	).Modify(dropdown)

	// Edit Class
	gr.New(&Modal{}).CreateElement(gr.Props{"id": "edit-class-modal", "title": "Edit " + pageType + " Classes"},
		gr.New(&EditClass{}).CreateElement(gr.Props{"apiType": apiType}),
	).Modify(dropdown)

	dropdownMenu.Modify(dropdown)

	return dropdown
}
