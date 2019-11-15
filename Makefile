



# path to the package relative to ./src/
BENCHMARK_TARGET := examples/readers


BENCHMARK_CONTAINER := gobenchmarkc
BENCHMARK_NETWORK := gobenchnet
BENCHMARK_IMAGE := gobenchi

package: run-package

setup: clean-collection
	docker-compose build --parallel --force-rm
	docker-compose run -d prometheus grafana


run:
	docker run \
		-d --rm \
		--cpus="3.0" \
		--network=$(BENCHMARK_NETWORK) \
		--network-alias=benchmark \
		--memory="4g" \
		$(BENCHMARK_IMAGE)


clean-container:
	docker rm -f $(BENCHMARK_CONTAINER)

run-package:
	@docker cp ./src $(BENCHMARK_CONTAINER):/go
	@docker exec -i $(BENCHMARK_CONTAINER) bash -c "cd /go/src/$(BENCHMARK_TARGET) && go get ./... && go build -o /benchmark ."
	@docker exec -d $(BENCHMARK_CONTAINER) bash -c "/benchmark > out.txt"



clean-collection:
	-@docker-compose down --remove-orphans


clean: clean-container


clean-collection-all:
	-@docker-compose down --volumes --remove-orphans