.PHONY:  build run test db-up db-down db-reset db-clean

build: 
	go build -o bin/integ 

run:
	go run main.go

test:
	go test -v ./...

db-up:
	docker compose up -d

db-down:
	docker compose down

db-clean:
	docker compose down -v
	
db-reset:
	docker compose exec db psql -U postgres -d postgres -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"