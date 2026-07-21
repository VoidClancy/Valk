.PHONY: run test build bench benchNS bench/raw bench/valkyrie bench/gorm bench/short bench/stress bench/prll bench/stress-gc bench/sqlite bench/postgres bench/all

run:
	go test -bench=. -benchmem -count=1

build:
	go build -o bin/bench

bench:
	go run ./cmd/benchsummary
db-reset:
	docker compose exec db psql -U postgres -d postgres -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

bench/sqlite:
	node prepareSchema.js sqlite
	rm -rf valk/migrations/*.sql
	rm -f dev.db
	DATABASE_URL="file:./dev.db" DATABASE_DIRECT_URL="file:./dev.db" ../bin/valk migrate init
	../bin/valk generate
	go generate ./ent
	BENCH_DB=sqlite go run ./cmd/benchsummary

bench/pg:
	node prepareSchema.js postgres
	rm -rf valk/migrations/*.sql
	make db-reset
	DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" DATABASE_DIRECT_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" ../bin/valk migrate init
	../bin/valk generate
	go generate ./ent
	BENCH_DB=postgres go run ./cmd/benchsummary

bench/all:
	@echo "=== Running SQLite Benchmarks ==="
	@make bench/sqlite
	@echo "=== Running PostgreSQL Benchmarks ==="
	@make bench/pg

benchNS:
	go test -bench=. -benchmem -count=1 -timeout 30m

bench/raw:
	go test -bench=BenchmarkRaw -benchmem -count=1

bench/valkyrie:
	go test -bench=BenchmarkValkyrie -benchmem -count=1

bench/gorm:
	go test -bench=BenchmarkGORM -benchmem -count=1

bench/short:
	go test -bench=. -benchmem -benchtime=100ms -count=1

NPROC   := $(shell getconf _NPROCESSORS_ONLN 2>/dev/null || echo 1)
PRLL_CPUS := $(shell awk 'BEGIN{n=$(NPROC);r="1";if(n>=2)r=r",2";if(n>=4)r=r",4";if(n>=8)r=r",8";if(n>8)r=r","n;print r}')

STRESS_COUNT     ?= 5
STRESS_BENCHTIME ?= 1m
STRESS_TIMEOUT   ?= 1h
PRLL_BENCHTIME   ?= 10s

bench/stress:
	go run ./cmd/benchsummary -count=$(STRESS_COUNT) -benchtime=$(STRESS_BENCHTIME) -timeout=$(STRESS_TIMEOUT)

bench/prll:
	go run ./cmd/benchsummary -count=1 -benchtime=$(PRLL_BENCHTIME) -cpu=$(PRLL_CPUS)

bench/stress-gc:
	GOGC=1 go run ./cmd/benchsummary -count=$(STRESS_COUNT) -benchtime=10s -timeout=$(STRESS_TIMEOUT)