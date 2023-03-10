/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ValidationSpec defines the desired state of Validation
type ValidationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Size is the size of the ansibleee deployment
	Size int32 `json:"size"`
	// Command string `json:"command"`

	//RestartPolicy is the policy applied to the Job on whether it needs to restart the Pod
	//+kubebuilder:default:="OnFailure"
	RestartPolicy string `json:"restartPolicy,omitempty"`
	// Container image for validation
	ContainerImage string `json:"containerImage"`
	// Playbook is the playbook that ansible will run on this execution
	Validation string `json:"validation,omitempty"`
	// Image is the container image that will execute the ansible command
	// +kubebuilder:default:="docker.io/matbu/validation"
	Image string `json:"image,omitempty"`
	// Command is the command executed by the image
	Command []string `json:"command,omitempty"`
	// Name is the name of the internal container inside the pod
	// +kubebuilder:default:="validationframework"
	Name string `json:"name,omitempty"`
	// Args for the command of the validation
	Args []string `json:"args,omitempty"`
	// Uid is the userid that will be used to run the container
	// +kubebuilder:default:=1001
	Uid int64 `json:"uid,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Validation is the Schema for the validations API
type Validation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ValidationSpec   `json:"spec,omitempty"`
	Status ValidationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ValidationList contains a list of Validation
type ValidationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Validation `json:"items"`
}

// ValidationStatus defines the observed state of Validation
type ValidationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Nodes are the names of the validation pods
	Nodes []string `json:"nodes"`
}

func init() {
	SchemeBuilder.Register(&Validation{}, &ValidationList{})
}
