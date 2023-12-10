# Service

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The hostname of the service | 
**Port** | **int32** | The port to expose the service on | 
**TargetPort** | Pointer to **int32** | The port to forward traffic to, if different from the port. Can be left blank if the same as the port. | [optional] 
**Protocol** | **string** | The protocol to use for the service. Either http, tcp or udp. | 
**DomainName** | Pointer to **string** | The domain name to use for the service. If left blank, the service will only be available internally. If set, the service will be available externally at the given domain name. | [optional] 

## Methods

### NewService

`func NewService(name string, port int32, protocol string, ) *Service`

NewService instantiates a new Service object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServiceWithDefaults

`func NewServiceWithDefaults() *Service`

NewServiceWithDefaults instantiates a new Service object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Service) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Service) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Service) SetName(v string)`

SetName sets Name field to given value.


### GetPort

`func (o *Service) GetPort() int32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *Service) GetPortOk() (*int32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *Service) SetPort(v int32)`

SetPort sets Port field to given value.


### GetTargetPort

`func (o *Service) GetTargetPort() int32`

GetTargetPort returns the TargetPort field if non-nil, zero value otherwise.

### GetTargetPortOk

`func (o *Service) GetTargetPortOk() (*int32, bool)`

GetTargetPortOk returns a tuple with the TargetPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetPort

`func (o *Service) SetTargetPort(v int32)`

SetTargetPort sets TargetPort field to given value.

### HasTargetPort

`func (o *Service) HasTargetPort() bool`

HasTargetPort returns a boolean if a field has been set.

### GetProtocol

`func (o *Service) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *Service) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *Service) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.


### GetDomainName

`func (o *Service) GetDomainName() string`

GetDomainName returns the DomainName field if non-nil, zero value otherwise.

### GetDomainNameOk

`func (o *Service) GetDomainNameOk() (*string, bool)`

GetDomainNameOk returns a tuple with the DomainName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomainName

`func (o *Service) SetDomainName(v string)`

SetDomainName sets DomainName field to given value.

### HasDomainName

`func (o *Service) HasDomainName() bool`

HasDomainName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


