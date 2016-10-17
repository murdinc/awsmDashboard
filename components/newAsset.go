package components

import (
	"fmt"

	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsmDashboard/helpers"
)

type NewAsset struct {
	*gr.This
}

// Implements the StateInitializer interface.
func (n NewAsset) GetInitialState() gr.State {
	return gr.State{"selectedClass": "", "querying": false, "error": "", "classData": nil}
}

func (n *NewAsset) selectClass(name string) {
	go func() {
		if endpoint := n.Props().String("classEndpoint"); endpoint != "" {
			n.SetState(gr.State{"querying": true})
			resp, err := helpers.QueryAPI("//localhost:8081/api" + endpoint + "/name/" + name)
			if !n.IsMounted() {
				return
			}
			if err != nil {
				n.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
				return
			}
			n.SetState(gr.State{"querying": false, "classData": resp})
		}
	}()
	n.SetState(gr.State{"selectedClass": name})
}

func (n NewAsset) Render() gr.Component {

	state := n.State()
	props := n.Props()

	// Class List placeholder
	response := el.Div()

	if state.Bool("querying") {
		gr.Text("Loading...").Modify(response)
	} else if errStr := state.String("error"); errStr != "" {
		gr.Text(errStr).Modify(response)
	} else {

		if state.String("selectedClass") == "" {

			// STEP 1

			if classes := props.Interface("classList"); classes != nil {
				classList := ClassListBuilder(classes, n.selectClass) // Build the class list
				classList.Modify(response)
			} else {
				gr.Text("Nothing here!").Modify(response)
			}

		} else if state.Interface("classData") != nil {

			// STEP 2

			gr.Text(state.String("selectedClass")).Modify(response)

			println(state.Interface("classData"))

		}
	}

	return response
}
