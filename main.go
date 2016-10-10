package main

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/grouter"
	"github.com/gopherjs/gopherjs/js"
	"github.com/murdinc/awsmDashboard/components"
)

var (
	brand = "awsm"

	pages = components.Pages{
		"Dashboard": components.Page{
			Route:       "/",
			ApiEndpoint: "/assets/instances",
		},
		"Instances": components.Page{
			Route:       "/instances",
			ApiEndpoint: "/assets/instances",
		},
		"Volumes": components.Page{
			Route:       "/volumes",
			ApiEndpoint: "/assets/volumes",
		},
		"Images": components.Page{
			Route:       "/images",
			ApiEndpoint: "/assets/images",
		},
		"Snapshots": components.Page{
			Route:       "/snapshots",
			ApiEndpoint: "/assets/snapshots",
		},
		"Vpcs": components.Page{
			Route:       "/vpcs",
			ApiEndpoint: "/assets/vpcs",
		},
		"Subnets": components.Page{
			Route:       "/subnets",
			ApiEndpoint: "/assets/subnets",
		},
		"Security Groups": components.Page{
			Route:       "/securitygroups",
			ApiEndpoint: "/assets/securitygroups",
		},
		"Addresses": components.Page{
			Route:       "/addresses",
			ApiEndpoint: "/assets/addresses",
		},
		"Alarms": components.Page{
			Route:       "/alarms",
			ApiEndpoint: "/assets/alarms",
		},
		"Keypairs": components.Page{
			Route:       "/keypairs",
			ApiEndpoint: "/assets/keypairs",
		},
		"Launch Configurations": components.Page{
			Route:       "/launchconfigurations",
			ApiEndpoint: "/assets/launchconfigurations",
		},
		"Load Balancers": components.Page{
			Route:       "/loadbalancers",
			ApiEndpoint: "/assets/loadbalancers",
		},
		"Scaling Policies": components.Page{
			Route:       "/scalingpolicies",
			ApiEndpoint: "/assets/scalingpolicies",
		},
		"SimpleDB Domains": components.Page{
			Route:       "/simpledbdomains",
			ApiEndpoint: "/assets/simpledbdomains",
		},
	}

	reactRouter   = js.Global.Get("ReactRouter")
	routerHistory = grouter.History{Object: reactRouter.Get("browserHistory")}

	// WithRouter makes this.props.router happen.
	appComponent = gr.New(new(app), gr.Apply(grouter.WithRouter))
)

func main() {

	var routes []grouter.Route

	for name, page := range pages {
		switch page.Route {
		case "/":
			routes = append(routes,
				grouter.NewIndexRoute(
					grouter.Components{"page": gr.New(&components.Layout{Brand: brand, ActivePage: name, Pages: pages}, gr.Apply(grouter.WithRouter))}))
		default:
			routes = append(routes,
				grouter.NewRoute(page.Route,
					grouter.Components{"page": gr.New(&components.Layout{Brand: brand, ActivePage: name, Pages: pages}, gr.Apply(grouter.WithRouter))}))
		}
	}

	router := grouter.New("/", appComponent, grouter.WithHistory(routerHistory)).With(routes...)

	mainComponent := gr.New(gr.NewSimpleRenderer(router))
	mainComponent.Render("react", gr.Props{})
}

type app struct {
	*gr.This
}

// Implements the Renderer interface.
func (a app) Render() gr.Component {
	return el.Div(
		a.Component("page"),
	)
}
