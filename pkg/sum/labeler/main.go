// Copyright (c) Efficient Go Authors
// Licensed under the Apache License 2.0.

package main

import (
	"context"
	"encoding/json"
	"flag"
	stdlog "log"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"
	"syscall"

	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/examples/pkg/metrics/httpmidleware"
	"github.com/felixge/fgprof"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gobwas/pool/pbytes"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/thanos-io/objstore/client"
)

const (
	labelObject1 = "labelObject1"
	labelObject2 = "labelObject2"
	labelObject3 = "labelObject3"
	labelObject4 = "labelObject4"
)

var (
	labelerFlags       = flag.NewFlagSet("labeler-v1", flag.ExitOnError)
	addr               = labelerFlags.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	objstoreConfigYAML = labelerFlags.String("objstore.config", "", "Configuration YAML for object storage to label objects against")
	labelerFunction    = labelerFlags.String("function", "labelObjectNaive", "The function to use for labeling. labelObjectNaive, "+labelObject1+", "+labelObject2+", "+labelObject3+","+labelObject4)
)

func main() {
	if err := runMain(context.Background(), os.Args[1:]); err != nil {
		// Use %+v for github.com/efficientgo/core/errors error to print with stack.
		stdlog.Fatalf("Error: %+v", errors.Wrapf(err, "%s", flag.Arg(0)))
	}
}

func runMain(ctx context.Context, args []string) (err error) {
	if err := labelerFlags.Parse(args); err != nil {
		return err
	}

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		version.NewCollector("metrics"),
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	if *objstoreConfigYAML == "" {
		return errors.New("missing -objstore.config flag")
	}

	logger := log.NewLogfmtLogger(os.Stderr)
	bkt, err := client.NewBucket(logger, []byte(*objstoreConfigYAML), reg, "labeler")
	if err != nil {
		return errors.Wrap(err, "bucket create")
	}

	l := &labeler{bkt: bkt}
	var labelObjectFunc labelFunc
	switch *labelerFunction {
	case "labelObjectNaive":
		l.tmpDir = "./tmp"
		if err := os.RemoveAll(l.tmpDir); err != nil {
			return errors.Wrap(err, "rm all")
		}
		if err := os.MkdirAll(l.tmpDir, os.ModePerm); err != nil {
			return errors.Wrap(err, "mkdir all")
		}

		labelObjectFunc = l.labelObjectNaive
	case labelObject1:
		labelObjectFunc = l.labelObject1
	case labelObject2:
		l.pool.New = func() any { return []byte(nil) }
		labelObjectFunc = l.labelObject2
	case labelObject3:
		l.bucketedPool = pbytes.New(1e3, 10e6)
		labelObjectFunc = l.labelObject3
	case labelObject4:
		// Yolo.
		labelerSet := [4]*labeler{
			{bkt: bkt},
			{bkt: bkt},
			{bkt: bkt},
			{bkt: bkt},
		}
		var used [4]bool
		l := sync.Mutex{}

		labelObjectFunc = func(ctx context.Context, objID string) (label, error) {
			l.Lock()
			found := -1
			for i, u := range used {
				if u {
					continue
				}
				found = i
			}
			if found == -1 {
				l.Unlock()
				return label{}, errors.New("Did not expect more requests than 4 at the same time.")
			}
			used[found] = true
			l.Unlock()

			ret, err := labelerSet[found].labelObject4(ctx, objID)
			l.Lock()
			used[found] = false
			l.Unlock()
			return ret, err
		}
	default:
		return errors.Newf("unknown function %v", *labelerFunction)

	}

	metricMiddleware := httpmidleware.NewMiddleware(reg, nil)
	m := http.NewServeMux()
	m.Handle("/metrics", metricMiddleware.WrapHandler("/metric", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	)))
	m.HandleFunc("/label_object", metricMiddleware.WrapHandler("/label_object", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		if err := r.ParseForm(); err != nil {
			httpErrHandle(w, http.StatusInternalServerError, err)
		}

		objectIDs := r.Form["object_id"]
		if len(objectIDs) == 0 {
			httpErrHandle(w, http.StatusBadRequest, errors.New("object_id parameter is required"))
			return
		} else if len(objectIDs) > 1 {
			httpErrHandle(w, http.StatusBadRequest, errors.New("only one object_id parameter is required"))
			return
		}

		// TODO(bwplotka): Discard request body.

		lbl, err := labelObjectFunc(ctx, objectIDs[0])
		if err != nil {
			httpErrHandle(w, http.StatusInternalServerError, err)
			return
		}

		b, err := json.Marshal(&lbl)
		if err != nil {
			httpErrHandle(w, http.StatusInternalServerError, err)
			return
		}

		if _, err := w.Write(b); err != nil {
			httpErrHandle(w, http.StatusInternalServerError, err)
			return
		}
	})))

	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/fgprof/profile", fgprof.Handler().ServeHTTP)

	srv := http.Server{Addr: *addr, Handler: m}

	g := &run.Group{}
	g.Add(func() error {
		level.Info(logger).Log("msg", "starting HTTP server", "addr", *addr)
		if err := srv.ListenAndServe(); err != nil {
			return errors.Wrap(err, "starting web server")
		}
		return nil
	}, func(error) {
		if err := srv.Close(); err != nil {
			level.Error(logger).Log("msg", "failed to stop web server", "err", err)
		}
	})
	g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	return g.Run()
}

func httpErrHandle(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte("{ \"error\": \" " + err.Error() + "\"}"))
}

type label struct {
	ObjID    string `json:"object_id"`
	Sum      int64  `json:"sum"`
	CheckSum []byte `json:"checksum"`
}
