// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/podinate/terraform-provider-podinate/internal/sdk/pkg/models/shared"
)

func (r *PodResourceModel) ToCreateSDKType() *shared.Pod {
	id := new(string)
	if !r.ID.IsUnknown() && !r.ID.IsNull() {
		*id = r.ID.ValueString()
	} else {
		id = nil
	}
	name := r.Name.ValueString()
	image := r.Image.ValueString()
	tag := r.Tag.ValueString()
	var volumes []shared.Volume = nil
	for _, volumesItem := range r.Volumes {
		name1 := volumesItem.Name.ValueString()
		size := volumesItem.Size.ValueInt64()
		mountPath := volumesItem.MountPath.ValueString()
		volumes = append(volumes, shared.Volume{
			Name:      name1,
			Size:      size,
			MountPath: mountPath,
		})
	}
	var environment []shared.EnvironmentVariable = nil
	for _, environmentItem := range r.Environment {
		key := environmentItem.Key.ValueString()
		value := environmentItem.Value.ValueString()
		secret := new(bool)
		if !environmentItem.Secret.IsUnknown() && !environmentItem.Secret.IsNull() {
			*secret = environmentItem.Secret.ValueBool()
		} else {
			secret = nil
		}
		environment = append(environment, shared.EnvironmentVariable{
			Key:    key,
			Value:  value,
			Secret: secret,
		})
	}
	var services []shared.Service = nil
	for _, servicesItem := range r.Services {
		name2 := servicesItem.Name.ValueString()
		port := servicesItem.Port.ValueInt64()
		targetPort := new(int64)
		if !servicesItem.TargetPort.IsUnknown() && !servicesItem.TargetPort.IsNull() {
			*targetPort = servicesItem.TargetPort.ValueInt64()
		} else {
			targetPort = nil
		}
		protocol := servicesItem.Protocol.ValueString()
		domainName := new(string)
		if !servicesItem.DomainName.IsUnknown() && !servicesItem.DomainName.IsNull() {
			*domainName = servicesItem.DomainName.ValueString()
		} else {
			domainName = nil
		}
		services = append(services, shared.Service{
			Name:       name2,
			Port:       port,
			TargetPort: targetPort,
			Protocol:   protocol,
			DomainName: domainName,
		})
	}
	status := new(string)
	if !r.Status.IsUnknown() && !r.Status.IsNull() {
		*status = r.Status.ValueString()
	} else {
		status = nil
	}
	createdAt := new(string)
	if !r.CreatedAt.IsUnknown() && !r.CreatedAt.IsNull() {
		*createdAt = r.CreatedAt.ValueString()
	} else {
		createdAt = nil
	}
	resourceID := new(string)
	if !r.ResourceID.IsUnknown() && !r.ResourceID.IsNull() {
		*resourceID = r.ResourceID.ValueString()
	} else {
		resourceID = nil
	}
	out := shared.Pod{
		ID:          id,
		Name:        name,
		Image:       image,
		Tag:         tag,
		Volumes:     volumes,
		Environment: environment,
		Services:    services,
		Status:      status,
		CreatedAt:   createdAt,
		ResourceID:  resourceID,
	}
	return &out
}

func (r *PodResourceModel) ToGetSDKType() *shared.Pod {
	out := r.ToCreateSDKType()
	return out
}

func (r *PodResourceModel) ToUpdateSDKType() *shared.Pod {
	out := r.ToCreateSDKType()
	return out
}

func (r *PodResourceModel) ToDeleteSDKType() *shared.Pod {
	out := r.ToCreateSDKType()
	return out
}

func (r *PodResourceModel) RefreshFromGetResponse(resp *shared.Pod) {
	if resp.CreatedAt != nil {
		r.CreatedAt = types.StringValue(*resp.CreatedAt)
	} else {
		r.CreatedAt = types.StringNull()
	}
	if len(r.Environment) > len(resp.Environment) {
		r.Environment = r.Environment[:len(resp.Environment)]
	}
	for environmentCount, environmentItem := range resp.Environment {
		var environment1 EnvironmentVariable
		environment1.Key = types.StringValue(environmentItem.Key)
		if environmentItem.Secret != nil {
			environment1.Secret = types.BoolValue(*environmentItem.Secret)
		} else {
			environment1.Secret = types.BoolNull()
		}
		environment1.Value = types.StringValue(environmentItem.Value)
		if environmentCount+1 > len(r.Environment) {
			r.Environment = append(r.Environment, environment1)
		} else {
			r.Environment[environmentCount].Key = environment1.Key
			r.Environment[environmentCount].Secret = environment1.Secret
			r.Environment[environmentCount].Value = environment1.Value
		}
	}
	if resp.ID != nil {
		r.ID = types.StringValue(*resp.ID)
	} else {
		r.ID = types.StringNull()
	}
	r.Image = types.StringValue(resp.Image)
	r.Name = types.StringValue(resp.Name)
	if resp.ResourceID != nil {
		r.ResourceID = types.StringValue(*resp.ResourceID)
	} else {
		r.ResourceID = types.StringNull()
	}
	if len(r.Services) > len(resp.Services) {
		r.Services = r.Services[:len(resp.Services)]
	}
	for servicesCount, servicesItem := range resp.Services {
		var services1 Service
		if servicesItem.DomainName != nil {
			services1.DomainName = types.StringValue(*servicesItem.DomainName)
		} else {
			services1.DomainName = types.StringNull()
		}
		services1.Name = types.StringValue(servicesItem.Name)
		services1.Port = types.Int64Value(servicesItem.Port)
		services1.Protocol = types.StringValue(servicesItem.Protocol)
		if servicesItem.TargetPort != nil {
			services1.TargetPort = types.Int64Value(*servicesItem.TargetPort)
		} else {
			services1.TargetPort = types.Int64Null()
		}
		if servicesCount+1 > len(r.Services) {
			r.Services = append(r.Services, services1)
		} else {
			r.Services[servicesCount].DomainName = services1.DomainName
			r.Services[servicesCount].Name = services1.Name
			r.Services[servicesCount].Port = services1.Port
			r.Services[servicesCount].Protocol = services1.Protocol
			r.Services[servicesCount].TargetPort = services1.TargetPort
		}
	}
	if resp.Status != nil {
		r.Status = types.StringValue(*resp.Status)
	} else {
		r.Status = types.StringNull()
	}
	r.Tag = types.StringValue(resp.Tag)
	if len(r.Volumes) > len(resp.Volumes) {
		r.Volumes = r.Volumes[:len(resp.Volumes)]
	}
	for volumesCount, volumesItem := range resp.Volumes {
		var volumes1 Volume
		volumes1.MountPath = types.StringValue(volumesItem.MountPath)
		volumes1.Name = types.StringValue(volumesItem.Name)
		volumes1.Size = types.Int64Value(volumesItem.Size)
		if volumesCount+1 > len(r.Volumes) {
			r.Volumes = append(r.Volumes, volumes1)
		} else {
			r.Volumes[volumesCount].MountPath = volumes1.MountPath
			r.Volumes[volumesCount].Name = volumes1.Name
			r.Volumes[volumesCount].Size = volumes1.Size
		}
	}
}

func (r *PodResourceModel) RefreshFromCreateResponse(resp *shared.Pod) {
	r.RefreshFromGetResponse(resp)
}

func (r *PodResourceModel) RefreshFromUpdateResponse(resp *shared.Pod) {
	r.RefreshFromGetResponse(resp)
}
