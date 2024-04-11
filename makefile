HAS_DOCKER_COMPOSE= $(shell which docker-compose)
export DB_TABLE := smartwallets
export DB_FORWARDER_TABLE := forwarders
export TEST_MONGO_DB_URL := mongodb://localhost
export DB_URL := postgres://rockside:password@localhost:5435/rockside?sslmode=disable

all: build

build: clean
	@mkdir bin
	go build -o bin

test:
	go test ./... -race

run: build
	./bin/tzkt

clean:
	rm -rf bin

.PHONY: build test run clean
