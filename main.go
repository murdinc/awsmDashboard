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
			Route:   "/",
			ApiType: "instances",
			Type:    "Dashboard",
		},
		"Instances": components.Page{
			Route:   "/instances",
			ApiType: "instances",
			Type:    "Instance",
		},
		"Volumes": components.Page{
			Route:   "/volumes",
			ApiType: "volumes",
			Type:    "Volume",
		},
		"Images": components.Page{
			Route:   "/images",
			ApiType: "images",
			Type:    "Image",
		},
		"Snapshots": components.Page{
			Route:   "/snapshots",
			ApiType: "snapshots",
			Type:    "Snapshot",
		},
		"Vpcs": components.Page{
			Route:   "/vpcs",
			ApiType: "vpcs",
			Type:    "Vpc",
		},
		"Subnets": components.Page{
			Route:   "/subnets",
			ApiType: "subnets",
			Type:    "Subnet",
		},
		"Security Groups": components.Page{
			Route:   "/securitygroups",
			ApiType: "securitygroups",
			Type:    "SecurityGroup",
		},
		"Addresses": components.Page{
			Route:   "/addresses",
			ApiType: "addresses",
			Type:    "Address",
		},
		"Alarms": components.Page{
			Route:   "/alarms",
			ApiType: "alarms",
			Type:    "Alarm",
		},
		"Keypairs": components.Page{
			Route:   "/keypairs",
			ApiType: "keypairs",
			Type:    "Keypair",
		},
		"Launch Configurations": components.Page{
			Route:   "/launchconfigurations",
			ApiType: "launchconfigurations",
			Type:    "Launch Configuration",
		},
		"Autoscale Groups": components.Page{
			Route:   "/autoscalegroups",
			ApiType: "autoscalegroups",
			Type:    "Autoscale Group",
		},
		"Load Balancers": components.Page{
			Route:   "/loadbalancers",
			ApiType: "loadbalancers",
			Type:    "Load Balancer",
		},
		"Scaling Policies": components.Page{
			Route:   "/scalingpolicies",
			ApiType: "scalingpolicies",
			Type:    "Scaling Policy",
		},
		"SimpleDB Domains": components.Page{
			Route:   "/simpledbdomains",
			ApiType: "simpledbdomains",
			Type:    "SimpleDB Domain",
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
