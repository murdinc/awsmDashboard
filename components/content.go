package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
)

type Content struct {
	*gr.This
}

func (c Content) Render() gr.Component {

	// Table placeholder
	response := el.Div()

	elem := el.Div(gr.CSS("content-wrapper"),
		el.Div(gr.CSS("content"),
			el.Header1(gr.Text(c.Props().String("Title"))),
			response,
		),
	)

	if assets := c.State().Interface("assetList"); assets != nil {
		table := TableBuilder(assets)
		table.Modify(response)
	} else {
		gr.Text("loading...").Modify(response)
	}

	return elem
}

// Implements the ComponentDidMount interface
func (c Content) ComponentDidMount() {
	resp, err := QueryAPI("//localhost:8081/api" + c.Props().String("ApiEndpoint"))
	if err != nil {
		panic(err)
		// TODO handle query failures
	}

	c.SetState(gr.State{"assetList": resp})
}

// Implements the ComponentWillUnmount interface
func (g Content) ComponentWillUnmount() {
	// TODO Handle query cancel
}

// Implements the ShouldComponentUpdate interface.
func (g Content) ShouldComponentUpdate(this *gr.This, next gr.Cops) bool {
	return g.State().HasChanged(next.State, "assetList")
}
