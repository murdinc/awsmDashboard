package components

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsm/config"
)

func EditClassFormBuilder(classBytes interface{}) *gr.Element {

	jsonParsed, _ := gabs.ParseJSON(classBytes.([]byte))
	className := jsonParsed.S("className").Data().(string)
	classType := jsonParsed.S("classType").Data().(string)
	class := jsonParsed.S("class").Bytes()

	switch classType {

	case "alarms":
		var cType config.AlarmClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "instances":
		var cType config.InstanceClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "volumes":
		var cType config.VolumeClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "images":
		var cType config.ImageClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	//case "keypairs":
	// TODO ?

	case "snapshots":
		var cType config.SnapshotClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "vpcs":
		var cType config.VpcClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "subnets":
		var cType config.SubnetClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "securitygroups":
		var cType config.SecurityGroupClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	//case "addresses":
	// TODO ?

	case "launchconfigurations":
		var cType config.LaunchConfigurationClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "loadbalancers":
		var cType config.LoadBalancerClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "scalingpolicies":
		var cType config.ScalingPolicyClass
		json.Unmarshal(class, &cType)
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	//case "simpledbdomains":
	// TODO ?

	default:
		println("Class Type not found in ClassFormBuilder switch:")
		println(classType)
	}

	return nil
}

func NewClassFormBuilder(className, classType string) *gr.Element {

	switch classType {

	case "alarms":
		var cType config.AlarmClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "instances":
		var cType config.InstanceClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "volumes":
		var cType config.VolumeClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "images":
		var cType config.ImageClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	//case "keypairs":
	// TODO ?

	case "snapshots":
		var cType config.SnapshotClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "vpcs":
		var cType config.VpcClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "subnets":
		var cType config.SubnetClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "securitygroups":
		var cType config.SecurityGroupClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	//case "addresses":
	// TODO ?

	case "launchconfigurations":
		var cType config.LaunchConfigurationClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "loadbalancers":
		var cType config.LoadBalancerClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	case "scalingpolicies":
		var cType config.ScalingPolicyClass
		keys, values := config.ExtractAwsmList(cType)
		return BuildClassForm(className, keys, values)

	//case "simpledbdomains":
	// TODO ?

	default:
		println("Class Type not found in ClassFormBuilder switch:")
		println(classType)
	}

	return nil
}

func BuildClassForm(className string, keys []string, values []string) *gr.Element {

	classEdit := el.Div(
		el.Header3(gr.Text(className)),
		el.HorizontalRule(),
	)

	classEditForm := el.Form()

	for i, key := range keys {

		el.Div(
			gr.CSS("form-group"),
			el.Label(
				gr.Text(key),
			),
			el.Input(
				attr.Type(gr.Text(key)),
				attr.ClassName("form-control"),
				attr.ID(gr.Text(key)),
				attr.Placeholder(key),
				attr.DefaultValue(values[i]),
			),
		).Modify(classEditForm)
	}

	classEditForm.Modify(classEdit)
	return classEdit

}
