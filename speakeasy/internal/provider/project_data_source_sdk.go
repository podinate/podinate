// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/podinate/terraform-provider-podinate/internal/sdk/pkg/models/shared"
)

func (r *ProjectDataSourceModel) RefreshFromGetResponse(resp *shared.Project) {
	if resp.ID != nil {
		r.ID = types.StringValue(*resp.ID)
	} else {
		r.ID = types.StringNull()
	}
	if resp.Name != nil {
		r.Name = types.StringValue(*resp.Name)
	} else {
		r.Name = types.StringNull()
	}
	if resp.ResourceID != nil {
		r.ResourceID = types.StringValue(*resp.ResourceID)
	} else {
		r.ResourceID = types.StringNull()
	}
}