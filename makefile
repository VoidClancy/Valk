.PHONY:  build run test install db-up db-down db-clean bi fmt fmt-check vet sandbox-gen sandbox-test

bi: build install
	

build: 
	go build -o bin/valkyrie 

install: build
	mkdir -p $(HOME)/go/bin
	ln -sf $(shell pwd)/bin/valkyrie $(HOME)/go/bin/valkyrie

run:
	go build -o bin/valkyrie && ./bin/valkyrie

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
	cd sandbox && go vet ./...

sandbox-gen: build
	cd sandbox && ../bin/valkyrie generate

sandbox-test: sandbox-gen
	cd sandbox && go test -v ./...

db-up:
	docker compose up -d

db-down:
	docker compose down

db-clean:
	docker compose down -v