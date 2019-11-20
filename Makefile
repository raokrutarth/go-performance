



# path to the package relative to ./src/
BENCHMARK_TARGET := examples/tree

BENCHMARK_BINARY := benchmark

# set a default name for the container that is re-computed
# when the container is started
CONTAINER_NAME ?= go-performance_benchmark_1

check_var := "docker-compose ps -q benchmark"

main: run run-main

test: run run-test create-pprof-profiles copy-profiles

run:
	@docker-compose up --no-build --detach --remove-orphans
	# wait until benchmark service from docker-compose is a running container
	while [ -z "$$(docker-compose ps -q benchmark)" ]; do sleep 1s; done

	# Set the global variable to the container ID
	$(eval CONTAINER_NAME := $$(docker-compose ps -q benchmark))
	@printf "Benchmark Container ID: %s\n" $(CONTAINER_NAME)

setup: clean-collection-volumes
	@docker-compose build --parallel --force-rm
	-@mkdir ./profiles > /dev/null

run-main: setup-benchmark-run
	@docker exec -i $(CONTAINER_NAME) go build -v -o /bin/$(BENCHMARK_BINARY) $(BENCHMARK_TARGET)
	@docker exec -d $(CONTAINER_NAME) $(BENCHMARK_BINARY)


run-test: setup-benchmark-run
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
		-test.blockprofile=/profiles/blockprofile.out \
		-test.trace=/profiles/trace.out

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


# copy-profiles copies the raw, PDF profiles and binaries from the container
# to the host machine
copy-profiles:
	docker cp $(CONTAINER_NAME):/profiles/. ./profiles
	docker cp $(CONTAINER_NAME):/bin/$(BENCHMARK_BINARY) ./profiles

setup-benchmark-run:
	-@docker exec -i $(CONTAINER_NAME) bash -c "pkill $(BENCHMARK_BINARY)"
	@docker cp ./src/. $(CONTAINER_NAME):/go/src
	@docker exec -i $(CONTAINER_NAME) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go get ./..."

clean:
	-@docker-compose down --remove-orphans

clean-collection-volumes:
	-@docker-compose down --volumes --remove-orphans