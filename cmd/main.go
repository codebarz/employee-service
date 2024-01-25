package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codebarz/employee-service/database"
	"github.com/codebarz/employee-service/entities/employees"
	"github.com/codebarz/employee-service/entities/roles"
	"github.com/codebarz/employee-service/rpc/proto/employeepb"
	"github.com/codebarz/employee-service/rpc/proto/rolepb"
	"github.com/codebarz/employee-service/services/employee"
	"github.com/codebarz/employee-service/services/role"
	roleservice "github.com/codebarz/employee-service/services/role"
	"github.com/go-kit/log"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	logger := initLogger()

	godotenv.Load()

	// r := chi.NewRouter()

	dbURL := os.Getenv("DB_URL")

	db := database.NewDatabase(logger)

	dbCfg := database.Config{DatabaseURL: dbURL}

	conn, err := db.OpenConnection(dbCfg)

	if err != nil {
		fmtErr := fmt.Sprintf("Err connecting to postgres DB. [ERROR]:%v", err)
		logger.Log(fmtErr)
		os.Exit(1)
	}

	connErr := conn.Ping()

	if connErr != nil {
		logger.Log(connErr)
	}

	logger.Log("DB connection successful")

	if err := db.Migrate(dbCfg); err != nil {
		logger.Log("Migration error: ", err)
		os.Exit(1)
	}

	defer func() {
		// db.l..Log("Closing DB connection", err)
		conn.Close()
	}()

	// v1Routes := chi.NewRouter()

	// v1Routes.Get("/health-check", services.Health)

	// r.Mount("/v1", v1Routes)

	// http.ListenAndServe(":9090", r)

	// rpc server
	rpc := grpc.NewServer()
	newRepo := roles.NewPgRepository(logger, conn)
	service := role.NewService(logger, newRepo)
	rs := roleservice.NewGRPCHandler(logger, service)

	employeeRepo := employees.NewPgRepository(logger, conn)
	employeeService := employee.NewService(logger, employeeRepo)
	es := employee.NewGRPCHandler(logger, employeeService)

	// register rpc server
	reflection.Register(rpc)
	rolepb.RegisterRoleServiceServer(rpc, rs)
	employeepb.RegisterEmployeeServiceServer(rpc, es)

	lis, err := net.Listen("tcp", ":9092")

	if err != nil {
		logger.Log("Can not listen")
		os.Exit(1)
	}

	rpc.Serve(lis)
}

func initLogger() log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	if os.Getenv("ENVIRONMENT") == "prod" || os.Getenv("ENVIRONMENT") == "stage" || os.Getenv("LOGFMT") == "json" {
		logger = log.NewJSONLogger(os.Stdout)
	}

	return logger
}
