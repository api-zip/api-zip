// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022, The Zip API Object Framework Authors and Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.
package zip

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	localSchemeBuilder = runtime.NewSchemeBuilder()
	Schemes            = runtime.NewScheme()
	Codecs             = serializer.NewCodecFactory(Schemes)
	ParameterCodec     = runtime.NewParameterCodec(Schemes)
)

// AddToScheme is a prototype to describe the method which can be used to append
// a new scheme to the list of global schemes.
type AddToScheme func(*runtime.Scheme) error

// Register accepts a slice of AddToScheme methods which are then registered
// against the list of global schemes.
func Register(schemes ...AddToScheme) error {
	for _, scheme := range schemes {
		localSchemeBuilder = append(localSchemeBuilder, scheme)
		if err := scheme(Schemes); err != nil {
			return err
		}
	}

	return nil
}
