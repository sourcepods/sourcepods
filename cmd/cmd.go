package cmd

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pressly/chi/middleware"
)

const (
	FlagAddr           = "addr"
	FlagAddrAPI        = "addr-api"
	FlagAPIPrefix      = "api-prefix"
	FlagDatabaseDriver = "database-driver"
	FlagDatabaseDSN    = "database-dsn"
	FlagLogJSON        = "log-json"
	FlagLogLevel       = "log-level"
	FlagSecret         = "secret"

	EnvAddr           = "GITPODS_ADDR"
	EnvAddrAPI        = "GITPODS_ADDR_API"
	EnvAPIPrefix      = "GITPODS_API_PREFIX"
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

func NewRequestLogger(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			level.Debug(logger).Log(
				"proto", r.Proto,
				"method", r.Method,
				"status", ww.Status(),
				"path", r.URL.Path,
				"duration", time.Since(start),
				"bytes", ww.BytesWritten(),
			)
		})
	}
}
