package cmd

import (
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

const (
	FlagAddr           = "addr"
	FlagAddrAPI        = "addr-api"
	FlagDatabaseDriver = "database-driver"
	FlagDatabaseDSN    = "database-dsn"
	FlagLogJSON        = "log-json"
	FlagLogLevel       = "log-level"
	FlagSecret         = "secret"

	EnvAddr           = "GITPODS_ADDR"
	EnvAddrAPI        = "GITPODS_ADDR_API"
	EnvDatabaseDriver = "GITPODS_DATABASE_DRIVER"
	EnvDatabaseDSN    = "GITPODS_DATABASE_DSN"
	EnvLogJSON        = "GITPODS_LOG_JSON"
	EnvLogLevel       = "GITPODS_LOG_LEVEL"
	EnvSecret         = "GITPODS_SECRET"
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
