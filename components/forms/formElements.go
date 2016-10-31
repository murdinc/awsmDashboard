package forms

import (
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
)

func textField(name, id string, state *gr.State, storeFunc func(*gr.Event)) *gr.Element {
	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		el.Input(
			attr.Type("text"),
			attr.ClassName("form-control"),
			attr.ID(id),
			attr.Placeholder(name),
			attr.Value(state.String(id)),
			evt.Change(storeFunc),
		),
	)
}

func numberField(name, id string, state *gr.State, storeFunc func(*gr.Event)) *gr.Element {
	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		el.Input(
			attr.Type("number"),
			attr.ClassName("form-control"),
			attr.ID(id),
			attr.Placeholder(name),
			attr.Value(state.Int(id)),
			evt.Change(storeFunc),
		),
	)
}

func textArea(name, id string, state *gr.State, storeFunc func(*gr.Event)) *gr.Element {
	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		el.TextArea(
			attr.ClassName("form-control"),
			attr.ID(id),
			attr.Placeholder(name),
			evt.Change(storeFunc),
			attr.Value(state.String(id)),
			attr.Rows(5),
		),
	)
}

func checkbox(name, id string, state *gr.State, storeFunc func(*gr.Event)) *gr.Element {

	label := "disabled"
	var checked gr.Modifier
	if state.Bool(id) {
		label = "enabled"
		checked = attr.Checked(true)
	}

	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		el.Div(
			gr.CSS("checkbox"),
			el.Label(
				el.Input(
					attr.Value(""), // to stop the warning about a uncontrolled components
					attr.Type("checkbox"),
					attr.ID(id),
					checked,
					evt.Change(storeFunc).StopPropagation(),
				),
				gr.Text(label),
			),
		),
	)
}

func selectOne(name, id string, options []string, state *gr.State, storeSelect func(string, interface{})) *gr.Element {
	opts := make([]interface{}, len(options))
	for i, option := range options {
		opts[i] = map[string]string{
			"value": option,
			"label": option,
		}
	}

	var value interface{}
	if selected, ok := state.Interface(id).(interface{}); ok {
		value = selected
	}

	onChange := func(vals interface{}) {
		storeSelect(id, vals)
	}

	reactSelect := gr.FromGlobal("Select")
	reactSelectElem := reactSelect.CreateElement(gr.Props{
		"name":               name,
		"value":              value,
		"options":            opts,
		"onChange":           onChange,
		"clearable":          false,
		"scrollMenuIntoView": false,
	})

	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		reactSelectElem,
	)
}

func selectMultiple(name, id string, options []string, state *gr.State, storeSelect func(string, interface{})) *gr.Element {
	opts := make([]interface{}, len(options))
	for i, option := range options {
		opts[i] = map[string]string{
			"value": option,
			"label": option,
		}
	}

	var value []interface{}
	if selected, ok := state.Interface(id).([]interface{}); ok {
		value = selected
	}

	onChange := func(vals interface{}) {
		storeSelect(id, vals)
	}

	reactSelect := gr.FromGlobal("Select")
	reactSelectElem := reactSelect.CreateElement(gr.Props{
		"name":               name,
		"value":              value,
		"options":            opts,
		"onChange":           onChange,
		"multi":              true,
		"scrollMenuIntoView": false,
	})

	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		reactSelectElem,
	)
}
