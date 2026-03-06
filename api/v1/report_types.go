/*
Copyright 2026.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ReportState string

const (
	ReportStatePending  ReportState = "Pending"
	ReportStateFinished ReportState = "Finished"
)

// ReportSpec defines the desired state of Report.
type ReportSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=10
	Priority int `json:"priority,omitempty"`
	// +kubebuilder:validation:Optional
	Details string `json:"details,omitempty"`
	// +kubebuilder:validation:Required
	Deadline metav1.Time `json:"deadline,omitempty"`
	// +kubebuilder:default="Pending"
	// +kubebuilder:validation:Enum=Pending;Finished
	ReportState ReportState `json:"repoState,omitempty"`
}

// ReportStatus defines the observed state of Report.
type ReportStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Priority",type="string",JSONPath=".spec.priority"
// +kubebuilder:printcolumn:name="DeadLine",type="string",JSONPath=".spec.deadline"
// +kubebuilder:printcolumn:name="ReportState",type="string",JSONPath=".spec.repoState"
// +kubebuilder:selectablefield:JSONPath=.spec.priority
// +kubebuilder:selectablefield:JSONPath=.spec.repoState

// Report is the Schema for the reports API.
type Report struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReportSpec   `json:"spec,omitempty"`
	Status ReportStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ReportList contains a list of Report.
type ReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Report `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Report{}, &ReportList{})
}
