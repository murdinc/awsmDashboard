package main

import (
	"github.com/bep/gr"
	"github.com/bep/gr/el"
	"github.com/bep/grouter"
	"github.com/gopherjs/gopherjs/js"
	"github.com/murdinc/awsmDashboard/components"
)

var (
	reactRouter   = js.Global.Get("ReactRouter")
	routerHistory = grouter.History{Object: reactRouter.Get("browserHistory")}

	dashboard            = gr.New(&components.Layout{Title: "Dashboard", Uri: "/assets/images"}, gr.Apply(grouter.WithRouter))
	instances            = gr.New(&components.Layout{Title: "Instances", Uri: "/assets/instances"}, gr.Apply(grouter.WithRouter))
	images               = gr.New(&components.Layout{Title: "Images", Uri: "/assets/images"}, gr.Apply(grouter.WithRouter))
	volumes              = gr.New(&components.Layout{Title: "Volumes", Uri: "/assets/volumes"}, gr.Apply(grouter.WithRouter))
	snapshots            = gr.New(&components.Layout{Title: "Snapshots", Uri: "/assets/snapshots"}, gr.Apply(grouter.WithRouter))
	vpcs                 = gr.New(&components.Layout{Title: "VPCs", Uri: "/assets/vpcs"}, gr.Apply(grouter.WithRouter))
	subnets              = gr.New(&components.Layout{Title: "Subnets", Uri: "/assets/subnets"}, gr.Apply(grouter.WithRouter))
	securitygroups       = gr.New(&components.Layout{Title: "Security Groups", Uri: "/assets/securitygroups"}, gr.Apply(grouter.WithRouter))
	addresses            = gr.New(&components.Layout{Title: "Addresses", Uri: "/assets/addresses"}, gr.Apply(grouter.WithRouter))
	alarms               = gr.New(&components.Layout{Title: "Alarms", Uri: "/assets/alarms"}, gr.Apply(grouter.WithRouter))
	keypairs             = gr.New(&components.Layout{Title: "Keypairs", Uri: "/assets/keypairs"}, gr.Apply(grouter.WithRouter))
	launchconfigurations = gr.New(&components.Layout{Title: "Launch Configurations", Uri: "/assets/launchconfigurations"}, gr.Apply(grouter.WithRouter))
	loadbalancers        = gr.New(&components.Layout{Title: "Load Balancers", Uri: "/assets/loadbalancers"}, gr.Apply(grouter.WithRouter))
	scalingpolicies      = gr.New(&components.Layout{Title: "Scaling Policies", Uri: "/assets/scalingpolicies"}, gr.Apply(grouter.WithRouter))
	simpledbdomains      = gr.New(&components.Layout{Title: "SimpleDB Domains", Uri: "/assets/simpledbdomains"}, gr.Apply(grouter.WithRouter))

	// WithRouter makes this.props.router happen.
	appComponent = gr.New(new(app), gr.Apply(grouter.WithRouter))

	router = grouter.New("/", appComponent, grouter.WithHistory(routerHistory)).With(
		grouter.NewIndexRoute(grouter.Components{"layout": dashboard}),
		grouter.NewRoute("/instances", grouter.Components{"layout": instances}),
		grouter.NewRoute("/volumes", grouter.Components{"layout": volumes}),
		grouter.NewRoute("/images", grouter.Components{"layout": images}),
		grouter.NewRoute("/snapshots", grouter.Components{"layout": snapshots}),
		grouter.NewRoute("/vpcs", grouter.Components{"layout": vpcs}),
		grouter.NewRoute("/subnets", grouter.Components{"layout": subnets}),
		grouter.NewRoute("/securitygroups", grouter.Components{"layout": securitygroups}),
		grouter.NewRoute("/addresses", grouter.Components{"layout": addresses}),
		grouter.NewRoute("/alarms", grouter.Components{"layout": alarms}),
		grouter.NewRoute("/keypairs", grouter.Components{"layout": keypairs}),
		grouter.NewRoute("/launchconfigurations", grouter.Components{"layout": launchconfigurations}),
		grouter.NewRoute("/loadbalancers", grouter.Components{"layout": loadbalancers}),
		grouter.NewRoute("/scalingpolicies", grouter.Components{"layout": scalingpolicies}),
		grouter.NewRoute("/simpledbdomains", grouter.Components{"layout": simpledbdomains}),
	)
)

func main() {
	mainComponent := gr.New(gr.NewSimpleRenderer(router))
	mainComponent.Render("react", gr.Props{})
}

type app struct {
	*gr.This
	//Nav components.Nav
}

// Implements the Renderer interface.
func (a app) Render() gr.Component {
	return el.Div(
		a.Component("layout"),
	)
}
