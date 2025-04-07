build:
	@go build -o bin/econestapi main.go

test:
	@go test -v ./...
	
run: build
	@./bin/econestapi

migration:
	@migrate create -ext sql -dir db/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run db/migrate/main.go up

migrate-down:
	@go run db/migrate/main.go down
