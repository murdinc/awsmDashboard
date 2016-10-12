package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
)

type Dropdown struct {
	*gr.This
	DropdownOptions []DropdownOption
}

type DropdownOption struct {
	Id   string
	Name string
}

func (d Dropdown) Render() gr.Component {

	if len(d.DropdownOptions) < 1 {
		return nil
	}

	dropdown := el.Div(
		gr.CSS("btn-group", "dropdown"),
		el.Button(
			gr.CSS("btn", "btn-primary", "btn-xs", "dropdown-toggle"),
			el.Italic(gr.CSS("fa", "fa-gear")),
			gr.Data("toggle", "dropdown"),
			gr.Aria("haspopup", "true"),
			gr.Aria("expanded", "true"),
		),
	)

	dropdownMenu := el.UnorderedList(
		gr.CSS("dropdown-menu"),
	)

	for _, option := range d.DropdownOptions {
		el.ListItem(el.Anchor(gr.Data("toggle", "modal"), gr.Data("target", "#"+option.Id), gr.Text(option.Name))).Modify(dropdownMenu)
		gr.New(&Modal{Id: option.Id, Title: option.Name}).CreateElement(gr.Props{}).Modify(dropdown)
	}

	dropdownMenu.Modify(dropdown)

	return dropdown
}
