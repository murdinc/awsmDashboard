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
			Route: "/",
			//AssetEndpoint: "/assets/instances",
			ClassEndpoint: "/classes/instances",
			PageType:      "Dashboard",
			ClassType:     "instances",
		},
		"Instances": components.Page{
			Route:         "/instances",
			AssetEndpoint: "/assets/instances",
			ClassEndpoint: "/classes/instances",
			PageType:      "Instance",
			ClassType:     "instances",
		},
		"Volumes": components.Page{
			Route:         "/volumes",
			AssetEndpoint: "/assets/volumes",
			ClassEndpoint: "/classes/volumes",
			PageType:      "Volume",
			ClassType:     "volumes",
		},
		"Images": components.Page{
			Route:         "/images",
			AssetEndpoint: "/assets/images",
			ClassEndpoint: "/classes/images",
			PageType:      "Image",
			ClassType:     "images",
		},
		"Snapshots": components.Page{
			Route:         "/snapshots",
			AssetEndpoint: "/assets/snapshots",
			ClassEndpoint: "/classes/snapshots",
			PageType:      "Snapshot",
			ClassType:     "snapshots",
		},
		"Vpcs": components.Page{
			Route:         "/vpcs",
			AssetEndpoint: "/assets/vpcs",
			ClassEndpoint: "/classes/vpcs",
			PageType:      "Vpc",
			ClassType:     "vpcs",
		},
		"Subnets": components.Page{
			Route:         "/subnets",
			AssetEndpoint: "/assets/subnets",
			ClassEndpoint: "/classes/subnets",
			PageType:      "Subnet",
			ClassType:     "subnets",
		},
		"Security Groups": components.Page{
			Route:         "/securitygroups",
			AssetEndpoint: "/assets/securitygroups",
			ClassEndpoint: "/classes/securitygroups",
			PageType:      "Security Group",
			ClassType:     "securitygroups",
		},
		"Addresses": components.Page{
			Route:         "/addresses",
			AssetEndpoint: "/assets/addresses",
			ClassEndpoint: "/classes/addresses",
			PageType:      "Address",
			ClassType:     "addresses",
		},
		"Alarms": components.Page{
			Route:         "/alarms",
			AssetEndpoint: "/assets/alarms",
			ClassEndpoint: "/classes/alarms",
			PageType:      "Alarm",
			ClassType:     "alarms",
		},
		"Keypairs": components.Page{
			Route:         "/keypairs",
			AssetEndpoint: "/assets/keypairs",
			ClassEndpoint: "/classes/keypairs",
			PageType:      "KeyPair",
			ClassType:     "keypairs",
		},
		"Launch Configurations": components.Page{
			Route:         "/launchconfigurations",
			AssetEndpoint: "/assets/launchconfigurations",
			ClassEndpoint: "/classes/launchconfigurations",
			PageType:      "Launch Configuration",
			ClassType:     "launchconfigurations",
		},
		"Autoscale Groups": components.Page{
			Route:         "/autoscalegroups",
			AssetEndpoint: "/assets/autoscalegroups",
			ClassEndpoint: "/classes/autoscalegroups",
			PageType:      "Autoscale Group",
			ClassType:     "autoscalegroups",
		},
		"Load Balancers": components.Page{
			Route:         "/loadbalancers",
			AssetEndpoint: "/assets/loadbalancers",
			ClassEndpoint: "/classes/loadbalancers",
			PageType:      "Load Balancer",
			ClassType:     "loadbalancers",
		},
		"Scaling Policies": components.Page{
			Route:         "/scalingpolicies",
			AssetEndpoint: "/assets/scalingpolicies",
			ClassEndpoint: "/classes/scalingpolicies",
			PageType:      "Scaling Policy",
			ClassType:     "scalingpolicies",
		},
		"SimpleDB Domains": components.Page{
			Route:         "/simpledbdomains",
			AssetEndpoint: "/assets/simpledbdomains",
			ClassEndpoint: "/classes/simpledbdomains",
			PageType:      "SimpleDB Domain",
			ClassType:     "simpledbdomains",
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
