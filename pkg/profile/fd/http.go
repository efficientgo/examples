// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package fd

import (
	"net/http"
	"net/http/pprof"

	"github.com/felixge/fgprof"
)

// ExampleHTTP is an example of HTTP exposure of Go profiles.
// Read more in "Efficient Go"; Example 9-5.
func ExampleHTTP() {
	m := http.NewServeMux()
	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/fgprof/profile", fgprof.Handler().ServeHTTP)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	srv := http.Server{Handler: m}

	// Start server...

	_ = srv
}
