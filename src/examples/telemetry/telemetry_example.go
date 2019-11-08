package main

import (
	"math/rand"
	"telemetry"
	"time"
)

func main() {
	telemetry.Initalize()

	rawMetricName := "example_guage"
	summaryMetricName := "example_summary"

	for {
		// use the telemetry API to expose a summary metric
		telemetry.ExportVariableValue(summaryMetricName, rand.Float64())

		// use the telemetry API to expose an in app counter
		telemetry.IncreaseRawValue(rawMetricName, 1)

		time.Sleep(5 * time.Second)
	}
}
