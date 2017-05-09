package main

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
)

const (
	FlagAddr     = "addr"
	FlagAddrAPI  = "addr-api"
	FlagEnv      = "env"
	FlagLogLevel = "loglevel"

	ProductionEnv = "production"
)

var FlagsUI = []cli.Flag{
	cli.StringFlag{
		Name:   FlagAddr,
		EnvVar: "GITPODS_ADDR",
		Usage:  "The address gitpods UI runs on",
		Value:  ":3010",
	},
	cli.StringFlag{
		Name:   FlagAddrAPI,
		EnvVar: "GITPODS_ADDR_API",
		Usage:  "The address gitpods API runs on",
		Value:  ":3020",
	},
	cli.StringFlag{
		Name:   FlagEnv,
		EnvVar: "GITPODS_ENV",
		Usage:  "The environment gitpod should run in",
		Value:  ProductionEnv,
	},
	cli.StringFlag{
		Name:   FlagLogLevel,
		EnvVar: "GITPODS_LOGLEVEL",
		Usage:  "The log level to filter logs with before printing",
		Value:  "info",
	},
}

func ActionUI(c *cli.Context) error {
	addr := c.String(FlagAddr)
	addrAPI := c.String(FlagAddrAPI)
	//env := c.String(FlagEnv)
	//loglevel := c.String(FlagLogLevel)

	//// Create the logger based on the environment: production/development/test
	//logger := newLogger(env, loglevel)

	// Create FileServer handler with buffalo's packr to serve file from disk or from within the binary.
	// The path is relative to this file.
	box := packr.NewBox("../../public")

	conf := HTMLConfig{API: addrAPI}

	r := NewUIRouter(box, conf)

	return http.ListenAndServe(addr, r)
}

func NewUIRouter(box packr.Box, conf HTMLConfig) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/", HomeHandler(box, conf)).Methods(http.MethodGet)
	r.Handle("/favicon.ico", http.FileServer(box)).Methods(http.MethodGet)
	r.Handle("/favicon.png", http.FileServer(box)).Methods(http.MethodGet)
	r.PathPrefix("/js").Handler(http.FileServer(box)).Methods(http.MethodGet)
	r.PathPrefix("/css").Handler(http.FileServer(box)).Methods(http.MethodGet)
	r.PathPrefix("/img").Handler(http.FileServer(box)).Methods(http.MethodGet)
	r.NotFoundHandler = HomeHandler(box, conf)

	return r
}

type HTMLConfig struct {
	API string `json:"api"`
}

func HomeHandler(box packr.Box, conf HTMLConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tem, err := template.New("index").Parse(box.String("index.html"))
		if err != nil {
			http.Error(w, "can't open index.html as template", http.StatusInternalServerError)
		}

		tem.Execute(w, conf)
	}
}

func newLogger(env string, loglevel string) log.Logger {
	var logger log.Logger

	if env == ProductionEnv {
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

	return log.With(logger, "ts", log.DefaultTimestampUTC)
}
