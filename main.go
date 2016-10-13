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
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Instance",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			ClassEndpoint: "/classes/instances",
		},
		"Instances": components.Page{
			Route: "/instances",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Instance",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/instances",
			ClassEndpoint: "/classes/instances",
		},
		"Volumes": components.Page{
			Route: "/volumes",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Volume",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/volumes",
			ClassEndpoint: "/classes/volumes",
		},
		"Images": components.Page{
			Route: "/images",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Image",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/images",
			ClassEndpoint: "/classes/images",
		},
		"Snapshots": components.Page{
			Route: "/snapshots",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Snapshot",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/snapshots",
			ClassEndpoint: "/classes/snapshots",
		},
		"Vpcs": components.Page{
			Route: "/vpcs",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Vpc",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/vpcs",
			ClassEndpoint: "/classes/vpcs",
		},
		"Subnets": components.Page{
			Route: "/subnets",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Subnet",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/subnets",
			ClassEndpoint: "/classes/subnets",
		},
		"Security Groups": components.Page{
			Route: "/securitygroups",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Security Group",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/securitygroups",
			ClassEndpoint: "/classes/securitygroups",
		},
		"Addresses": components.Page{
			Route: "/addresses",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Address",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/addresses",
			ClassEndpoint: "/classes/addresses",
		},
		"Alarms": components.Page{
			Route: "/alarms",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Alarm",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/alarms",
			ClassEndpoint: "/classes/alarms",
		},
		"Keypairs": components.Page{
			Route: "/keypairs",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Key Pair",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/keypairs",
			ClassEndpoint: "/classes/keypairs",
		},
		"Launch Configurations": components.Page{
			Route: "/launchconfigurations",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Launch Configuration",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/launchconfigurations",
			ClassEndpoint: "/classes/launchconfigurations",
		},
		"Load Balancers": components.Page{
			Route: "/loadbalancers",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Load Balancer",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/loadbalancers",
			ClassEndpoint: "/classes/loadbalancers",
		},
		"Scaling Policies": components.Page{
			Route: "/scalingpolicies",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New Scaling Policy",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/scalingpolicies",
			ClassEndpoint: "/classes/scalingpolicies",
		},
		"SimpleDB Domains": components.Page{
			Route: "/simpledbdomains",
			DropdownOptions: []components.DropdownOption{
				components.DropdownOption{
					Name: "New SimpleDB Domain",
					Id:   "New",
				},
				components.DropdownOption{
					Name: "Edit Classes",
					Id:   "Edit",
				},
			},
			AssetEndpoint: "/assets/simpledbdomains",
			ClassEndpoint: "/classes/simpledbdomains",
		},
	}

	reactRouter   = js.Global.Get("ReactRouter")
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
