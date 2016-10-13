package components

import (
	"fmt"

	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsmDashboard/helpers"
)

type Content struct {
	*gr.This
	Page Page
}

// Implements the StateInitializer interface.
func (c Content) GetInitialState() gr.State {
	return gr.State{"querying": false, "error": nil, "assetList": nil}
}

func (c Content) Render() gr.Component {

	// Table placeholder
	response := el.Div()

	elem := el.Div(gr.CSS("content-wrapper"),
		el.Div(gr.CSS("content-header"),
			el.Header1(
				gr.Text(c.Props().String("activePage")+" "),
			),
			gr.New(&Dropdown{DropdownOptions: c.Page.DropdownOptions, ClassEndpoint: c.Page.ClassEndpoint}).CreateElement(gr.Props{}),
		),
		el.Div(gr.CSS("content"),
			response,
		),
	)

	if assets := c.State().Interface("assetList"); assets != nil {
		table := AssetTableBuilder(assets) // Build the table
		table.Modify(response)
	} else if c.State().Bool("querying") {
		gr.Text("Loading...").Modify(response)
	} else if errStr := c.State().Interface("error"); errStr != nil {
		gr.Text(errStr).Modify(response)
	} else {
		gr.Text("Nothing here!").Modify(response)
	}

	return elem
}

// Implements the ComponentDidMount interface
func (c Content) ComponentDidMount() {

	if endpoint := c.Page.AssetEndpoint; endpoint != "" {

		c.SetState(gr.State{"querying": true})

		resp, err := helpers.QueryAPI("//localhost:8081/api" + endpoint)
		if !c.IsMounted() {
			return
		}
		if err != nil {
			c.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
			return
		}

		c.SetState(gr.State{"querying": false, "assetList": resp})

	}
}

// Implements the ShouldComponentUpdate interface.
func (c Content) ShouldComponentUpdate(this *gr.This, next gr.Cops) bool {
	return c.State().HasChanged(next.State, "assetList") &&
		c.State().HasChanged(next.State, "querying") &&
		c.State().HasChanged(next.State, "error")
}

/*

	// Implements the ChildContext interface.
	func (c Content) GetChildContext() gr.Context {
		return gr.Context{}
	}

	// Implements the ComponentWillUpdate interface
	func (c Content) ComponentWillUpdate(next gr.Cops) {
		log.Println("ComponentWillUpdate")
		log.Println(c.Props().String("activePage"))
	}

	// Implements the ComponentWillReceiveProps interface
	func (c Content) ComponentWillReceiveProps(data gr.Cops) {
		log.Println("ComponentWillReceiveProps")
		log.Println(c.Props().String("activePage"))
	}

	// Implements the ComponentDidUpdate interface
	func (c Content) ComponentDidUpdate(data gr.Cops) {
		log.Println("ComponentDidUpdate")
		log.Println(c.Props().String("activePage"))
	}

	// Implements the ComponentWillMount interface
	func (c Content) ComponentWillMount() {
		log.Println("ComponentWillMount")
		log.Println(c.Props().String("activePage"))
	}

	// Implements the ComponentWillUnmount interface
	func (c Content) ComponentWillUnmount() {
		log.Println("ComponentWillUnmount")
		log.Println(c.Props().String("activePage"))
	}

*/
