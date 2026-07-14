.PHONY:  build build-prod run test install db-up db-down db-clean bi fmt fmt-check tidy tidy-check vulncheck vet integration-gen integration-test bench race lint test-sqlite test-pg test-dbs ci-local

bi: build install

e2e: test integration-test

race:
	go test -race ./... && cd integration && go test -race ./...

bench:
	cd integration && go test -bench=. -benchmem -benchtime=2s -count=3

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

tidy:
	go mod tidy
	cd integration && go mod tidy

tidy-check:
	go mod tidy
	cd integration && go mod tidy
	@git diff --exit-code go.mod go.sum integration/go.mod integration/go.sum || (echo "go.mod or go.sum is not tidy. Run: make tidy"; exit 1)

vet:
	go vet ./...
	cd integration && go vet ./...

lint:
	-$(shell go env GOPATH)/bin/staticcheck ./...
	-$(shell go env GOPATH)/bin/gocritic check ./...

vulncheck:
	govulncheck ./...
	cd integration && govulncheck ./...

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
	
db-reset:
	docker compose exec db psql -U postgres -d postgres -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
 
test-sqlite: bi test
	node integration/prepareSchema.js sqlite
	cd integration && ../bin/valk -g
	rm -f integration/valk/migrations/*.sql
	rm -f integration/dev.db
	cd integration && DATABASE_URL="file:./dev.db" DATABASE_DIRECT_URL="file:./dev.db" ../bin/valk -m init
	cd integration && go test -tags sqlite -v ./...

test-pg: bi db-reset test
	node integration/prepareSchema.js postgres
	cd integration && ../bin/valk -g
	rm -f integration/valk/migrations/*.sql
	cd integration && DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" DATABASE_DIRECT_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" ../bin/valk -m init
	cd integration && go test -v ./...

test-dbs: test-sqlite test-pg

ci: fmt fmt-check tidy-check vet test vulncheck test-dbs