package cmd

import (
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const (
	FlagAPIPrefix       = "api-prefix"
	FlagAPIURL          = "api-url"
	FlagDatabaseDriver  = "database-driver"
	FlagDatabaseDSN     = "database-dsn"
	FlagGRPCAddr        = "grpc-addr"
	FlagHTTPAddr        = "http-addr"
	FlagHTTPPrivateAddr = "http-private-addr"
	FlagLogJSON         = "log-json"
	FlagLogLevel        = "log-level"
	FlagMigrationsPath  = "migrations-path"
	FlagRoot            = "root"
	FlagSSHAddr         = "ssh-addr"
	FlagSSHHostKeyPath  = "ssh-host-key"
	FlagStorageGRPCURL  = "storage-grpc-url"
	FlagStorageHTTPURL  = "storage-http-url"
	FlagTracingURL      = "tracing-url"

	//EnvDatabaseDSN is the data source name string to connect to the database with
	EnvDatabaseDSN = "GITPODS_DATABASE_DSN"
)

func NewLogger(json bool, loglevel string) log.Logger {
	var logger log.Logger

	if json {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	} else {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	}

	switch strings.ToLower(loglevel) {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	default:
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	return log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
}
