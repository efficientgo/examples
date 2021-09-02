package ref

import (
	"fmt"
	"time"
)

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

type LabelMatcher_Type int32

const (
	LabelMatcher_EQ  LabelMatcher_Type = 0 // =
	LabelMatcher_NEQ LabelMatcher_Type = 1 // !=
	LabelMatcher_RE  LabelMatcher_Type = 2 // =~
	LabelMatcher_NRE LabelMatcher_Type = 3 // !~
)

func (l *LabelMatcher_Type) String() string {
	switch *l {
	case LabelMatcher_RE:
		return "=~"
	case LabelMatcher_NEQ:
		return "!~"
	case LabelMatcher_EQ:
		return "="
	case LabelMatcher_NRE:
		return "!="
	}
	return "unknown"
}

type LabelMatcher struct {
	Type  LabelMatcher_Type
	Name  string
	Value string
}

func (l *LabelMatcher) String() string {
	return fmt.Sprintf("%s%s'%s'", l.Name, l.Type.String(), l.Value)
}
