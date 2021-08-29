package export1

import (
	"context"
	"io"
	"math"
	"time"

	"github.com/efficientgo/examples/pkg/parquet-export/ref"
	"github.com/efficientgo/examples/pkg/parquet-export/ref/chunkenc"
	"github.com/pkg/errors"
	"github.com/xitongsys/parquet-go/writer"
	"google.golang.org/grpc"
)

// This "naive" version of export logic for Efficient Go book purposes.

var aggregationPeriod = int64((5 * time.Minute) / time.Millisecond) // Hardcoded 5 minutes.

// Export5mAggregations transforms selected data from Thanos system to Parquet format, suitable for analytic use.
func Export5mAggregations(ctx context.Context, address string, metricSelector []*LabelMatcher, minTime, maxTime int64, w io.Writer) (seriesNum int, samplesNum int, _ error) {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return 0, 0, errors.Wrap(err, "dial")
	}
	stream, err := NewStoreClient(cc).Series(ctx, &SeriesRequest{Matchers: metricSelector, MinTime: minTime, MaxTime: maxTime})
	if err != nil {
		return 0, 0, err
	}

	var series []Series
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, 0, errors.Wrap(err, "stream read")
		}

		if w := r.GetWarning(); w != "" {
			return 0, 0, errors.New(w)
		}
		series = append(series, *r.GetSeries())
		seriesNum++
	}

	pw, err := writer.NewParquetWriterFromWriter(w, new(ref.Aggregation), 1)
	if err != nil {
		return 0, 0, errors.Wrap(err, "new parquet writer")
	}

	var aggr []*ref.Aggregation
	for _, s := range series {
		a, sn, err := aggregate(s)
		if err != nil {
			return 0, 0, err
		}
		aggr = append(aggr, a...)
		samplesNum += sn
	}
	for _, a := range aggr {
		if err := pw.Write(a); err != nil {
			return 0, 0, err
		}
	}
	return seriesNum, samplesNum, pw.WriteStop()
}

func aggregate(s Series) (aggr []*ref.Aggregation, samplesNum int, _ error) {
	curr := newAggregationFromSeries(s.Labels)
	for _, c := range s.Chunks {
		r, err := chunkenc.FromData(chunkenc.Encoding(c.Raw.Type+1), c.Raw.Data)
		if err != nil {
			return nil, 0, err
		}
		iter := r.Iterator(nil)
		for iter.Next() {
			samplesNum++

			t, v := iter.At()

			if curr.Timestamp < t {
				aggr = append(aggr, curr)
				curr = newAggregationFromSeries(s.Labels)
			}
			if curr.Count == 0 {
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
			return nil, 0, err
		}
	}

	if curr.Count > 0 {
		aggr = append(aggr, curr)
	}

	return aggr, samplesNum, nil
}

// newAggregationFromSeries returns empty aggregation.
// For simplicity, we assume static labels in sorted order.
func newAggregationFromSeries(labels []*Label) *ref.Aggregation {
	return &ref.Aggregation{
		Timestamp: math.MaxInt64,
		Min:       math.MaxInt64,
		Max:       math.MinInt64,

		TargetLabel:  labels[0].Value,
		NameLabel:    labels[1].Value,
		ClusterLabel: labels[2].Value,
		ReplicaLabel: labels[3].Value,
	}
}
