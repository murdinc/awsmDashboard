package components

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
	"github.com/murdinc/awsmDashboard/helpers"
)

type EditClass struct {
	*gr.This
}

// Implements the StateInitializer interface
func (e EditClass) GetInitialState() gr.State {
	return gr.State{"step": 1, "selectedClass": "", "querying": false, "error": "", "classData": nil}
}

func (e EditClass) ComponentWillMount() {
	e.getClassList()
}

func (e EditClass) doneButton(*gr.Event) {
	e.SetState(gr.State{"success": ""})
	hideAllModals()
}

func (e EditClass) Render() gr.Component {

	state := e.State()
	props := e.Props()

	// Class List placeholder
	response := el.Div()

	if state.Bool("querying") {
		gr.Text("Loading...").Modify(response)

	} else if errStr := state.String("error"); errStr != "" {
		helpers.ErrorElem(errStr).Modify(response)
	} else {

		if state.Int("step") == 1 {

			// STEP 1

			if classes := state.Interface("classList"); classes != nil {
				classList := ClassListBuilder(classes, e.selectClass) // Build the class list
				classList.Modify(response)
			} else {

				helpers.ErrorElem("No existing " + props.String("apiType") + " classes found!").Modify(response)

				buttons := el.Div(
					gr.CSS("btn-toolbar"),
				)

				// Done
				el.Button(
					evt.Click(e.doneButton).PreventDefault(),
					gr.CSS("btn", "btn-primary"),
					gr.Text("Done"),
				).Modify(buttons)

				buttons.Modify(response)

			}

		} else if state.Int("step") == 2 {

			// STEP 2

			classForm, classJson := EditClassFormBuilder(state.Interface("classData").([]byte))

			classForm.CreateElement(
				gr.Props{
					"className":     classJson.S("className").Data().(string),
					"class":         classJson.S("class").Bytes(),
					"backButton":    e.stepTwoBack,
					"apiType":       props.String("apiType"),
					"hasDelete":     true,
					"hideAllModals": hideAllModals,
					"editClass":     true,
				},
			).Modify(response)
		}
	}

	return response
}

func (e EditClass) getClassList() {
	go func() {
		if apiType := e.Props().String("apiType"); apiType != "" {
			e.SetState(gr.State{"querying": true})
			endpoint := "//localhost:8081/api/classes/" + apiType
			resp, err := helpers.GetAPI(endpoint)
			if !e.IsMounted() {
				return
			}
			if err != nil {
				e.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
				return
			}

			respParsed, _ := gabs.ParseJSON(resp)

			success, ok := respParsed.S("success").Data().(bool)
			if !ok || !success {
				println("no existing " + e.Props().String("apiType") + " classes found")
				e.SetState(gr.State{"querying": false})
				return
			}

			e.SetState(gr.State{"querying": false, "classList": resp})
		}
	}()
}

func (e *EditClass) selectClass(name string) {
	e.SetState(gr.State{"querying": true})
	go func() {
		if apiType := e.Props().String("apiType"); apiType != "" {
			endpoint := "//localhost:8081/api/classes/" + apiType + "/name/" + name
			resp, err := helpers.GetAPI(endpoint)
			if !e.IsMounted() {
				return
			}
			if err != nil {
				e.SetState(gr.State{"querying": false, "error": fmt.Sprintf("Error while querying endpoint: %s", endpoint)})
				return
			}
			e.SetState(gr.State{"classData": resp})
		}
		e.SetState(gr.State{"querying": false, "step": 2, "selectedClass": name})
	}()
}

func (e EditClass) stepTwoBack() {
	e.getClassList()
	e.SetState(gr.State{"step": 1})
}
