package components

import (
	"encoding/json"

	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/bep/gr/evt"
	"github.com/murdinc/awsm/config"
)

func ClassListBuilder(cl interface{}, onClick func(string)) *gr.Element {

	classList := cl.([]byte)

	jsonParsed, _ := gabs.ParseJSON(classList)
	classType := jsonParsed.S("classType").Data().(string)
	classes, _ := jsonParsed.S("classes").ChildrenMap()

	classListGroup := el.Div(
		el.Div(
			gr.CSS("list-group"),
		),
	)

	switch classType {

	case "alarms":
		var cType config.AlarmClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	case "instances":
		var cType config.InstanceClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	case "volumes":
		var cType config.VolumeClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	case "images":
		var cType config.ImageClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	//case "keypairs":
	// TODO ?

	case "snapshots":
		var cType config.SnapshotClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	case "vpcs":
		var cType config.VpcClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	case "subnets":
		var cType config.SubnetClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	case "securitygroups":
		var cType config.SecurityGroupClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	//case "addresses":
	// TODO ?

	case "launchconfigurations":
		var cType config.LaunchConfigurationClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	case "autoscalegroups":
		var cType config.AutoscaleGroupClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	case "loadbalancers":
		var cType config.LoadBalancerClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	case "scalingpolicies":
		var cType config.ScalingPolicyClass
		for className, class := range classes {
			cJson := class.Bytes()
			json.Unmarshal(cJson, &cType)
			keys, values := config.ExtractAwsmClass(cType)
			buildClassButton(className, keys, values, classListGroup, onClick)
		}

	//case "simpledbdomains":
	// TODO ?

	default:
		println("Class Type not found in ClassListBuilder switch:")
		println(classType)
	}

	return classListGroup
}

func buildClassButton(name string, keys []string, values []string, classListGroup *gr.Element, onClick func(string)) {

	clickListener := func(event *gr.Event) {
		onClick(name)
	}

	button := el.Button(
		attr.Type("button"),
		gr.CSS("list-group-item"),
		attr.ID(name),
		evt.Click(clickListener),
		el.Header5(
			gr.CSS("list-group-item-heading"),
			gr.Text(name),
		),
	)

	description := el.Div(
		gr.CSS("list-group-item-text"),
	)

	left := el.DescriptionList(gr.CSS("dl-horizontal"))
	right := el.DescriptionList(gr.CSS("dl-horizontal"))

	keyLen := len(keys)

	for i, key := range keys {
		if i < (keyLen / 2) {
			el.DefinitionTerm(gr.Text(key)).Modify(left)
			el.Description(gr.Text(values[i])).Modify(left)
		} else {
			el.DefinitionTerm(gr.Text(key)).Modify(right)
			el.Description(gr.Text(values[i])).Modify(right)
		}
	}

	el.Div(
		gr.CSS("row"),
		el.Div(gr.CSS("col-sm-6"),
			left),
		el.Div(gr.CSS("col-sm-6"),
			right),
	).Modify(description)

	description.Modify(button)
	button.Modify(classListGroup)

}
