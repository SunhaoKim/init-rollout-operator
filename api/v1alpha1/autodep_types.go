/*
Copyright 2021.

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
// AutodepSpec defines the desired state of Autodep
type AutodepSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Autodep. Edit autodep_types.go to remove/update
	Depimage           string `json:"depimage,omitempty"`
	Depenv             string `json:"depenv,omitempty"`
	Deptype            string `json:"deptype,omitempty"`
	DepimagePullSecret string `json:"depimagepullsecret,omitempty"`
	//AutoDeploy         bool   `json:"autodeploy"`
	SvcPort int32 `json:"svcport,omitempty"`
}

// AutodepStatus defines the observed state of Autodep
type AutodepStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Autodep is the Schema for the autodeps API
type Autodep struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutodepSpec   `json:"spec,omitempty"`
	Status AutodepStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AutodepList contains a list of Autodep
type AutodepList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Autodep `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Autodep{}, &AutodepList{})
}
