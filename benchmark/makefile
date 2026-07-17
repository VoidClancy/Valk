.PHONY: run test build bench benchNS bench/raw bench/valkyrie bench/gorm bench/short

run:
	go test -bench=. -benchmem -count=1

build:
	go build -o bin/bench

bench:
	go run ./cmd/benchsummary

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
