package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
)

type ButtonList struct {
	*gr.This
	Buttons []Button
}

type Button struct {
	Name        string
	Description string
}

func (b ButtonList) Render() gr.Component {

	buttonList := el.Div(
		gr.CSS("list-group"),
	)

	for _, button := range b.Buttons {
		el.Button(
			attr.Type("button"),
			gr.CSS("list-group-item"),
			el.Header5(
				gr.CSS("list-group-item-heading"),
				gr.Text(button.Name),
			),
			el.Paragraph(
				gr.CSS("list-group-item-text"),
				gr.Text(button.Description),
			),
		).Modify(buttonList)
	}

	return buttonList
}
