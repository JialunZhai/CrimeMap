package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	CrimeMapNamespace = "crime_map"
	HBaseClientSubsystem = "hbase_client"
)

var InvalidRowsCounter = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: CrimeMapNamespace,
	Subsystem: HBaseClientSubsystem,
	Name: "invalid_rows_counter",
	Help: "The total number of rows containing invalid values.",
})
