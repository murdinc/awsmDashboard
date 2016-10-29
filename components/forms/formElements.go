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
			attr.ClassName("form-control"),
			attr.ID(id),
			attr.Placeholder(name),
			attr.Value(state.String(id)),
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
	checked := attr.Target(gr.CSS(""))
	if state.Bool(id) {
		label = "enabled"
		checked = attr.Checked(gr.CSS(""))
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
					evt.Change(storeFunc),
				),
				gr.Text(label),
			),
		),
	)
}

func selectOne(name, id string, options []string, state *gr.State, storeFunc func(*gr.Event)) *gr.Element {

	selectOpts := el.Select(
		attr.ClassName("form-control"),
		attr.ID(id),
		evt.Change(storeFunc),
		attr.Value(state.String(id)),
		el.Option(
			//attr.ID("nil"),
			//attr.Value("nil"),
			//attr.Disabled("true"),
			gr.Text("Select "+name),
		),
	)

	for _, option := range options {
		opt := el.Option(
			attr.ID(option),
			attr.Value(option),
			gr.Text(option),
		)

		opt.Modify(selectOpts)
	}

	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		selectOpts,
	)
}

func selectMultiple(name, id string, options []string, state *gr.State, storeFunc func(*gr.Event)) *gr.Element {

	selectOpts := el.Select(
		attr.Multiple(attr.ClassName("")),
		attr.ClassName("form-control"),
		attr.ID(id),
		evt.Change(storeFunc),
	)

	if selected, ok := state.Interface(id).([]interface{}); ok {
		attr.Value(selected).Modify(selectOpts)
	}

	for _, option := range options {
		opt := el.Option(
			attr.ID(option),
			attr.Value(option),
			gr.Text(option),
		)

		opt.Modify(selectOpts)
	}

	return el.Div(
		gr.CSS("form-group"),
		el.Label(
			gr.Text(name),
		),
		selectOpts,
	)
}
