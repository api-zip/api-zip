// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022, The Zip API Object Framework Authors and Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.
package zip

import "context"

// StreamClient provides the interface for connecting to a remote service which
// returns a stream.
type StreamClient[In ReferenceObject, Out any] struct {
	fn     Stream[In, Out]
	config *ClientConfig
}

// NewStreamClient instantiates a new stream-enable client for the given stream
// provider.
func NewStreamClient[In ReferenceObject, Out any](
	ctx context.Context,
	fn Stream[In, Out],
	opts ...ClientOption,
) (
	StreamStrategy[In, Out],
	error,
) {
	client := StreamClient[In, Out]{
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

// Channel implements StreamStrategy
func (client *StreamClient[In, Out]) Channel(
	ctx context.Context,
	req In,
) (
	chan Out,
	chan error,
	error,
) {
	// Handle OnBefore requests
	for _, before := range client.config.before {
		ret, err := before(ctx, req)
		if err != nil {
			return nil, nil, err
		}

		req = ret.(In)
	}

	// Invoke the method which returns the channel
	events, errs, err := client.fn(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	// Return early if the returning object cannot be handled by OnAfter
	// interceptors (due to type constraint) or if there are no OnAfter handlers.
	if _, ok := any(*new(Out)).(ReferenceObject); !ok || len(client.config.after) == 0 {
		return events, errs, nil
	}

	// Handle OnAfter requests by intercepting the channel
	go func(events *chan Out, errs *chan error) {
	loop:
		for {
			select {
			case event := <-*events:
				in := any(event).(ReferenceObject)

				for _, after := range client.config.after {
					out, err := after(ctx, req, in)
					if err != nil {
						*errs <- err
						continue
					}

					// Re-assign the result from the OnAfter callback such that its value
					// is passed to the next OnAfter invocation.
					in = any(out).(ReferenceObject)
				}

				*events <- in.(Out)

			case <-ctx.Done():
				break loop
			}
		}
	}(&events, &errs)

	return events, errs, nil
}
