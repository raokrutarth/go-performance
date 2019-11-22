



# path to the package relative to ./src/
BENCHMARK_TARGET := examples/tree

BENCHMARK_BINARY := benchmark

# set a default name for the container that is re-computed
# when the container is started
CONTAINER_ID ?= UNDEFINED

setup: clean-collection-volumes
	@docker-compose build --parallel --force-rm
	@mkdir -p ./profiles


main: run run-main

test: run run-test create-pprof-profiles copy-profiles

run:
	@docker-compose up --no-build --detach --remove-orphans
	# wait until the benchmark service from docker-compose is an active container
	@while [ -z "$$(docker-compose ps -q benchmark)" ]; do sleep 2s; done

	# Get the dynamically assigned container ID
	$(eval CONTAINER_ID = `docker-compose ps -q benchmark`)
	@printf "Benchmark Container ID: %s\n" $(CONTAINER_ID)
	@printf "Benchmark target: %s" $(BENCHMARK_TARGET)

run-main: setup-benchmark-run
	@docker exec -i $(CONTAINER_ID) go build -v -o /bin/$(BENCHMARK_BINARY) $(BENCHMARK_TARGET)

	# run the benchmark binary within the container in detached mode
	@docker exec -i $(CONTAINER_ID) $(BENCHMARK_BINARY)


run-test: setup-benchmark-run
	# compile the test to a binary
	docker exec -i $(CONTAINER_ID) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go test -c -i -o /bin/$(BENCHMARK_BINARY)"
	# run the test as a binary to allow process_exporter stats
	docker exec -i $(CONTAINER_ID) $(BENCHMARK_BINARY) \
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
	docker exec -i $(CONTAINER_ID) go tool pprof -pdf -lines -sample_index=inuse_space \
		/bin/$(BENCHMARK_BINARY) /profiles/memprofile.out > ./profiles/inuse_heap.pdf

	docker exec -i $(CONTAINER_ID) go tool pprof -pdf -lines -sample_index=alloc_space \
		/bin/$(BENCHMARK_BINARY) /profiles/memprofile.out > ./profiles/allocated_heap.pdf

	docker exec -i $(CONTAINER_ID) go tool pprof -pdf -lines \
		/bin/$(BENCHMARK_BINARY) /profiles/cpuprofile.out > ./profiles/cpu.pdf

	docker exec -i $(CONTAINER_ID) go tool pprof -pdf -lines \
		/bin/$(BENCHMARK_BINARY) /profiles/mutexprofile.out > ./profiles/mutex.pdf

	docker exec -i $(CONTAINER_ID) go tool pprof -pdf -lines \
		/bin/$(BENCHMARK_BINARY) /profiles/blockprofile.out > ./profiles/block.pdf


# copy-profiles copies the raw, PDF profiles and binaries from the container
# to the host machine
copy-profiles:
	docker cp $(CONTAINER_ID):/profiles/. ./profiles
	docker cp $(CONTAINER_ID):/bin/$(BENCHMARK_BINARY) ./profiles

setup-benchmark-run:
	# kill the benchmark binary if it is already running
	-@docker exec -i $(CONTAINER_ID) bash -c "pkill $(BENCHMARK_BINARY)"

	# copy/overwrite Go sources to container
	@docker cp ./src/. $(CONTAINER_ID):/go/src

	# install dependencies, if needed, by the target package
	@docker exec -i $(CONTAINER_ID) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go get ./..."

clean:
	-@docker-compose down --remove-orphans

clean-collection-volumes:
	-@docker-compose down --volumes --remove-orphans