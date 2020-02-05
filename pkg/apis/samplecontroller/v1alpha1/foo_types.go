package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FooSpec defines the desired state of Foo
type FooSpec struct {
	// the name of deployment which is owned by foo
	// +kubebuilder:validation:Format=string
	DeploymentName string `json:"deploymentName"`

	// the replicas of deployment which is owned by foo
	// +kubebuilder:validation:Minimum=1
	Replicas *int32 `json:"replicas"`
}

// FooStatus defines the observed state of Foo
type FooStatus struct {
	// this is equal deployment.status.availableReplicas
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is the Schema for the foos API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=foos,scope=Namespaced
type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooSpec   `json:"spec,omitempty"`
	Status FooStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList contains a list of Foo
type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Foo `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Foo{}, &FooList{})
}
