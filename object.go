// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022, The Zip API Object Framework Authors and Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.
package zip

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Referencable is an interface for allowing the object to return a unique
// reference for lookups.
type Referencable interface {
	Reference() (string, error)
}

// ReferenceObject combines the low-level Kubernetes-centric runtime.Object
// and the Zip API Referencable to allow mutable, copyable objects to be
// referenced at runtime.
type ReferenceObject interface {
	Referencable
	runtime.Object
}

// Object represents a generic low-level Kubernetes-inspired object which
// contains information about the object type, its metadata, a specification
// which is used to describe the object and its current status.
type Object[Spec, Status any] struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the Object.
	Spec Spec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the Object.
	Status Status `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// Reference implements Referencable
func (obj *Object[_, _]) Reference() (string, error) {
	return obj.ObjectMeta.Name, nil
}

// SetGroupVersionKind sets the API version and kind of the object reference.
func (obj *Object[_, _]) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	obj.APIVersion, obj.Kind = gvk.ToAPIVersionAndKind()
}

// GroupVersionKind returns the API version and kind of the object reference.
func (obj *Object[_, _]) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(obj.APIVersion, obj.Kind)
}

// GetObjectKind implements runtime.Object.
func (obj *Object[_, _]) GetObjectKind() schema.ObjectKind {
	return obj
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out.
func (in *Object[Spec, Status]) DeepCopyInto(out *Object[Spec, Status]) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec = out.Spec
	in.Status = out.Status
}

// DeepCopy copies the receiver, creating a new Object.
func (in *Object[Spec, Status]) DeepCopy() *Object[Spec, Status] {
	if in == nil {
		return nil
	}
	out := new(Object[Spec, Status])
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject implements runtime.Object.
func (in Object[_, _]) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
