package forms

import (
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
)

var (
	reactSelect          = gr.FromGlobal("Select")
	reactCreatableSelect = gr.FromGlobal("Select", "Creatable")
)

func textField(name, id string, v interface{}, storeFunc func(*gr.Event)) *gr.Element {

	var value string

	valueStr, ok := v.(string)
	if ok {
		value = valueStr
	}

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
			attr.Value(value),
			evt.Change(storeFunc),
		),
	)
}

func numberField(name, id string, value interface{}, storeFunc func(*gr.Event)) *gr.Element {

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
			attr.Value(value),
			evt.Change(storeFunc),
		),
	)
}

func textArea(name, id string, value string, storeFunc func(*gr.Event)) *gr.Element {
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
			attr.Value(value),
			attr.Rows(5),
		),
	)
}

func checkbox(name, id string, value bool, storeFunc func(*gr.Event)) *gr.Element {

	label := "disabled"
	var checked gr.Modifier
	if value {
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

func selectOne(name, id string, options []string, value interface{}, storeSelect func(string, interface{})) *gr.Element {
	opts := make([]interface{}, len(options))
	for i, option := range options {
		opts[i] = map[string]string{
			"value": option,
			"label": option,
		}
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

func selectMultiple(name, id string, options []string, value interface{}, storeSelect func(string, interface{})) *gr.Element {
	opts := make([]interface{}, len(options))
	for i, option := range options {
		opts[i] = map[string]string{
			"value": option,
			"label": option,
		}
	}

	onChange := func(vals interface{}) {
		storeSelect(id, vals)
	}

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

func createableSelectMultiple(name, id string, options []string, s interface{}, storeSelect func(string, interface{})) *gr.Element {

	var selected []interface{}
	selectedSlice, ok := s.([]interface{})
	if ok {
		selected = selectedSlice
	}

	opts := make([]interface{}, len(options))
	for i, option := range options {
		opts[i] = map[string]string{
			"value": option,
			"label": option,
		}
	}

	// Creatable API seems to be a bit in limbo currently, so doing this to account for the experienced wonkiness
	var value []interface{}
	for _, sel := range selected {
		selStr, ok := sel.(string)
		if ok {
			newVal := make(map[string]string)
			newVal["value"] = selStr
			newVal["label"] = selStr
			value = append(value, newVal)
		}
	}

	onChange := func(vals interface{}) {
		storeSelect(id, vals)
	}

	reactSelectElem := reactCreatableSelect.CreateElement(gr.Props{
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
