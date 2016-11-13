package components

import (
	"fmt"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
	"github.com/murdinc/awsmDashboard/helpers"
)

type EditWidgets struct {
	*gr.This
}

// Implements the StateInitializer interface
func (e EditWidgets) GetInitialState() gr.State {
	return gr.State{"step": 1, "selectedWidget": "", "querying": false, "error": "", "widgetData": nil}
}

func (e EditWidgets) ComponentWillMount() {
	e.getWidgetsList()
}

func (e EditWidgets) doneButton(*gr.Event) {
	e.SetState(gr.State{"success": ""})
	hideAllModals()
}

func (e EditWidgets) Render() gr.Component {

	state := e.State()
	props := e.Props()

	// Widget List placeholder
	response := el.Div()

	if state.Bool("querying") {
		gr.Text("Loading...").Modify(response)

	} else if errStr := state.String("error"); errStr != "" {
		helpers.ErrorElem(errStr).Modify(response)
	} else {

		if state.Int("step") == 1 {

			// STEP 1

			if widgets := state.Interface("widgetList"); widgets != nil {
				widgetList := WidgetListBuilder(widgets, e.selectWidget) // Build the widget list
				widgetList.Modify(response)
			} else {

				helpers.ErrorElem("No existing " + props.String("apiType") + " widgets found!").Modify(response)

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

			widgetForm, widgetJson := EditWidgetFormBuilder(state.Interface("widgetData").([]byte))

			widgetForm.CreateElement(
				gr.Props{
					"widgetName":    e.State().String("selectedWidget"),
					"widget":        widgetJson.Bytes(),
					"backButton":    e.stepTwoBack,
					"apiType":       props.String("apiType"),
					"hasDelete":     true,
					"hideAllModals": hideAllModals,
				},
			).Modify(response)
		}
	}

	return response
}

func (e EditWidgets) getWidgetsList() {
	go func() {
		if apiType := e.Props().String("apiType"); apiType != "" {
			e.SetState(gr.State{"querying": true})
			endpoint := "//localhost:8081/api/" + apiType + "/widgets"
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
				println("no existing " + e.Props().String("apiType") + " widgets found")
				e.SetState(gr.State{"querying": false})
				return
			}

			e.SetState(gr.State{"querying": false, "widgetList": resp})
		}
	}()
}

func (e *EditWidgets) selectWidget(name string) {
	state := e.State()

	e.SetState(gr.State{"querying": true})

	if wl := state.Interface("widgetList"); wl != nil {
		widgetList, _ := gabs.ParseJSON(wl.([]byte))
		widgetData := widgetList.S("widgets").S(name).Bytes()
		e.SetState(gr.State{"querying": false, "step": 2, "selectedWidget": name, "widgetData": widgetData})
	}
}

func (e EditWidgets) stepTwoBack() {
	e.getWidgetsList()
	e.SetState(gr.State{"step": 1})
}
