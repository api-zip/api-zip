// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022, The Zip API Object Framework Authors and Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.
package zip

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ObjectList represents a list of Objects.
type ObjectList[Spec, Status any] struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items contains the list of machines.
	Items []Object[Spec, Status] `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// Reference implements Referencable.
func (list *ObjectList[_, _]) Reference() (string, error) {
	return "", fmt.Errorf("cannot reference list")
}

// SetGroupVersionKind sets the API version and kind of the object reference
func (list *ObjectList[_, _]) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	list.APIVersion, list.Kind = gvk.ToAPIVersionAndKind()
}

// GroupVersionKind returns the API version and kind of the object reference
func (list *ObjectList[_, _]) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(list.APIVersion, list.Kind)
}

// GetObjectKind implements runtime.Object.
func (list *ObjectList[_, _]) GetObjectKind() schema.ObjectKind {
	return list
}

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in *ObjectList[Spec, Status]) DeepCopyInto(out *ObjectList[Spec, Status]) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Object[Spec, Status], len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy copies the receiver, creating a new ObjectList.
func (list *ObjectList[Spec, Status]) DeepCopy() *ObjectList[Spec, Status] {
	if list == nil {
		return nil
	}
	out := new(ObjectList[Spec, Status])
	list.DeepCopyInto(out)
	return out
}

// DeepCopyObject implements runtime.Object.
func (list *ObjectList[_, _]) DeepCopyObject() runtime.Object {
	if c := list.DeepCopy(); c != nil {
		return c
	}
	return nil
}
