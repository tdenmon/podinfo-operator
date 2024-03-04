/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Resources struct {
	MemoryLimit resource.Quantity `json:"memoryLimit"`
	CpuRequest  resource.Quantity `json:"cpuRequest"`
}

type Image struct {
	ImageRepository string `json:"repository"`
	ImageTag        string `json:"tag"`
}

type Ui struct {
	UiColor   string `json:"color"`
	UiMessage string `json:"message"`
}

type Redis struct {
	RedisEnabled bool `json:"enabled"`
}

// PodInfoSpec defines the desired state of PodInfo
type PodInfoSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ReplicaCount int32     `json:"replicaCount"`
	ResourceInfo Resources `json:"resources"`
	ImageInfo    Image     `json:"image"`
	UiInfo       Ui        `json:"ui"`
	RedisInfo    Redis     `json:"redis"`
}

// PodInfoStatus defines the observed state of PodInfo
type PodInfoStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ReplicaCount int32     `json:"replicaCount"`
	ResourceInfo Resources `json:"resources"`
	ImageInfo    Image     `json:"image"`
	UiInfo       Ui        `json:"ui"`
	RedisInfo    Redis     `json:"redis"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodInfo is the Schema for the podinfoes API
type PodInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodInfoSpec   `json:"spec,omitempty"`
	Status PodInfoStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodInfoList contains a list of PodInfo
type PodInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodInfo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodInfo{}, &PodInfoList{})
}
