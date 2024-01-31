package config

import (
	arg "github.com/alexflint/go-arg"
	"github.com/greyfinance/grey-go-libs/log"
	"github.com/pkg/errors"
)

// Config is the config struct
type Config struct {
	log.LoglevelEnv
	Environment        string `arg:"--env,env:ENVIRONMENT"`
	ListenHTTP         string `arg:"--listen-http,env:LISTEN_HTTP"`
	ListenGRPC         string `arg:"--listen-grpc,env:LISTEN_GRPC"`
	ListenHTTPLiveness string `arg:"--listen-http-liveness,env:LISTEN_HTTP_LIVENESS"`

	PostgresDBURL string `arg:"--postgres-db-url,env:POSTGRES_DB_URL"`
	DisableTLS    bool   `arg:"--disable-tls,env:DISABLE_TLS"`

	RabbitMQURL string `arg:"--rabbitmq-url,env:RABBITMQ_URL"`
}

// New creates a new config struct with sane defaults
func New() (Config, error) {
	c := Config{
		ListenHTTP:         ":8080",
		ListenGRPC:         ":8083",
		ListenHTTPLiveness: ":8084",
		PostgresDBURL:      "postgresql://127.0.0.1:5432/employee",
		DisableTLS:         true,
		RabbitMQURL:        "amqp://rabbitmq:rabbitmqpass@localhost:5672",
	}

	err := errors.Wrap(errors.WithStack(arg.Parse(&c)), "failed to parse config")
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
