



# path to the package relative to ./src/
BENCHMARK_PACKAGE := examples/tree


BENCHMARK_CONTAINER := gobenchmarkc

package: collection run-package

run-package:
	@docker cp ./src $(BENCHMARK_CONTAINER):/go
	@docker exec -i $(BENCHMARK_CONTAINER) bash -c "cd /go/src/$(BENCHMARK_PACKAGE) && go get ./... && go build -o /benchmark ."
	@docker exec -d $(BENCHMARK_CONTAINER) bash -c "/benchmark > out.txt"

collection: clean-collection
	@docker-compose up -d

collection-with-build: clean-collection
	@docker-compose up --build -d

clean-collection:
	-@docker-compose down --remove-orphans


clean: clean-collection


clean-collection-all:
	-@docker-compose down --volumes --remove-orphans