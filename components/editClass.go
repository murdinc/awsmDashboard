package components

import (
	"fmt"

	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsmDashboard/helpers"
)

type EditClass struct {
	*gr.This
}

// Implements the StateInitializer interface.
func (e EditClass) GetInitialState() gr.State {
	return gr.State{"selectedClass": "", "querying": false, "error": "", "classData": nil}
}

func (e *EditClass) selectClass(name string) {
	go func() {
		if endpoint := e.Props().String("classEndpoint"); endpoint != "" {
			e.SetState(gr.State{"querying": true})
			resp, err := helpers.QueryAPI("//localhost:8081/api" + endpoint + "/name/" + name)
			if !e.IsMounted() {
				return
			}
			if err != nil {
				e.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
				return
			}
			e.SetState(gr.State{"querying": false, "classData": resp})
		}
	}()
	e.SetState(gr.State{"selectedClass": name})
}

func (e EditClass) Render() gr.Component {

	state := e.State()
	props := e.Props()

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
				classList := ClassListBuilder(classes, e.selectClass) // Build the class list
				classList.Modify(response)
			} else {
				gr.Text("Nothing here!").Modify(response)
			}

		} else if state.Interface("classData") != nil {

			// STEP 2

			EditClassFormBuilder(state.Interface("classData")).Modify(response)

		}
	}

	return response
}
