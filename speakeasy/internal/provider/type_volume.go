// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type Volume struct {
	MountPath types.String `tfsdk:"mount_path"`
	Name      types.String `tfsdk:"name"`
	Size      types.Int64  `tfsdk:"size"`
}
