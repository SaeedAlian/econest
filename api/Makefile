build:
	@go build -o bin/econestapi main.go

test:
	@ENV="test" go test -v ./...
	
run: build
	@ENV="devel" ./bin/econestapi

run-prod: build
	@ENV="prod" ./bin/econestapi

run-super-admin-cli: build
	@ENV="devel" ./bin/econestapi --cli

migration:
	@migrate create -ext sql -dir db/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run db/migrate/main.go up

migrate-down:
	@go run db/migrate/main.go down
