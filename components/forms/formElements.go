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

func TextField(name, key string, v interface{}, storeFunc func(*gr.Event)) *gr.Element {

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
			attr.Name(key),
			attr.Placeholder(name),
			attr.Value(value),
			evt.Change(storeFunc),
		),
	)
}

func NumberField(name, key string, value interface{}, storeFunc func(*gr.Event)) *gr.Element {

	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		el.Input(
			attr.Type("number"),
			attr.ClassName("form-control"),
			attr.Name(key),
			attr.Placeholder(name),
			attr.Value(value),
			evt.Change(storeFunc),
		),
	)
}

func TextArea(name, key string, value string, storeFunc func(*gr.Event)) *gr.Element {
	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		el.TextArea(
			attr.ClassName("form-control"),
			attr.Name(key),
			attr.Placeholder(name),
			evt.Change(storeFunc),
			attr.Value(value),
			attr.Rows(8),
		),
	)
}

func Checkbox(name, key string, value bool, storeFunc func(*gr.Event)) *gr.Element {

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
					attr.Name(key),
					checked,
					evt.Change(storeFunc).StopPropagation(),
				),
				gr.Text(label),
			),
		),
	)
}

func Toggle(falseName, trueName, key string, value interface{}, storeFunc func(*gr.Event)) *gr.Element {

	valBool, ok := value.(bool)
	if !ok {
		valBool = false
	}

	label := falseName
	var checked gr.Modifier
	if valBool {
		label = trueName
		checked = attr.Checked(true)
	}

	return el.Div(
		gr.CSS("form-group"),
		el.Div(
			el.Label(
				gr.Text(label),
			),
		),
		el.Div(
			gr.CSS("switch"),
			el.Label(
				el.Input(
					attr.Value(""), // to stop the warning about a uncontrolled components
					attr.Type("checkbox"),
					attr.Name(key),
					checked,
					evt.Change(storeFunc).StopPropagation(),
				),
				el.Div(gr.CSS("slider")),
			),
		),
	)
}

func SelectOne(name, key string, options []string, value interface{}, storeSelect func(string, interface{})) *gr.Element {
	opts := make([]interface{}, len(options))
	for i, option := range options {
		opts[i] = map[string]string{
			"value": option,
			"label": option,
		}
	}

	onChange := func(vals interface{}) {
		storeSelect(key, vals)
	}

	reactSelect := gr.FromGlobal("Select")
	reactSelectElem := reactSelect.CreateElement(gr.Props{
		"name":               key,
		"value":              value,
		"options":            opts,
		"onChange":           onChange,
		"clearable":          true,
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

func SelectOneMeta(name, key string, options []string, optionsMeta map[string]string, value interface{}, storeSelect func(string, interface{})) *gr.Element {
	opts := make([]interface{}, len(options))
	for i, option := range options {
		opts[i] = map[string]string{
			"value": option,
			"label": option + " - " + optionsMeta[option],
		}
	}

	onChange := func(vals interface{}) {
		storeSelect(key, vals)
	}

	//reactSelect := gr.FromGlobal("Select")
	reactSelectElem := reactSelect.CreateElement(gr.Props{
		"name":               name,
		"value":              value,
		"options":            opts,
		"onChange":           onChange,
		"clearable":          true,
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

func CreateableSelectMeta(name, key string, options []string, optionsMeta map[string]string, value interface{}, storeSelect func(string, interface{})) *gr.Element {

	selStr, ok := value.(string)
	if !ok {
		selStr = ""
	}
	existing := false

	opts := make([]interface{}, len(options))
	for i, option := range options {
		opts[i] = map[string]string{
			"value": option,
			"label": option + " - " + optionsMeta[option],
		}
		if option == selStr {
			existing = true
		}
	}

	// Creatable API seems to be a bit in limbo currently, so doing this to account for the experienced wonkiness
	if !existing {
		newVal := make(map[string]string)
		newVal["value"] = selStr
		newVal["label"] = selStr
		opts = append(opts, newVal)
	}

	onChange := func(vals interface{}) {
		storeSelect(key, vals)
	}

	reactSelectElem := reactCreatableSelect.CreateElement(gr.Props{
		"name":               name,
		"value":              value,
		"options":            opts,
		"onChange":           onChange,
		"multi":              false,
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

func SelectMultiple(name, key string, options []string, value interface{}, storeSelect func(string, interface{})) *gr.Element {
	opts := make([]interface{}, len(options))
	for i, option := range options {
		opts[i] = map[string]string{
			"value": option,
			"label": option,
		}
	}

	onChange := func(vals interface{}) {
		storeSelect(key, vals)
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

func CreateableSelectMultiple(name, key string, options []string, s interface{}, storeSelect func(string, interface{})) *gr.Element {

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
		storeSelect(key, vals)
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

func DisableEnter(event *gr.Event) {
	key := event.Get("key").String()
	keyCode := event.Get("keyCode").String()

	if (keyCode == "13" || key == "Enter" || key == "Return") && event.Target().Get("type").String() != "textarea" {
		event.Object.Call("preventDefault")
	}
}

func CaptureEnter(callback func(*gr.Event)) func(*gr.Event) {
	return func(event *gr.Event) {
		key := event.Get("key").String()
		keyCode := event.Get("keyCode").String()

		if keyCode == "13" || key == "Enter" || key == "Return" {
			event.Object.Call("preventDefault")
			callback(event)
		}
	}
}
