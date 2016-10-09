package components

import (
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/bep/grouter"
)

type Layout struct {
	*gr.This
	Title   string
	Uri     string
	Content *gr.ReactComponent
}

// Implements the Renderer interface.
func (l Layout) Render() gr.Component {

	return el.Div(
		gr.CSS("main-wrapper"),

		// Nav
		gr.New(&Nav{Brand: "awsm"}).CreateElement(l.This.Props()),

		//Content
		gr.New(&Content{}).CreateElement(gr.Props{"Title": l.Title, "Uri": l.Uri}),
	)
}

func (l Layout) createLinkListItem(path, Title string) gr.Modifier {
	return el.ListItem(
		grouter.MarkIfActive(l.Props(), path),
		attr.Role("presentation"),
		grouter.Link(path, Title))
}

func (l Layout) onClick(event *gr.Event) {
	l.SetState(gr.State{"counter": l.State().Int("counter") + 1})
}

// Implements the ShouldComponentUpdate interface.
func (l Layout) ShouldComponentUpdate(
	next gr.Cops) bool {

	return l.State().HasChanged(next.State, "counter")
}
