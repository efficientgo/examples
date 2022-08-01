package fd

import (
	"net/http"
	"net/http/pprof"

	"github.com/felixge/fgprof"
)

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
