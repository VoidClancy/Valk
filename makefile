.PHONY:  build run test install db-up db-down db-clean bi

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

db-up:
	docker compose up -d

db-down:
	docker compose down

db-clean:
	docker compose down -v