# go-performance

Playground to measure performance implications of Golang code snippets.

## Prereqs

- Docker v??
- Docker-compose v??
- ability to run make file

## Usage

**Supported targets**

`test-file`

Given a _test.go file containing a Benchmark...() function, run the benchmark as a binary and dump pprof profiles.

`package`

Given a go package with a main() function, run the application on the container.
__Can import the telemetry package APIs to expose internal Golang runtime statistics or export internal numbers to the grafana dashboard.__

`file`

Given a go file/package with a main() function, run the application on the container.
__Can import the telemetry package APIs to expose internal Golang runtime statistics or export internal numbers to the grafana dashboard.__

### Benchmarking `_test.go` files

**If the benchmark or test can run for long enough** to see meaningful data on the system level charts, write a `_test.go` file (in a package if needed) and point Makefile to use the test file as below.

```Makefile
make test-file BENCHMARK_TARGET=src/example/benchmark_example_test.go
```

### Benchmarking Packages

Add the desired package containing a main function to src/ and set the package as a benchmark target as below.

```Makefile
make package BENCHMARK_TARGET=src/examples/telemetry
```


