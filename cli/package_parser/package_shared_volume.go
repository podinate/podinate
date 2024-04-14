package package_parser

import (
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/johncave/podinate/cli/sdk"
	"github.com/zclconf/go-cty/cty"
)

type SharedVolume struct {
	ID        string
	Name      *string `cty:"name"`
	ProjectID string  `cty:"project_id"`
	Size      int     `cty:"size"`
	Class     *string `cty:"class"`
}

var sharedVolumeHCLSpec = &hcldec.BlockMapSpec{
	TypeName:   "shared_volume",
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
			Required: false,
		},
		"size": &hcldec.AttrSpec{
			Name:     "size",
			Type:     cty.Number,
			Required: true,
		},
		"class": &hcldec.AttrSpec{
			Name:     "class",
			Type:     cty.String,
			Required: false,
		},
	},
}

// ToSDK converts a shared volume from the package parser to the SDK type
func (sv *SharedVolume) ToSDK() (*sdk.SharedVolume, error) {
	theProject, err := sdk.GetProjectByID(sv.ProjectID)
	if err != nil {
		return nil, err
	}
	return &sdk.SharedVolume{
		ID:      sv.ID,
		Name:    sv.Name,
		Project: theProject,
		Size:    sv.Size,
		Class:   sv.Class,
	}, nil
}
