package components

import (
	"github.com/Jeffail/gabs"
	"github.com/bep/gr"
	"github.com/bep/gr/attr"
	"github.com/bep/gr/el"
	"github.com/murdinc/awsm/config"
	"github.com/murdinc/awsmDashboard/components/forms"
)

func EditClassFormBuilder(classBytes []byte) (*gr.ReactComponent, *gabs.Container) {

	jsonParsed, _ := gabs.ParseJSON(classBytes)
	classType := jsonParsed.S("classType").Data().(string)

	switch classType {

	case "alarms":
		return gr.New(&forms.AlarmClassForm{}), jsonParsed

	case "instances":
		return gr.New(&forms.InstanceClassForm{}), jsonParsed

	case "volumes":
		return gr.New(&forms.VolumeClassForm{}), jsonParsed

	case "images":
		return gr.New(&forms.ImageClassForm{}), jsonParsed

	case "snapshots":
		return gr.New(&forms.SnapshotClassForm{}), jsonParsed

	case "vpcs":
		return gr.New(&forms.VpcClassForm{}), jsonParsed

	case "subnets":
		return gr.New(&forms.SubnetClassForm{}), jsonParsed

	case "securitygroups":
		return gr.New(&forms.SecurityGroupClassForm{}), jsonParsed

	case "launchconfigurations":
		return gr.New(&forms.LaunchConfigurationClassForm{}), jsonParsed

	case "loadbalancers":
		return gr.New(&forms.LoadBalancerClassForm{}), jsonParsed

	case "scalingpolicies":
		return gr.New(&forms.ScalingPolicyClassForm{}), jsonParsed

	case "autoscalegroups":
		return gr.New(&forms.AutoscaleGroupClassForm{}), jsonParsed

	case "keypairs":
		return gr.New(&forms.KeyPairClassForm{}), jsonParsed

	//case "simpledbdomains":
	//return gr.New(&forms.SimpleDBDomainClassForm{}), jsonParsed

	//case "addresses":
	//return gr.New(&forms.AddressClassForm{}), jsonParsed

	default:
		println("Class Type not found in EditClassFormBuilder switch:")
		println(classType)
	}

	return nil, nil
}

func NewClassFormBuilder(classType string) *gr.ReactComponent {

	switch classType {

	case "alarms":
		return gr.New(&forms.AlarmClassForm{})

	case "instances":
		return gr.New(&forms.InstanceClassForm{})

	case "volumes":
		return gr.New(&forms.VolumeClassForm{})

	case "images":
		return gr.New(&forms.ImageClassForm{})

	case "keypairs":
		return gr.New(&forms.KeyPairClassForm{})

	case "snapshots":
		return gr.New(&forms.SnapshotClassForm{})

	case "vpcs":
		return gr.New(&forms.VpcClassForm{})

	case "subnets":
		return gr.New(&forms.SubnetClassForm{})

	case "securitygroups":
		return gr.New(&forms.SecurityGroupClassForm{})

	case "launchconfigurations":
		return gr.New(&forms.LaunchConfigurationClassForm{})

	case "loadbalancers":
		return gr.New(&forms.LoadBalancerClassForm{})

	case "scalingpolicies":
		return gr.New(&forms.ScalingPolicyClassForm{})

	case "autoscalegroups":
		return gr.New(&forms.AutoscaleGroupClassForm{})

	//case "simpledbdomains":
	// TODO ?

	//case "addresses":
	//return gr.New(&forms.InstanceClassForm{})

	default:
		println("Class Type not found in NewClassFormBuilder switch:")
		println(classType)
	}

	return nil
}

func BuildClassForm(className string, cType interface{}) *gr.Element {

	keys, values := config.ExtractAwsmClass(cType)

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
