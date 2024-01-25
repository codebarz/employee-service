.PHONY: proto
proto: ## Compile GRPC Protobuf
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    rpc/proto/*/*.proto
	
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

