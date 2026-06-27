.PHONY:  build run test

build: 
	go build -o valkyrie 

run:
	go build -o valkyrie && ./valkyrie

test:
	go test -v ./...