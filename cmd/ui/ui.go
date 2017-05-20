package main

import (
	"html/template"
	"net/http"

	"github.com/gitpods/gitpods/cmd"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gobuffalo/packr"
	"github.com/pressly/chi"
	"github.com/urfave/cli"
)

type uiConf struct {
	Addr     string
	AddrAPI  string
	LogJson  bool
	LogLevel string
}

var (
	uiConfig = uiConf{}

	uiFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagAddr,
			EnvVar:      cmd.EnvAddr,
			Usage:       "The address gitpods UI runs on",
			Value:       ":3010",
			Destination: &uiConfig.Addr,
		},
		cli.StringFlag{
			Name:        cmd.FlagAddrAPI,
			EnvVar:      cmd.EnvAddrAPI,
			Usage:       "The address gitpods API runs on",
			Value:       ":3020",
			Destination: &uiConfig.AddrAPI,
		},
		cli.BoolFlag{
			Name:        cmd.FlagLogJson,
			EnvVar:      cmd.EnvLogJson,
			Usage:       "The logger will log json lines",
			Destination: &uiConfig.LogJson,
		},
		cli.StringFlag{
			Name:        cmd.FlagLogLevel,
			EnvVar:      cmd.EnvLogLevel,
			Usage:       "The log level to filter logs with before printing",
			Value:       "info",
			Destination: &uiConfig.LogLevel,
		},
	}
)

func ActionUI(c *cli.Context) error {
	logger := cmd.NewLogger(uiConfig.LogJson, uiConfig.LogLevel)
	logger = log.WithPrefix(logger, "app", "ui")

	// Create FileServer handler with buffalo's packr to serve file from disk or from within the binary.
	// The path is relative to this file.
	box := packr.NewBox("../../public")

	homeHandler := HomeHandler(box, HTMLConfig{
		API: uiConfig.AddrAPI,
	})

	r := chi.NewRouter()
	//r.Use(handler.LoggerMiddleware(logger))

	r.Get("/", homeHandler)
	r.NotFound(homeHandler)

	r.Handle("/favicon.ico", http.FileServer(box))
	r.Handle("/favicon.png", http.FileServer(box))
	r.Handle("/img/*", http.FileServer(box))
	r.Handle("/css/*", http.FileServer(box))
	r.Handle("/js/*", http.FileServer(box))

	level.Info(logger).Log("msg", "starting gitpods ui", "addr", uiConfig.Addr)
	return http.ListenAndServe(uiConfig.Addr, r)
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
