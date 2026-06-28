.PHONY:  build run test

build: 
	go build -o bin/valkyrie 

run:
	go build -o bin/valkyrie && ./bin/valkyrie

test:
	go test -v ./...