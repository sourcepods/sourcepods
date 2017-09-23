package main

import (
	"html/template"
	"net/http"

	"github.com/gitpods/gitpods/cmd"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gobuffalo/packr"
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
			Name:        cmd.FlagAPIURL,
			EnvVar:      cmd.EnvAPIURL,
			Usage:       "The address gitpods API runs on",
			Value:       ":3020",
			Destination: &uiConfig.AddrAPI,
		},
		cli.StringFlag{
			Name:        cmd.FlagHTTPAddr,
			EnvVar:      cmd.EnvHTTPAddr,
			Usage:       "The address gitpods UI runs on",
			Value:       ":3010",
			Destination: &uiConfig.Addr,
		},
		cli.BoolFlag{
			Name:        cmd.FlagLogJSON,
			EnvVar:      cmd.EnvLogJSON,
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
	logger = log.WithPrefix(logger, "app", c.App.Name)

	// Create FileServer handler with buffalo's packr to serve file from disk or from within the binary.
	// The path is relative to this file.
	box := packr.NewBox("../../ui/build/web")

	homeHandler := HomeHandler(box, HTMLConfig{
		API: uiConfig.AddrAPI,
	})

	r := chi.NewRouter()
	r.Use(cmd.NewRequestLogger(logger))

	r.Get("/", homeHandler)
	r.NotFound(homeHandler)

	r.Handle("/components/*", http.FileServer(box))
	r.Handle("/favicon.ico", http.FileServer(box))
	r.Handle("/favicon.png", http.FileServer(box))
	r.Handle("/img/*", http.FileServer(box))
	r.Handle("/main.dart.js", http.FileServer(box))

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
