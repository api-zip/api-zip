// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022, The Zip API Object Framework Authors and Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.
package zip

import "context"

// Method represents a generic API function call.
type Method[In, Out ReferenceObject] func(context.Context, In) (Out, error)

// MethodStrategy is a generic method handler which standardizes implementation
// of a client such that the action performed on a specific endpoint for a given
// input type In which returns a known output type Out is consistently
// implemented.
type MethodStrategy[In, Out ReferenceObject] interface {
	// Do performs the request by invoking the provided generic Method.
	Do(context.Context, In) (Out, error)
}

// Stream represents a generic stream-enabled API function call.
type Stream[In ReferenceObject, Out any] func(context.Context, In) (chan Out, chan error, error)

// StreamMethodStrategy is a generic handler which standardizes the
// implementation of a streaming client which performs actions on an input In
// which resolves to a stream (implemented as a Go channel) Out.
type StreamStrategy[In ReferenceObject, Out any] interface {
	// Channel performs the request with the provide input In and returns a
	// channel which allows returning a stream of content.  Any issues from the
	// stream itself will be propagated back through the error channel.
	// Initialization errors are returned through the standard error.
	Channel(context.Context, In) (chan Out, chan error, error)
}
