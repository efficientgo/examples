package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"io"
	stdlog "log"
	"net/http"
	httppprof "net/http/pprof"
	"os"
	"syscall"

	"github.com/efficientgo/examples/pkg/metrics/httpmidleware"
	"github.com/efficientgo/examples/pkg/sum"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/thanos-io/objstore"
	"github.com/thanos-io/objstore/client"
)

var (
	labelerFlags       = flag.NewFlagSet("labeler", flag.ExitOnError)
	addr               = labelerFlags.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	objstoreConfigYAML = labelerFlags.String("objstore.config", "", "Configuration YAML for object storage to label objects against")
)

func main() {
	if err := runMain(os.Args[1:]); err != nil {
		// Use %+v for github.com/pkg/errors error to print with stack.
		stdlog.Fatalf("Error: %+v", errors.Wrapf(err, "%s", flag.Arg(0)))
	}
}

func runMain(args []string) (err error) {
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

	tmpDir := "./tmp"
	if err := os.RemoveAll(tmpDir); err != nil {
		return errors.Wrap(err, "rm all")
	}
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "mkdir all")
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

		lbl, err := labelObjectNaive(ctx, tmpDir, bkt, objectIDs[0])
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

	m.HandleFunc("/debug/pprof/", httppprof.Index)
	m.HandleFunc("/debug/pprof/cmdline", httppprof.Cmdline)
	m.HandleFunc("/debug/pprof/profile", httppprof.Profile)
	m.HandleFunc("/debug/pprof/symbol", httppprof.Symbol)
	m.HandleFunc("/debug/pprof/trace", httppprof.Trace)

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
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))
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

func labelObjectNaive(ctx context.Context, tmpDir string, bkt objstore.BucketReader, objID string) (_ label, err error) {
	rc, err := bkt.Get(ctx, objID)
	if err != nil {
		return label{}, err
	}

	// Download file first.
	// TODO(bwplotka): This is naive for book purposes.
	f, err := os.CreateTemp(tmpDir, "cached-*")
	if err != nil {
		return label{}, err
	}
	defer func() { _ = os.RemoveAll(f.Name()) }()

	h := sha256.New()

	// Write to both checksum hash and file.
	if _, err := io.Copy(f, io.TeeReader(rc, h)); err != nil {
		return label{}, err
	}
	if err := rc.Close(); err != nil {
		return label{}, err
	}

	s, err := sum.Sum(f.Name())
	if err != nil {
		return label{}, err
	}

	// Get/calculate other attributes...

	return label{
		ObjID:    objID,
		Sum:      s,
		CheckSum: h.Sum(nil),
		// ...
	}, nil
}
