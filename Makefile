# Build
build:
	go build cmd/main.go

# Run
run:
	go run cmd/main.go

migration-up:
	migrate -database ${POSTGRESQL_URL} -path database/migrations up

migration-down:
	migrate -database ${POSTGRESQL_URL} -path database/migrations down
