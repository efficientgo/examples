package ref

import "time"

// Aggregation is simplified parquet format we could expect from TSDB data.
// In practice this ignores variability in labels, stale markers and counter notion of counter resets.
// This is ignored for example simplicity.
// NOTE: See https://github.com/xitongsys/parquet-go#example-of-type-and-encoding to understand the parquet struct tags.
type Aggregation struct {
	NameLabel    string  `parquet:"name=__name__, type=BYTE_ARRAY, convertedtype=UTF8"`
	TargetLabel  string  `parquet:"name=__blockgen_target__, type=BYTE_ARRAY, convertedtype=UTF8"`
	ClusterLabel string  `parquet:"name=cluster, type=BYTE_ARRAY, convertedtype=UTF8"`
	ReplicaLabel string  `parquet:"name=replica, type=BYTE_ARRAY, convertedtype=UTF8"`
	Timestamp    int64   `parquet:"name=_timestamp_millis, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
	Count        int64   `parquet:"name=_count, type=INT64"`
	Sum          float64 `parquet:"name=_sum, type=DOUBLE"`
	Min          float64 `parquet:"name=_min, type=DOUBLE"`
	Max          float64 `parquet:"name=_max, type=DOUBLE"`
}

// TimestampFromTime returns a new millisecond timestamp from a time.
func TimestampFromTime(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/int64(time.Millisecond)
}

// TimeFromTimestamp returns a new time.Time object from a millisecond timestamp.
func TimeFromTimestamp(ts int64) time.Time {
	return time.Unix(ts/1000, (ts%1000)*int64(time.Millisecond)).UTC()
}
