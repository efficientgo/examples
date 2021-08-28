package ref

import "time"

// Aggregation is simplified parquet format we could expect from TSDB data.
// In practice this ignored variability in labels, stale markers and counter notion of counter resets.
// This is ignored for example simplicity.
// NOTE: See https://github.com/xitongsys/parquet-go#example-of-type-and-encoding to understand the parquet struct tags.
type Aggregation struct {
	Name      string  `parquet:"name=__name__, type=BYTE_ARRAY"`
	Timestamp int64   `parquet:"name=_timestamp_millis, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
	Count     float64 `parquet:"name=_count, type=DOUBLE"`
	Sum       float64 `parquet:"name=_sum, type=DOUBLE"`
	Min       float64 `parquet:"name=_min, type=DOUBLE"`
	Max       float64 `parquet:"name=_max, type=DOUBLE"`
}

// TimestampFromTime returns a new millisecond timestamp from a time.
func TimestampFromTime(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/int64(time.Millisecond)
}

// TimeFromTimestamp returns a new time.Time object from a millisecond timestamp.
func TimeFromTimestamp(ts int64) time.Time {
	return time.Unix(ts/1000, (ts%1000)*int64(time.Millisecond)).UTC()
}
