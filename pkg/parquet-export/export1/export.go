// Package export1, This "naive" version of export logic for Efficient Go book purposes.
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

var aggregationPeriod = int64((5 * time.Minute) / time.Millisecond) // Hardcoded 5 minutes.

// Export5mAggregations transforms selected data from Thanos system to Parquet format, suitable for analytic use.
func Export5mAggregations(ctx context.Context, address string, metricSelector []*ref.LabelMatcher, minTime, maxTime int64, w io.Writer) (seriesNum int, samplesNum int, _ error) {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return 0, 0, errors.Wrap(err, "dial")
	}

	stream, err := NewStoreClient(cc).Series(ctx, &SeriesRequest{Matchers: convertLabelMatchers(metricSelector), MinTime: minTime, MaxTime: maxTime})
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

	// Create 5m aggregations.
	var aggr []*ref.Aggregation
	for _, s := range series {
		var saggr []*ref.Aggregation

		curr := newAggregationFromSeries(s.Labels)
		if err := everySample(s, func(t int64, v float64) {
			samplesNum++
			if curr.Timestamp < t {
				saggr = append(saggr, curr)
				curr = newAggregationFromSeries(s.Labels)
			}
			if curr.Count == 0 {
				curr.Timestamp = t + aggregationPeriod
			}
			curr.Count++
		}); err != nil {
			return 0, 0, nil
		}
		if curr.Count > 0 {
			saggr = append(saggr, curr)
		}

		// Min.
		if err := everySampleAndAggr(s, saggr, func(t int64, v float64, aggr *ref.Aggregation) {
			if aggr.Min > v {
				aggr.Min = v
			}
		}); err != nil {
			return 0, 0, nil
		}

		// Max.
		if err := everySampleAndAggr(s, saggr, func(t int64, v float64, aggr *ref.Aggregation) {
			if aggr.Max < v {
				aggr.Max = v
			}
		}); err != nil {
			return 0, 0, nil
		}

		// Sum.
		if err := everySampleAndAggr(s, saggr, func(t int64, v float64, aggr *ref.Aggregation) {
			aggr.Sum += v
		}); err != nil {
			return 0, 0, nil
		}
		aggr = append(aggr, saggr...)
	}

	for _, a := range aggr {
		if err := pw.Write(a); err != nil {
			return 0, 0, err
		}
	}
	return seriesNum, samplesNum, pw.WriteStop()
}

func convertLabelMatchers(matchers []*ref.LabelMatcher) []*LabelMatcher {
	var ret []*LabelMatcher
	for _, m := range matchers {
		ret = append(ret, &LabelMatcher{
			Type:  LabelMatcher_Type(m.Type),
			Name:  m.Name,
			Value: m.Value,
		})
	}
	return ret
}

func everySample(s Series, f func(t int64, v float64)) error {
	for _, c := range s.Chunks {
		r, err := chunkenc.FromData(chunkenc.Encoding(c.Raw.Type+1), c.Raw.Data)
		if err != nil {
			return err
		}
		iter := r.Iterator(nil)
		for iter.Next() {
			f(iter.At())
		}
		if err := iter.Err(); err != nil {
			return err
		}
	}
	return nil
}

func everySampleAndAggr(s Series, aggr []*ref.Aggregation, f func(t int64, v float64, aggr *ref.Aggregation)) error {
	ai := -1
	cnt := int64(0)
	return everySample(s, func(t int64, v float64) {
		if cnt == 0 {
			ai++
			cnt = aggr[ai].Count
		}
		f(t, v, aggr[ai])
		cnt--
	})
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
