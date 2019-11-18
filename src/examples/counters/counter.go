package main

import (
	"math/rand"
	"time"

	"golang.performance.com/telemetry"
)

/**
	Simple example to demonstrate the use of the telemery API.
		- Increment a custom in-app counter every 5 seconds
		- Record a random value that will be added to a summary
**/

func main() {
	telemetry.Initialize()

	rawMetricName := "example_gauge"
	summaryMetricName := "example_summary"

	for {
		// use the telemetry API to expose a summary metric and observe a random value
		telemetry.ExportVariableValue(summaryMetricName, rand.Float64())

		// use the telemetry API to expose an in app counter
		telemetry.IncreaseRawValue(rawMetricName, 1)

		time.Sleep(5 * time.Second)
	}
}
