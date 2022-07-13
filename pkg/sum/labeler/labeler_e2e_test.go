package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/efficientgo/e2e"
	e2edb "github.com/efficientgo/e2e/db"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/examples/pkg/sum/sumtestutil"
	"github.com/efficientgo/tools/core/pkg/testutil"
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

func TestLabeler_Label(t *testing.T) {
	e, err := e2e.NewDockerEnvironment("labeler")
	testutil.Ok(t, err)
	t.Cleanup(e.Close)

	mon, err := e2emonitoring.Start(e)
	testutil.Ok(t, err)

	minio := e2edb.NewMinio(e, "object-storage", bktName)
	testutil.Ok(t, e2e.StartAndWaitReady(minio))

	// Run program we want to test and benchmark.
	labeler := e2e.NewInstrumentedRunnable(e, "labeler").
		WithPorts(map[string]int{
			"http": 8080,
		}, "http").
		Init(e2e.StartOptions{
			Image: "labeler:latest",
			Command: e2e.NewCommand(
				"/labeler",
				"-listen-address=8080",
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

	// Schedule constant traffic from single client.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for ctx.Err() == nil {
			r, err := http.Get("http://" + labeler.Endpoint("http") + "/label_object?object_id=test-input")
			testutil.Ok(t, err)
			testutil.Ok(t, printResp(r))
		}
	}()

	testutil.Ok(t, mon.OpenUserInterfaceInBrowser())
	testutil.Ok(t, e2einteractive.RunUntilEndpointHit())
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
