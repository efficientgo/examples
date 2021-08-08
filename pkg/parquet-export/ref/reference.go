package ref

// Aggregation is simplified parquet format we could expect from TSDB data.
// In practice this ignored variability in labels, stale markers and counter notion of counter resets.
// This is ignored for example simplicity.
type Aggregation struct {
	Name      string  `parquet:"name=__name__, type=BYTE_ARRAY"`
	Timestamp int64   `parquet:"name=,_timestamp_millis type=TIMESTAMP_MILLIS"`
	Count     float64 `parquet:"name=,_count type=DOUBLE"`
	Sum       float64 `parquet:"name=,_sum type=DOUBLE"`
	Min       float64 `parquet:"name=,_min type=DOUBLE"`
	Max       float64 `parquet:"name=,_max type=DOUBLE"`
}
