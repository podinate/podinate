package package_parser

import (
	hcldec "github.com/hashicorp/hcl2/hcldec"
	"github.com/johncave/podinate/cli/sdk"
	"github.com/zclconf/go-cty/cty"
)

type Pod struct {
	ID          string
	ProjectID   string `cty:"project_id"`
	Name        string `cty:"name"`
	Image       string `cty:"image"`
	Tag         string `cty:"tag"`
	Environment map[string]struct {
		Value  string `cty:"value"`
		Secret *bool  `cty:"secret"`
	} `cty:"environment"`
}

// GetHCLSpect returns the HCL spec of the pod type
var podHCLSpec = &hcldec.BlockMapSpec{
	TypeName:   "pod",
	LabelNames: []string{"id"},
	Nested: &hcldec.ObjectSpec{
		"project_id": &hcldec.AttrSpec{
			Name:     "project_id",
			Type:     cty.String,
			Required: true,
		},
		"name": &hcldec.AttrSpec{
			Name:     "name",
			Type:     cty.String,
			Required: true,
		},
		"image": &hcldec.AttrSpec{
			Name:     "image",
			Type:     cty.String,
			Required: true,
		},
		"tag": &hcldec.AttrSpec{
			Name:     "tag",
			Type:     cty.String,
			Required: true,
		},
		"environment": &hcldec.BlockMapSpec{
			TypeName:   "environment",
			LabelNames: []string{"key"},
			Nested: &hcldec.ObjectSpec{
				"value": &hcldec.AttrSpec{
					Name:     "value",
					Type:     cty.String,
					Required: true,
				},
				"secret": &hcldec.AttrSpec{
					Name:     "secret",
					Type:     cty.Bool,
					Required: false,
				},
			},
		},
	},
}

// Tosdk returns the API client representation of the pod
func (p *Pod) ToSDK() (*sdk.Pod, error) {
	theProject, err := sdk.ProjectGetByID(p.ProjectID)
	if err != nil {
		return nil, err
	}

	out := &sdk.Pod{
		Project: theProject,
		ID:      p.ID,
		Name:    p.Name,
		Image:   p.Image,
		Tag:     p.Tag,
	}

	for k, v := range p.Environment {
		new := sdk.EnvironmentVariable{
			Key:   k,
			Value: v.Value,
		}
		if v.Secret != nil {
			new.Secret = *v.Secret
		} else {
			new.Secret = false
		}
		out.Environment = append(out.Environment, new)
	}

	return out, nil
}
