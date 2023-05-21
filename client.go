// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022, The Zip API Object Framework Authors and Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.
package zip

import "context"

// MethodClient is the generic Zip API client.
type MethodClient[In, Out ReferenceObject] struct {
	fn     Method[In, Out]
	config *ClientConfig
}

// NewMethodClient instantiates a new Zip API method client.
func NewMethodClient[In, Out ReferenceObject](
	ctx context.Context,
	fn Method[In, Out],
	opts ...ClientOption,
) (
	MethodStrategy[In, Out],
	error,
) {
	client := MethodClient[In, Out]{
		fn:     fn,
		config: &ClientConfig{},
	}

	for _, opt := range opts {
		if err := opt(client.config); err != nil {
			return nil, err
		}
	}

	return &client, nil
}

// Do implements MethodStrategy.
func (client *MethodClient[In, Out]) Do(ctx context.Context, req In) (Out, error) {
	for _, before := range client.config.before {
		ret, err := before(ctx, req)
		if err != nil {
			return *new(Out), err
		}

		req = ret.(In)
	}

	o, err := client.fn(ctx, req)
	if err != nil {
		return o, err
	}

	for _, after := range client.config.after {
		out, err := after(ctx, req, o)
		if err != nil {
			return *new(Out), err
		}

		if out, ok := out.(Out); ok {
			o = out
		}
	}

	return o, nil
}
