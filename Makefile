



# path to the package relative to ./src/
BENCHMARK_TARGET := examples/readers

package: run run-package

run:
	@docker-compose up --no-build --detach --remove-orphans
	@$(eval CONTAINER_NAME := $(shell docker-compose ps -q benchmark))

setup: clean
	@docker-compose build --parallel --force-rm

run-package: stop-running-benchmark
	@docker cp ./src $(CONTAINER_NAME):/go
	@docker exec -i $(CONTAINER_NAME) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go get ./... && go build -o /benchmark ."
	@docker exec -d $(CONTAINER_NAME) bash -c "/benchmark > /var/log/benchmark.log"

stop-running-benchmark:
	-@docker exec -i $(CONTAINER_NAME) bash -c "pkill benchmark"


clean:
	-@docker-compose down --remove-orphans

clean-collection-volumes:
	-@docker-compose down --volumes --remove-orphans