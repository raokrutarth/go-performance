package telemetry

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	restEndpoint     = "/metrics"
	merticsServePort = ":3535"
	varNameTag       = "variable_name"
	valueNameTag     = "value_name"
)

var (
	summaryObjectives = map[float64]float64{0.5: 0.05, 0.95: 0.001}

	valueSummaries = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "value_summaries",
			Help:       "Median and outliers of exposed variables/values",
			Objectives: summaryObjectives,
		},
		[]string{varNameTag},
	)

	arbitraryValues = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "value_gauge",
			Help: "Value of arbitrary exposed in-app numbers",
		},
		[]string{valueNameTag},
	)
)

func Initialize() {
	prometheus.MustRegister(arbitraryValues)
	prometheus.MustRegister(valueSummaries)

	go func() {
		http.Handle(restEndpoint, promhttp.Handler())
		http.ListenAndServe(merticsServePort, nil)
	}()

	// sanity check for the telemetry dashboards
	testMetricUpdater()
}

func SetRawValue(valueName string, value float64) {
	arbitraryValues.WithLabelValues(valueName).Set(value)
}

// IncreaseRawValue increases/inatalizes a raw counter by delta. Visible on the
// grafana dashboard in the ___ graph
func IncreaseRawValue(valueName string, delta float64) {
	arbitraryValues.WithLabelValues(valueName).Add(delta)
}

// IncreaseRawValue decreases/inatalizes a raw counter by delta. Visible on the
// grafana dashboard in the ___ graph. (delta should be positive)
func DecreaseRawValue(valueName string, delta float64) {
	if delta < 0 {
		log.Printf("expected delta greater than 0 but got %.2f", delta)
		return
	}
	arbitraryValues.WithLabelValues(valueName).Sub(delta)
}

// ExportVariableValue exposes the median and outliers of a given value.
// Use when the same "value/number" identified by valueName can have
// different values throughout the system. e.g. function execution times
func ExportVariableValue(valueName string, value float64) {
	valueSummaries.WithLabelValues(valueName).Observe(value)
}

// testMetricUpdater allows the /metrics consumer something to
// keep updating
func testMetricUpdater() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop() // needed to free OS resources in case of a panic in the loop

	go func() {
		for range ticker.C {
			IncreaseRawValue("per_second_ticker", 1)
		}
	}()
}
