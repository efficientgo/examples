package export1

import (
	"context"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/efficientgo/examples/pkg/parquet-export/ref"
	"github.com/efficientgo/examples/pkg/parquet-export/ref/chunkenc"
	"github.com/pkg/errors"
	"github.com/xitongsys/parquet-go/writer"
	"google.golang.org/grpc"
)

var aggregationPeriod = int64((5 * time.Minute) / time.Millisecond) // Hardcoded 5 minutes.

func Export5mAggregations(ctx context.Context, address string, metricSelector []*LabelMatcher, minTime, maxTime int64, w io.Writer) error {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	stream, err := NewStoreClient(cc).Series(ctx, &SeriesRequest{Matchers: metricSelector, MinTime: minTime, MaxTime: maxTime})
	if err != nil {
		return err
	}

	var series []*Series
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if w := r.GetWarning(); w != "" {
			return errors.New(w)
		}
		fmt.Println(r.GetSeries().Labels)
		series = append(series, r.GetSeries())
	}

	pw, err := writer.NewParquetWriterFromWriter(w, new(ref.Aggregation), 1)

	var aggr []ref.Aggregation
	for _, s := range series {
		curr := newAggregation(s.Labels[0].Value)

		for _, c := range s.Chunks {
			r, err := chunkenc.FromData(chunkenc.Encoding(c.Raw.Type), c.Raw.Data)
			if err != nil {
				return err
			}
			iter := r.Iterator(nil)
			for iter.Next() {
				t, v := iter.At()

				if curr.Count == 0 {
					curr.Timestamp = t + aggregationPeriod
				} else if curr.Timestamp < t {
					aggr = append(aggr, curr)
					curr := newAggregation(s.Labels[0].Value)
					curr.Timestamp = t + aggregationPeriod
				}

				curr.Count++
				curr.Sum += v
				if curr.Min > v {
					curr.Min = v
				}
				if curr.Max < v {
					curr.Max = v
				}
			}
			if iter.Err() != nil {
				return err
			}
		}

		if curr.Count > 0 {
			aggr = append(aggr, curr)
		}
	}
	for _, a := range aggr {
		if err := pw.Write(a); err != nil {
			return err
		}
	}
	return pw.WriteStop()
}

func newAggregation(name string) ref.Aggregation {
	return ref.Aggregation{
		Name:      name,
		Timestamp: math.MinInt64,
		Min:       math.MaxInt64,
		Max:       math.MinInt64,
	}
}
