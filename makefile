<<<<<<< HEAD
.PHONY:  build build-prod run test install db-up db-down db-clean bi fmt fmt-check vet integration-gen integration-test bench race lint
=======
.PHONY:  build build-prod run test install db-up db-down db-clean bi fmt fmt-check tidy tidy-check vulncheck vet integration-gen integration-test bench bench-sqlite bench-pg bench-all race lint test-sqlite test-pg test-dbs ci-local
>>>>>>> 2c6bd3d (fix: generate migrations folder on running migrations to avoid file not present err)

bi: build install

e2e: test integration-test

race:
	go test -race ./... && cd integration && go test -race ./...

bench:
	cd integration && go test -bench=. -benchmem -benchtime=2s -count=3

bench-sqlite:
	cd benchmark && make bench/sqlite && cd ..

bench-pg:
	cd benchmark && make bench/postgres && cd ..

bench-all:
	cd benchmark && make bench/all && cd ..

build: 
	go build -o bin/valk 

build-prod:
	go build -ldflags="-s -w" -o bin/valk

install: build
	mkdir -p $(HOME)/go/bin
	ln -sf $(shell pwd)/bin/valk $(HOME)/go/bin/valk

run:
	go build -o bin/valk && ./bin/valk

test:
	go test -v ./...

fmt:
	gofmt -w .

fmt-check:
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Go code is not formatted. Run: make fmt"; \
		gofmt -d .; \
		exit 1; \
	fi

vet:
	go vet ./...
	cd integration && go vet ./...

lint:
	-$(shell go env GOPATH)/bin/staticcheck ./...
	-$(shell go env GOPATH)/bin/gocritic check ./...

integration-gen: build
	cd integration && ../bin/valk generate

integration-test: integration-gen
	cd integration && go test -v ./...

db-up:
	docker compose up -d

db-down:
	docker compose down

db-clean:
	docker compose down -v