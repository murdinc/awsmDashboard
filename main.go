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
			Route:      "/",
			ApiType:    "volumes",
			Type:       "Dashboard",
			HasClasses: true,
		},
		"Instances": components.Page{
			Route:      "/instances",
			ApiType:    "instances",
			Type:       "Instance",
			HasClasses: true,
		},
		"Volumes": components.Page{
			Route:      "/volumes",
			ApiType:    "volumes",
			Type:       "Volume",
			HasClasses: true,
		},
		"Images": components.Page{
			Route:      "/images",
			ApiType:    "images",
			Type:       "Image",
			HasClasses: true,
		},
		"Snapshots": components.Page{
			Route:      "/snapshots",
			ApiType:    "snapshots",
			Type:       "Snapshot",
			HasClasses: true,
		},
		"Vpcs": components.Page{
			Route:      "/vpcs",
			ApiType:    "vpcs",
			Type:       "Vpc",
			HasClasses: true,
		},
		"Subnets": components.Page{
			Route:      "/subnets",
			ApiType:    "subnets",
			Type:       "Subnet",
			HasClasses: true,
		},
		"Security Groups": components.Page{
			Route:      "/securitygroups",
			ApiType:    "securitygroups",
			Type:       "Security Group",
			HasClasses: true,
		},
		"Addresses": components.Page{
			Route:      "/addresses",
			ApiType:    "addresses",
			Type:       "Address",
			HasClasses: false,
		},
		"Alarms": components.Page{
			Route:      "/alarms",
			ApiType:    "alarms",
			Type:       "Alarm",
			HasClasses: true,
		},
		"Key Pairs": components.Page{
			Route:      "/keypairs",
			ApiType:    "keypairs",
			Type:       "Key Pair",
			HasClasses: true,
		},
		"Launch Configurations": components.Page{
			Route:      "/launchconfigurations",
			ApiType:    "launchconfigurations",
			Type:       "Launch Configuration",
			HasClasses: true,
		},
		"Autoscale Groups": components.Page{
			Route:      "/autoscalegroups",
			ApiType:    "autoscalegroups",
			Type:       "Autoscale Group",
			HasClasses: true,
		},
		"Load Balancers": components.Page{
			Route:      "/loadbalancers",
			ApiType:    "loadbalancers",
			Type:       "Load Balancer",
			HasClasses: true,
		},
		"Scaling Policies": components.Page{
			Route:      "/scalingpolicies",
			ApiType:    "scalingpolicies",
			Type:       "Scaling Policy",
			HasClasses: true,
		},
		"SimpleDB Domains": components.Page{
			Route:      "/simpledbdomains",
			ApiType:    "simpledbdomains",
			Type:       "SimpleDB Domain",
			HasClasses: false,
		},
	}

	reactRouter = js.Global.Get("ReactRouter")

	routerHistory = grouter.History{Object: reactRouter.Get("browserHistory")}
	appComponent  = gr.New(new(app))
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
