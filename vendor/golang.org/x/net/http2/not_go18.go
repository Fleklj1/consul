// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !go1.8

package http2

import "net/http"

func configureServer18(h1 *http.Server, h2 *Server) error {
	// No IdleTimeout to sync prior to Go 1.8.
	return nil
}

func shouldLogPanic(panicValue interface{}) bool {
	return panicValue != nil
}
