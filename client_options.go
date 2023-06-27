// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022, The Zip API Object Framework Authors and Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.
package zip

import (
	"context"
	"reflect"

	"k8s.io/apiserver/pkg/storage"
)

// OnBefore is a prototype for a function that can be called on a referencable
// object before it is invoked by the client.
type OnBefore func(context.Context, ReferenceObject) (any, error)

// OnAfter is a prototype for a function that can be called on a referencable
// object after it has been returned by the client.
type OnAfter func(context.Context, ReferenceObject, ReferenceObject) (any, error)

// ClientConfig contains configuration which is passed to a Zip API client.
type ClientConfig struct {
	store  Store
	before []OnBefore
	after  []OnAfter
}

// ClientOption is a Zip API client config option-method handler.
type ClientOption func(*ClientConfig) error

// StoreRehydrationMode is a switch that dictates how the store's before call is
// used to manipulate the
type StoreRehydrationMode string

const (
	// Always rehydrate the reference object
	StoreRehydrationAlways = StoreRehydrationMode("always")

	// Only rehydratae the reference object when its Spec is nil
	StoreRehydrationSpecNil = StoreRehydrationMode("specnil")

	// Never rehydrate the reference object
	StoreRehydrationNever = StoreRehydrationMode("never")
)

// WithStore sets the Abstract client's store to the specified interface
// implementation.  An additional positional argument for the rehydration mode
// is used to configure when the store on the before call is used and how it
// manipulates the value of the reference object.
func WithStore[Spec, Status any](store Store, mode StoreRehydrationMode) ClientOption {
	return func(config *ClientConfig) error {
		if store == nil {
			return nil
		}

		config.store = store

		// Inject a before handler which saves the in-bound reference object to the
		// store.
		if err := WithBefore(func(ctx context.Context, req ReferenceObject) (any, error) {
			// If this object is listable, attempt to retrieve from a list from
			// the  store instead.
			if list, ok := req.(*ObjectList[Spec, Status]); ok {
				var ret ObjectList[Spec, Status]
				err := store.GetList(ctx, "", storage.ListOptions{}, &ret)
				if err != nil {
					return list, err
				}

				switch {
				// When hydration is set to never, return both the obj and err which is
				// then handled appropriately by the client.  Without this, the obj
				// would still be returned without the potential err.
				case mode == StoreRehydrationNever:
					return list, err

				// Always rehydrate
				case mode == StoreRehydrationAlways:
					if err == nil && &ret != nil {
						list.Items = ret.Items
					}

				case mode == StoreRehydrationSpecNil && (list.Items == nil || len(list.Items) == 0):
					if err == nil && &ret != nil {
						list.Items = ret.Items
					}
				}

				return list, nil
			}

			// Cast the referencable object, which we know is a spec-and-status
			// object.
			obj := req.(*Object[Spec, Status])

			// If this object is not referencable, do not attempt to retrieve
			// the object, simply return the input request.
			ref, err := req.Reference()
			if err != nil {
				return req, nil
			}

			var ret Object[Spec, Status]
			err = store.Get(ctx, ref, storage.GetOptions{}, &ret)

			switch {
			// When hydration is set to never, return both the obj and err which is
			// then handled appropriately by the client.  Without this, the obj would
			// still be returned without the potential err.
			case mode == StoreRehydrationNever:
				return obj, err

			// Always rehydrate.
			case mode == StoreRehydrationAlways:
				if err == nil && &ret != nil {
					*obj = ret
				}

			// Look up the object if the Spec of the Object is nil and a request to
			// hydrate the contents is desired.  This is useful in scenarios where
			// other attributes of the Object are used, e.g. those that define the
			// Reference() method of the ReferenceObject interface.
			case mode == StoreRehydrationSpecNil && reflect.DeepEqual(obj.Spec, *new(Spec)):
				if err == nil && &ret != nil {
					*obj = ret
				}
			}

			return obj, nil
		})(config); err != nil {
			return err
		}

		// Inject an after handler which saves the out-bound reference object to the
		// store.  If the out-bound reference object is nil, delete it from the
		// store.
		if err := WithAfter(func(ctx context.Context, before, after ReferenceObject) (any, error) {
			// If this object is not referencable, do not attempt to retrieve
			// the object, simply return the input request.
			ref, err := before.Reference()
			if err != nil {
				return after, nil
			}

			// If the returned object is empty, we should delete the reference
			if after == (*Object[Spec, Status])(nil) {
				return nil, store.Delete(ctx, ref, after, nil, nil, nil)
			}

			// Update the provided "before" reference with the returned "after" object
			if err := store.Create(ctx, ref, before, after, 0); err != nil {
				return after, err
			}

			return after, nil
		})(config); err != nil {
			return err
		}

		return nil
	}
}

// WithBefore provides pre-call functions which manipulate the inbound object
// before the client invokes its method strategy.
func WithBefore(before ...OnBefore) ClientOption {
	return func(config *ClientConfig) error {
		if config.before == nil {
			config.before = []OnBefore{}
		}
		config.before = append(config.before, before...)
		return nil
	}
}

// WithAfter provides post-call functions which manipulate the outbound object
// after the client has invoked its method strategy.
func WithAfter(after ...OnAfter) ClientOption {
	return func(config *ClientConfig) error {
		if config.after == nil {
			config.after = []OnAfter{}
		}
		config.after = append(config.after, after...)
		return nil
	}
}
