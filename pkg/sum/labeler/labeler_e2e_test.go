package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/efficientgo/core/testutil"
	"github.com/efficientgo/e2e"
	e2edb "github.com/efficientgo/e2e/db"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/examples/pkg/sum/sumtestutil"
	"github.com/go-kit/log"
	"github.com/thanos-io/objstore/client"
	"github.com/thanos-io/objstore/providers/s3"
	"gopkg.in/yaml.v3"
)

const bktName = "test"

func marshal(t testing.TB, i interface{}) string {
	t.Helper()

	b, err := yaml.Marshal(i)
	testutil.Ok(t, err)

	return string(b)
}

// TestLabeler_LabelObject runs interactive macro benchmark for `labeler` program.
// Prerequisites:
// * `docker` CLI and docker engine installed.
// * Run `make docker` from root project to build `labeler:latest` docker image.
// Read more in "Efficient Go"; Example 8-19, 8-20,
func TestLabeler_LabelObject(t *testing.T) {
	t.Skip("Comment this line if you want to run it - it's interactive test. Won't be useful in CI")

	e, err := e2e.NewDockerEnvironment("labeler")
	testutil.Ok(t, err)
	t.Cleanup(e.Close)

	// Start monitoring.
	mon, err := e2emonitoring.Start(e)
	testutil.Ok(t, err)
	testutil.Ok(t, mon.OpenUserInterfaceInBrowser())

	// Start storage.
	minio := e2edb.NewMinio(e, "object-storage", bktName)
	testutil.Ok(t, e2e.StartAndWaitReady(minio))

	// Run program we want to test and benchmark.
	labeler := e2e.NewInstrumentedRunnable(e, "labeler").
		WithPorts(map[string]int{"http": 8080}, "http").
		Init(e2e.StartOptions{
			Image:     "labeler:test",
			LimitCPUs: 4.0,
			Command: e2e.NewCommand(
				"/labeler",
				"-listen-address=:8080",
				"-objstore.config="+marshal(t, client.BucketConfig{
					Type: client.S3,
					Config: s3.Config{
						Bucket:    bktName,
						AccessKey: e2edb.MinioAccessKey,
						SecretKey: e2edb.MinioSecretKey,
						Endpoint:  minio.InternalEndpoint(e2edb.AccessPortName),
						Insecure:  true,
					},
				}),
			),
		})
	testutil.Ok(t, e2e.StartAndWaitReady(labeler))

	// Add test file.
	testutil.Ok(t, uploadTestInput(minio, "object1.txt", 2e6))

	// Start continuous profiling (not present in examples 8-19, 8-20).
	parca := e2e.NewInstrumentedRunnable(e, "parca").
		WithPorts(map[string]int{"http": 7070}, "http").
		Init(e2e.StartOptions{
			Image: "ghcr.io/parca-dev/parca:main-4e20a666",
			Command: e2e.NewCommand("/bin/sh", "-c",
				`cat << EOF > /shared/data/config.yml && \
    /parca --config-path=/shared/data/config.yml
object_storage:
  bucket:
    type: "FILESYSTEM"
    config:
      directory: "./data"
scrape_configs:
- job_name: "labeler"
  scrape_interval: "15s"
  static_configs:
    - targets: [ '`+labeler.InternalEndpoint("http")+`' ]
  profiling_config:
    pprof_config:
      fgprof:
        enabled: true
        path: /debug/fgprof/profile
        delta: true
EOF
`),
			User:      strconv.Itoa(os.Getuid()),
			Readiness: e2e.NewTCPReadinessProbe("http"),
		})
	testutil.Ok(t, e2e.StartAndWaitReady(parca))
	testutil.Ok(t, e2einteractive.OpenInBrowser("http://"+parca.Endpoint("http")))

	// Load test labeler from 1 clients with k6 and export result to Prometheus.
	k6 := e.Runnable("k6").Init(e2e.StartOptions{
		Command: e2e.NewCommandRunUntilStop(),
		Image:   "grafana/k6:0.39.0",
	})
	testutil.Ok(t, e2e.StartAndWaitReady(k6))

	url := fmt.Sprintf(
		"http://%s/label_object?object_id=object1.txt",
		labeler.InternalEndpoint("http"),
	)
	testutil.Ok(t, k6.Exec(e2e.NewCommand(
		"/bin/sh", "-c",
		`cat << EOF | k6 run -u 1 -d 5m -
import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
	const res = http.get('`+url+`');
	check(res, {
		'is status 200': (r) => r.status === 200,
		'response': (r) =>
			r.body.includes('{"object_id":"object1.txt","sum":6221600000,"checksum":"SUUr'),
	});
	sleep(0.5)
}
EOF`)))

	// Once done, wait for user input so user can explore the results in Prometheus UI and logs.
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}

func uploadTestInput(m e2e.Runnable, objID string, numLen int) error {
	bkt, err := s3.NewBucketWithConfig(log.NewNopLogger(), s3.Config{
		Bucket:    bktName,
		AccessKey: e2edb.MinioAccessKey,
		SecretKey: e2edb.MinioSecretKey,
		Endpoint:  m.Endpoint(e2edb.AccessPortName),
		Insecure:  true,
	}, "test")
	if err != nil {
		return err
	}

	b := bytes.Buffer{}
	if _, err := sumtestutil.CreateTestInputWithExpectedResult(&b, numLen); err != nil {
		return err
	}

	return bkt.Upload(context.Background(), objID, &b)
}
