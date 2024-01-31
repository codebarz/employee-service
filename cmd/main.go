package main

import (
	"net/http"
	"os"

	"github.com/codebarz/employee-service/config"
	"github.com/codebarz/employee-service/database"
	"github.com/codebarz/employee-service/entities/employees"
	"github.com/codebarz/employee-service/entities/roles"
	"github.com/codebarz/employee-service/rpc/proto/employeepb"
	"github.com/codebarz/employee-service/rpc/proto/rolepb"
	"github.com/codebarz/employee-service/services/employee"
	"github.com/codebarz/employee-service/services/role"
	roleservice "github.com/codebarz/employee-service/services/role"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/greyfinance/grey-go-libs/http/handler/version"
	"github.com/greyfinance/grey-go-libs/http/httputils"
	"github.com/greyfinance/grey-go-libs/http/renderer"
	"github.com/greyfinance/grey-go-libs/http/starter"
	"github.com/greyfinance/grey-go-libs/log/levelfilter"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_kit "github.com/grpc-ecosystem/go-grpc-middleware/logging/kit"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

func main() {

	logger := initLogger()
	level.Info(logger).Log("msg", "starting application...")

	godotenv.Load()

	cfg, err := config.New()
	if err != nil {
		level.Error(logger).Log("config err", err)
		os.Exit(1)
	}

	dbCfg := database.Config{PostgresDBURL: cfg.PostgresDBURL, DisableTLS: cfg.DisableTLS}
	if err := database.Migrate(dbCfg); err != nil {
		level.Error(logger).Log("migration err:", err)
		os.Exit(1)
	}

	db, err := database.Open(dbCfg)
	if err != nil {
		level.Error(logger).Log("connecting to db err", err)
		os.Exit(1)
	}

	defer func() {
		level.Info(logger).Log("Database Stopping", err)
		db.Close()
	}()

	// HTTP Routes
	router := chi.NewRouter()
	// render := renderer.NewRenderer(logger)
	renderer.RegisterDefaultRoutes(router, logger)

	//Public route group
	router.Group(func(r chi.Router) {
		r.Use(
			middleware.Recoverer,
			func(next http.Handler) http.Handler {
				return &ochttp.Handler{
					IsPublicEndpoint: true,
					Handler:          next,
				}
			},
		)

		// r.With(renderer.NewPromLatencyExporter("http_request_duration_events_seconds")).
		// 	Mount("/api/v1/events", eventservice.NewHTTPHandler(eventSvc, render))
	})

	// Servers
	httpServer := httputils.NewServerWithDefaultTimeouts(logger)
	httpServer.Handler = router
	httpServer.Addr = cfg.ListenHTTP

	// For Kubernetes to handle graceful shutdown
	livenessServer := httputils.NewServerWithDefaultTimeouts(logger)
	livenessServer.Handler = http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("liveliness probe ok"))
	})
	livenessServer.Addr = cfg.ListenHTTPLiveness

	// gRPC server
	grpcOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_kit.UnaryServerInterceptor(logger, grpc_kit.WithLevels(grpc_kit.DefaultClientCodeToLevel)),
		)),
		grpc.StatsHandler(&ocgrpc.ServerHandler{IsPublicEndpoint: false}),
	}

	// rpc server
	newRepo := roles.NewPgRepository(logger, db)
	service := role.NewService(logger, newRepo)

	employeeRepo := employees.NewPgRepository(logger, db)
	employeeService := employee.NewService(logger, employeeRepo)

	grpcServer := grpc.NewServer(grpcOpts...)

	rolepb.RegisterRoleServiceServer(grpcServer, roleservice.NewGRPCHandler(logger, service))
	employeepb.RegisterEmployeeServiceServer(grpcServer, employee.NewGRPCHandler(logger, employeeService))
	reflection.Register(grpcServer)
	grpc_prometheus.Register(grpcServer)

	// start servers
	servers := starter.New().
		WithHTTP(httpServer, livenessServer).
		WithGRPC(grpcServer, cfg.ListenGRPC, nil)
	servers.Log = logger
	if err := servers.RunUntilInterrupt(); err != nil {
		level.Error(logger).Log("msg", "failed to start HTTP/gRPC servers", "err", err)
		os.Exit(1)
	}

}

func initLogger() log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	if os.Getenv("ENVIRONMENT") == "prod" || os.Getenv("ENVIRONMENT") == "stage" || os.Getenv("LOGFMT") == "json" {
		logger = log.NewJSONLogger(os.Stdout)
	}
	logger = levelfilter.FromEnv(logger)
	logger = log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,

		"commit", version.Commit,
	)
	return logger
}
