



# path to the package relative to ./src/
BENCHMARK_TARGET := examples/set

BENCHMARK_BINARY := benchmark

package: run run-package

benchmark: run run-benchmark

run:
	@docker-compose up --no-build --detach --remove-orphans
	@sleep 2s
	@$(eval CONTAINER_NAME := $(shell docker-compose ps -q benchmark))

setup: clean
	@docker-compose build --parallel --force-rm
	-@mkdir ./profiles

run-package: setup-benchmark-run
	@docker exec -i $(CONTAINER_NAME) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go get ./... && go build -o /bin/$(BENCHMARK_BINARY) ."
	@docker exec -i $(CONTAINER_NAME) $(BENCHMARK_BINARY)


run-benchmark: setup-benchmark-run
	docker exec -i $(CONTAINER_NAME) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go test -c -i -o /bin/$(BENCHMARK_BINARY)"
	docker exec -i $(CONTAINER_NAME) $(BENCHMARK_BINARY) \
		-test.v \
		-test.bench=. \
		-test.benchmem -test.memprofile=/memprofile.out \
		-test.cpuprofile=/cpuprofile.out \
		-test.mutexprofile=/mutexprofile.out \
		-test.blockprofile=/blockprofile.out
	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines /bin/$(BENCHMARK_BINARY) /memprofile.out > ./profiles/memory.pdf
	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines -sample_index=inuse_space /bin/$(BENCHMARK_BINARY) /memprofile.out > ./profiles/in_use.pdf
	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines -sample_index=alloc_space /bin/$(BENCHMARK_BINARY) /memprofile.out > ./profiles/allocated.pdf

	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines /bin/$(BENCHMARK_BINARY) /cpuprofile.out > ./profiles/cpu.pdf
	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines /bin/$(BENCHMARK_BINARY) /mutexprofile.out > ./profiles/mutex.pdf
	docker exec -i $(CONTAINER_NAME) go tool pprof -pdf -lines /bin/$(BENCHMARK_BINARY) /blockprofile.out > ./profiles/block.pdf


setup-benchmark-run:
	-@docker exec -i $(CONTAINER_NAME) bash -c "pkill $(BENCHMARK_BINARY)"
	-@docker exec -i $(CONTAINER_NAME) bash -c "rm -rf /go/src/*"
	@docker cp ./src $(CONTAINER_NAME):/go
	@docker exec -i $(CONTAINER_NAME) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go get ./..."


clean:
	-@docker-compose down --remove-orphans

clean-collection-volumes:
	-@docker-compose down --volumes --remove-orphans