package helpers

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
)

func ErrorElem(errStr string) *gr.Element {
	if errStr != "" {
		return el.Div(
			gr.CSS("alert", "alert-danger", "fade", "in"),
			el.Strong(gr.Text("Error! ")),
			gr.Text(errStr),
		)
	}

	return el.Div()
}

func SuccessElem(successStr string) *gr.Element {
	if successStr != "" {
		return el.Div(
			gr.CSS("alert", "alert-success", "fade", "in"),
			el.Strong(gr.Text("Success! ")),
			gr.Text(successStr),
		)
	}

	return el.Div()
}
