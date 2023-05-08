package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

const (
	CrimeMapNamespace    = "crime_map"
	HBaseClientSubsystem = "hbase_client"
)

var InvalidRowsCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: CrimeMapNamespace,
	Subsystem: HBaseClientSubsystem,
	Name:      "invalid_rows_counter",
	Help:      "The total number of rows containing invalid values.",
})

var ScanRequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: CrimeMapNamespace,
	Subsystem: HBaseClientSubsystem,
	Name:      "scan_requests_counter",
	Help:      "The total number of SCAN requests sent to HBase.",
})

var ReturnedRowsCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: CrimeMapNamespace,
	Subsystem: HBaseClientSubsystem,
	Name:      "scan_responses_counter",
	Help:      "The total number of rows returned by SCAN from HBase.",
})

var ScanRequestDurationSec = promauto.NewHistogram(prometheus.HistogramOpts{
	Namespace: CrimeMapNamespace,
	Subsystem: HBaseClientSubsystem,
	Name:      "scan_request_duration_sec",
	Buckets:   durationSecBuckets(10000*time.Microsecond, 1*time.Hour, 10),
	Help:      "How long it took to complete a SCAN request from HBase, in **seconds**.",
})

// exponentialBucketRange returns prometheus.ExponentialBuckets specified in
// terms of a min and max value, rather than needing to explicitly calculate the
// number of buckets.
func exponentialBucketRange(min, max, factor float64) []float64 {
	buckets := []float64{}
	current := min
	for current < max {
		buckets = append(buckets, current)
		current *= factor
	}
	return buckets
}

func durationSecBuckets(min, max time.Duration, factor float64) []float64 {
	return exponentialBucketRange(float64(min.Seconds()), float64(max.Seconds()), factor)
}
