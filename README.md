# Employee Service

### Database Migrations

Run

```bash
make migration-up
```

Down migration

```bash
make migration-down
```

### Run service

```bash
make run
```

### Buld service

```bash
make build
```

### Genrate GRPC proto buffers

```bash
make proto
```

## Test Service

Can be tested using <a href="https://github.com/fullstorydev/grpcurl">`grpcurl`</a> CLI tool. For example to create employee

```bash
grpcurl --plaintext -d '{"first_name": "Tega", "last_name": "Oke", "role"
: "e8dfbedc-f369-4863-9e9f-ea4e3831f16a", "email": "tega@grey.co"}' localhost:9092 EmployeeService.CreateEmployee
```

or with <a href="https://learning.postman.com/docs/sending-requests/grpc/first-grpc-request/">postman</a>
