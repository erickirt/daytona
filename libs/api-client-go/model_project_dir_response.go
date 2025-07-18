/*
Daytona

Daytona AI platform API Docs

API version: 1.0
Contact: support@daytona.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"encoding/json"
)

// checks if the ProjectDirResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ProjectDirResponse{}

// ProjectDirResponse struct for ProjectDirResponse
type ProjectDirResponse struct {
	Dir *string `json:"dir,omitempty"`
}

// NewProjectDirResponse instantiates a new ProjectDirResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProjectDirResponse() *ProjectDirResponse {
	this := ProjectDirResponse{}
	return &this
}

// NewProjectDirResponseWithDefaults instantiates a new ProjectDirResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProjectDirResponseWithDefaults() *ProjectDirResponse {
	this := ProjectDirResponse{}
	return &this
}

// GetDir returns the Dir field value if set, zero value otherwise.
func (o *ProjectDirResponse) GetDir() string {
	if o == nil || IsNil(o.Dir) {
		var ret string
		return ret
	}
	return *o.Dir
}

// GetDirOk returns a tuple with the Dir field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProjectDirResponse) GetDirOk() (*string, bool) {
	if o == nil || IsNil(o.Dir) {
		return nil, false
	}
	return o.Dir, true
}

// HasDir returns a boolean if a field has been set.
func (o *ProjectDirResponse) HasDir() bool {
	if o != nil && !IsNil(o.Dir) {
		return true
	}

	return false
}

// SetDir gets a reference to the given string and assigns it to the Dir field.
func (o *ProjectDirResponse) SetDir(v string) {
	o.Dir = &v
}

func (o ProjectDirResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ProjectDirResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Dir) {
		toSerialize["dir"] = o.Dir
	}
	return toSerialize, nil
}

type NullableProjectDirResponse struct {
	value *ProjectDirResponse
	isSet bool
}

func (v NullableProjectDirResponse) Get() *ProjectDirResponse {
	return v.value
}

func (v *NullableProjectDirResponse) Set(val *ProjectDirResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableProjectDirResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableProjectDirResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProjectDirResponse(val *ProjectDirResponse) *NullableProjectDirResponse {
	return &NullableProjectDirResponse{value: val, isSet: true}
}

func (v NullableProjectDirResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProjectDirResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
