package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/efficientgo/e2e"
	e2edb "github.com/efficientgo/e2e/db"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/examples/pkg/sum/sumtestutil"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/go-kit/log"
	"github.com/pkg/errors"
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

func TestLabeler_Label(t *testing.T) {
	e, err := e2e.NewDockerEnvironment("labeler")
	testutil.Ok(t, err)
	t.Cleanup(e.Close)

	mon, err := e2emonitoring.Start(e)
	testutil.Ok(t, err)
	testutil.Ok(t, mon.OpenUserInterfaceInBrowser())

	minio := e2edb.NewMinio(e, "object-storage", bktName)
	testutil.Ok(t, e2e.StartAndWaitReady(minio))

	// Run program we want to test and benchmark.
	labeler := e2e.NewInstrumentedRunnable(e, "labeler").
		WithPorts(map[string]int{"http": 8080}, "http").
		Init(e2e.StartOptions{
			Image: "labeler:latest",
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
	testutil.Ok(t, uploadTestInput(minio, "test-input", 2e6))

	// Load test labeler from 1 clients with k6 and export result to Prometheus.
	testutil.Ok(t, runK6LoadTest(e, `
import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  http.get('http://`+labeler.InternalEndpoint("http")+`/label_object?object_id=test-input');
  sleep(0.1);
}`, 5, 120*time.Second, "todo"))

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//go func() {
	//	for ctx.Err() == nil {
	//		now := time.Now()
	//		r, err := http.Get("http://" + labeler.Endpoint("http") + "/label_object?object_id=test-input")
	//		testutil.Ok(t, err)
	//		testutil.Ok(t, assertResp(`{}`, r))
	//		fmt.Println("client latency:", time.Since(now).String())
	//	}
	//}()

	// Once done, wait for user input so user can explore the results in Prometheus UI.
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
}

func runK6LoadTest(e e2e.Environment, script string, clients int, duration time.Duration, endpoint string) error {
	k6 := e.Runnable("k6").Init(e2e.StartOptions{Command: e2e.NewCommandRunUntilStop(), Image: "grafana/k6:0.39.0"})
	if err := e2e.StartAndWaitReady(k6); err != nil {
		return err
	}
	return k6.Exec(e2e.NewCommand("/bin/sh", "-c", `cat << EOF | k6 run --vus 5 --duration 120s -
`+script+`
EOF`))
}

func assertResp(expected string, resp *http.Response) error {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected code, expected 200, got: %v", resp.Status)
	}
	if expected == string(b) {
		return errors.Errorf("unexpected response, expected %v, got: %v", expected, string(b))
	}
	return nil
}

func printResp(resp *http.Response) error {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Got", resp.Status, " with body", b)
	return nil
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
