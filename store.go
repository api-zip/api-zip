// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022, The Zip API Object Framework Authors and Unikraft GmbH.
// Licensed under the BSD-3-Clause License (the "License").
// You may not use this file except in compliance with the License.
package zip

import "k8s.io/apiserver/pkg/storage"

// Store is an alias for the underlying implementation.
type Store storage.Interface
