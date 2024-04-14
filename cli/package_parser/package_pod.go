package package_parser

import (
	hcldec "github.com/hashicorp/hcl2/hcldec"
	"github.com/johncave/podinate/cli/sdk"
	"github.com/zclconf/go-cty/cty"
)

type Pod struct {
	ID          string
	ProjectID   string   `cty:"project_id"`
	Name        string   `cty:"name"`
	Image       string   `cty:"image"`
	Tag         *string  `cty:"tag"`
	Command     []string `cty:"command"`
	Arguments   []string `cty:"arguments"`
	Environment map[string]struct {
		Value  string `cty:"value"`
		Secret *bool  `cty:"secret"`
	} `cty:"environment"`
	Service map[string]struct {
		Port       int     `cty:"port"`
		TargetPort *int    `cty:"target_port"`
		Protocol   *string `cty:"protocol"`
		DomainName *string `cty:"domain_name"`
	} `cty:"service"`
	Volume map[string]struct {
		Size  int     `cty:"size"`
		Path  string  `cty:"path"`
		Class *string `cty:"class"`
	} `cty:"volume"`
	SharedVolume *[]struct {
		VolumeID string `cty:"volume_id"`
		Path     string `cty:"path"`
	} `cty:"shared_volume"`
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
			Required: false,
		},
		"command": &hcldec.AttrSpec{
			Name:     "command",
			Type:     cty.List(cty.String),
			Required: false,
		},
		"arguments": &hcldec.AttrSpec{
			Name:     "arguments",
			Type:     cty.List(cty.String),
			Required: false,
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
		"service": &hcldec.BlockMapSpec{
			TypeName:   "service",
			LabelNames: []string{"name"},
			Nested: &hcldec.ObjectSpec{
				"port": &hcldec.AttrSpec{
					Name:     "port",
					Type:     cty.Number,
					Required: true,
				},
				"target_port": &hcldec.AttrSpec{
					Name:     "target_port",
					Type:     cty.Number,
					Required: false,
				},
				"protocol": &hcldec.AttrSpec{
					Name:     "protocol",
					Type:     cty.String,
					Required: false,
				},
				"domain_name": &hcldec.AttrSpec{
					Name:     "domain_name",
					Type:     cty.String,
					Required: false,
				},
			},
		},
		"volume": &hcldec.BlockMapSpec{
			TypeName:   "volume",
			LabelNames: []string{"name"},
			Nested: &hcldec.ObjectSpec{
				"size": &hcldec.AttrSpec{
					Name:     "size",
					Type:     cty.Number,
					Required: true,
				},
				"path": &hcldec.AttrSpec{
					Name:     "path",
					Type:     cty.String,
					Required: true,
				},
				"class": &hcldec.AttrSpec{
					Name:     "class",
					Type:     cty.String,
					Required: false,
				},
			},
		},
		"shared_volume": &hcldec.BlockListSpec{
			TypeName: "shared_volume",
			MinItems: 0,
			Nested: &hcldec.ObjectSpec{
				"volume_id": &hcldec.AttrSpec{
					Name:     "volume_id",
					Type:     cty.String,
					Required: true,
				},
				"path": &hcldec.AttrSpec{
					Name:     "path",
					Type:     cty.String,
					Required: true,
				},
			},
		},
	},
}

// Tosdk returns the API client representation of the pod
func (p *Pod) ToSDK() (*sdk.Pod, error) {
	theProject, err := sdk.GetProjectByID(p.ProjectID)
	if err != nil {
		return nil, err
	}

	// Get all services
	var services sdk.ServiceSlice

	for k, v := range p.Service {
		new := sdk.Service{
			Name: k,
			Port: v.Port,
		}
		if v.TargetPort != nil {
			new.TargetPort = v.TargetPort
		}
		if v.Protocol != nil {
			new.Protocol = *v.Protocol
		}
		if v.DomainName != nil {
			new.DomainName = v.DomainName
		}
		services = append(services, new)
	}

	// Get all volumes
	var volumes sdk.VolumeSlice
	for k, v := range p.Volume {
		new := sdk.Volume{
			Name: k,
			Size: v.Size,
			Path: v.Path,
		}
		if v.Class != nil {
			new.Class = *v.Class
		}
		volumes = append(volumes, new)
	}

	var sharedVolumes sdk.SharedVolumeAttachmentSlice
	if p.SharedVolume != nil {
		for _, v := range *p.SharedVolume {
			new := sdk.SharedVolumeAttachment{
				ID:   v.VolumeID,
				Path: v.Path,
			}
			sharedVolumes = append(sharedVolumes, new)
		}
	}

	out := &sdk.Pod{
		Project:       theProject,
		ID:            p.ID,
		Name:          p.Name,
		Image:         p.Image,
		Command:       p.Command,
		Arguments:     p.Arguments,
		Services:      services,
		Volumes:       volumes,
		SharedVolumes: sharedVolumes,
	}

	if p.Tag != nil {
		out.Tag = p.Tag
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
