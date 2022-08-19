package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/efficientgo/e2e"
	e2edb "github.com/efficientgo/e2e/db"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/thanos-io/objstore/client"
	"github.com/thanos-io/objstore/providers/s3"
)

func TestLabeler_LabelObject_Cmp(t *testing.T) {
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

	// Add test file.
	testutil.Ok(t, uploadTestInput(minio, "object.10M.txt", 1e7))
	testutil.Ok(t, uploadTestInput(minio, "object.100M.txt", 1e8))

	labelers := map[string]e2e.Runnable{labelObject1: nil, labelObject2: nil, labelObject3: nil}
	for labelerFunc := range labelers {
		// Run program we want to test and benchmark.
		labelers[labelerFunc] = e2e.NewInstrumentedRunnable(e, labelerFunc).
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
					"-function="+labelerFunc,
				),
			})
	}

	// Start continuous profiling.
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
    - targets:
      - '`+labelers[labelObject1].InternalEndpoint("http")+`'
      - '`+labelers[labelObject2].InternalEndpoint("http")+`'
      - '`+labelers[labelObject3].InternalEndpoint("http")+`'
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

	for _, labelerFunc := range []string{labelObject1, labelObject2, labelObject3} {
		l := labelers[labelerFunc]

		testutil.Ok(t, e2e.StartAndWaitReady(l))

		url10M := fmt.Sprintf("http://%s/label_object?object_id=object.10M.txt", l.InternalEndpoint("http"))
		url100M := fmt.Sprintf("http://%s/label_object?object_id=object.100M.txt", l.InternalEndpoint("http"))

		testutil.Ok(t, k6.Exec(e2e.NewCommand(
			"/bin/sh", "-c",
			`cat << EOF | k6 run -u 4 -d 1m -
import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
	const res = http.get('`+url10M+`');
	check(res, {
		'is status 200': (r) => r.status === 200,
		'response': (r) =>
			r.body.includes('{"object_id":"object.10M.txt","sum":31108000000,"checksum":""'),
	});
	sleep(0.5)

	const res2 = http.get('`+url100M+`');
	check(res2, {
		'is status 200': (r) => r.status === 200,
		'response': (r) =>
			r.body.includes('{"object_id":"object100M.txt","sum":311080000000,"checksum":""'),
	});
	sleep(0.5)
}
EOF`)))

		testutil.Ok(t, l.Stop())
	}
	// Once done, wait for user input so user can explore the results in Prometheus UI and logs.
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}
