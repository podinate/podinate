package package_parser

import (
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/johncave/podinate/cli/sdk"
	"github.com/zclconf/go-cty/cty"
)

type Project struct {
	ID        string
	Name      string `cty:"name"`
	AccountID string `cty:"account_id"`
}

var projectHCLSpec = &hcldec.BlockMapSpec{
	TypeName:   "project",
	LabelNames: []string{"id"},
	Nested: &hcldec.ObjectSpec{
		"name": &hcldec.AttrSpec{
			Name:     "name",
			Type:     cty.String,
			Required: true,
		},
		"account_id": &hcldec.AttrSpec{
			Name:     "account_id",
			Type:     cty.String,
			Required: false,
		},
	},
}

// ToSDK converts a project from the package parser to the SDK type
func (p *Project) ToSDK() *sdk.Project {
	return &sdk.Project{
		ID:   p.ID,
		Name: p.Name,
	}
}
