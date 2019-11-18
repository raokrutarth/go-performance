



# path to the package relative to ./src/
BENCHMARK_TARGET := examples/set

BENCHMARK_BINARY := benchmark

# set a default name for the container that is re-computed
# when the container is started
CONTAINER_NAME ?= go-performance_benchmark_1

package: run run-package

benchmark: run run-benchmark create-pprof-profiles

run:
	@docker-compose up --no-build --detach --remove-orphans
	$(eval CONTAINER_NAME := $(shell docker-compose ps -q benchmark))
	@printf "Benchmark container: %s\n" $(CONTAINER_NAME)

setup: clean
	@docker-compose build --parallel --force-rm
	-@mkdir ./profiles

run-package: setup-benchmark-run
	@docker exec -i $(CONTAINER_NAME) go build -v -o /bin/$(BENCHMARK_BINARY) $(BENCHMARK_TARGET)"
	@docker exec -i $(CONTAINER_NAME) $(BENCHMARK_BINARY)


run-benchmark: setup-benchmark-run
	# compile the test to a binary
	docker exec -i $(CONTAINER_NAME) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go test -c -i -o /bin/$(BENCHMARK_BINARY)"
	# run the test as a binary to allow process_exporter stats
	docker exec -i $(CONTAINER_NAME) $(BENCHMARK_BINARY) \
		-test.v \
		-test.bench=. \
		-test.benchmem \
		-test.memprofile=/profiles/memprofile.out \
		-test.cpuprofile=/profiles/cpuprofile.out \
		-test.mutexprofile=/profiles/mutexprofile.out \
		-test.blockprofile=/profiles/blockprofile.out

create-pprof-profiles:
	# generate pprof profile PDFs
	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines -sample_index=inuse_space \
		/bin/$(BENCHMARK_BINARY) /profiles/memprofile.out > ./profiles/inuse_heap.pdf

	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines -sample_index=alloc_space \
		/bin/$(BENCHMARK_BINARY) /profiles/memprofile.out > ./profiles/allocated_heap.pdf

	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines \
		/bin/$(BENCHMARK_BINARY) /profiles/cpuprofile.out > ./profiles/cpu.pdf

	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines \
		/bin/$(BENCHMARK_BINARY) /profiles/mutexprofile.out > ./profiles/mutex.pdf

	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines \
		/bin/$(BENCHMARK_BINARY) /profiles/blockprofile.out > ./profiles/block.pdf


setup-benchmark-run:
	-@docker exec -i $(CONTAINER_NAME) bash -c "pkill $(BENCHMARK_BINARY)"
	-@docker exec -i $(CONTAINER_NAME) bash -c "rm -rf /go/src/*"
	@docker cp ./src $(CONTAINER_NAME):/go
	@docker exec -i $(CONTAINER_NAME) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go get ./..."


clean:
	-@docker-compose down --remove-orphans

clean-collection-volumes:
	-@docker-compose down --volumes --remove-orphans