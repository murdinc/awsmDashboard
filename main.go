package main

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/grouter"
	"github.com/gopherjs/gopherjs/js"
	"github.com/murdinc/awsmDashboard/components"
)

var (
	reactRouter = js.Global.Get("ReactRouter")

	routerHistory = grouter.History{Object: reactRouter.Get("browserHistory")}

	dashboard = gr.New(&components.Layout{Title: "Dashboard"}, gr.Apply(grouter.WithRouter))
	instances = gr.New(&components.Layout{Title: "Instances"}, gr.Apply(grouter.WithRouter))
	images    = gr.New(&components.Layout{Title: "Images"}, gr.Apply(grouter.WithRouter))
	volumes   = gr.New(&components.Layout{Title: "Volumes"}, gr.Apply(grouter.WithRouter))

	// WithRouter makes this.props.router happen.
	appComponent = gr.New(new(app), gr.Apply(grouter.WithRouter))

	router = grouter.New("/", appComponent, grouter.WithHistory(routerHistory)).With(
		//router = grouter.New("/", appComponent).With(
		grouter.NewIndexRoute(grouter.Components{"layout": dashboard}),
		grouter.NewRoute("/instances", grouter.Components{"layout": instances}),
		grouter.NewRoute("/volumes", grouter.Components{"layout": volumes}),
		grouter.NewRoute("/images", grouter.Components{"layout": images}),
	)
)

func main() {
	mainComponent := gr.New(gr.NewSimpleRenderer(router))
	mainComponent.Render("react", gr.Props{})
}

type app struct {
	*gr.This
	Nav components.Nav
}

// Implements the Renderer interface.
func (a app) Render() gr.Component {
	return el.Div(
		a.Component("layout"),
	)
}
