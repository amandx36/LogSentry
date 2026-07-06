package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var LogsProcessed = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "logs_processed_total",
		Help: "Total processed logs",
	},
)

var ErrorLogs = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "logs_error_total",
		Help: "Total error logs",
	},
)

var WarnLogs = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "logs_warn_total",
		Help: "Total warning logs",
	},
)

var InfoLogs = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "logs_info_total",
		Help: "Total info logs",
	},
)

var UnknownLogs = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "logs_unknown_total",
		Help: "Total unknown logs",
	},
)

var DBInsert = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "db_insert_total",
		Help: "Database inserts",
	},
)

var WatcherEvents = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "watcher_events_total",
		Help: "Watcher events",
	},
)

var ParserFailures = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "parser_failures_total",
		Help: "Parser failures",
	},
)
var LiveEvents = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "log_write_events_total",
		Help: "Total write events detected by the watcher",
	},
)

var ReadFailures = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "read_failures_total",
		Help: "Total failures while reading appended log data",
	},
)
var FilesProcessed = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "files_processed_total",
		Help: "Total log files processed",
	},
)

var ActiveWorkers = promauto.NewGauge(
	prometheus.GaugeOpts{
		Name: "active_workers",
		Help: "Current active worker goroutines",
	},
)

var ParseDuration = promauto.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "parse_duration_seconds",
		Help:    "Time taken to parse a file",
		Buckets: prometheus.DefBuckets,
	},
)
