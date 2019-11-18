# go-performance

Playground to measure performance implications of Golang code snippets.

## Prereqs

- Docker
- docker-compose
- GNU Make

## Usage

**Supported targets**

`Go Benchmark test`

Given a package containing _test.go file containing at least one Benchmark...() function, run all benchmarks as a binary and dump pprof profiles.

`package main`

Given a go package with a main() function, run the application on the container.
__Can import the telemetry package APIs to expose internal Golang runtime statistics or export internal numbers to the grafana dashboard.__

### Benchmarking `_test.go` files

**If the benchmark or test can run for long enough** to see meaningful data on the system level charts, write a `_test.go` file (in a package if needed) and point Makefile to use the test file as below.

```Makefile
make setup (only once)
make test BENCHMARK_TARGET=example/set
```

### Benchmarking Packages

Add the desired package containing a main function to src/ and set the package as a benchmark target as below. See `src/examples/counters` for an example.

```Makefile
make setup (only once)
make main BENCHMARK_TARGET=examples/counters
```

The grafana UI should be visibile on port 3333 of the host machine


