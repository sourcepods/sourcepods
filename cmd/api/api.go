package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/lib/pq"
	"github.com/oklog/run"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sourcepods/sourcepods/cmd"
	"github.com/sourcepods/sourcepods/pkg/api"
	apiv1 "github.com/sourcepods/sourcepods/pkg/api/v1"
	"github.com/sourcepods/sourcepods/pkg/authorization"
	"github.com/sourcepods/sourcepods/pkg/session"
	"github.com/sourcepods/sourcepods/pkg/sourcepods/repository"
	"github.com/sourcepods/sourcepods/pkg/sourcepods/user"
	"github.com/sourcepods/sourcepods/pkg/storage"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/urfave/cli"
)

type apiConf struct {
	HTTPAddr        string
	HTTPPrivateAddr string
	APIPrefix       string
	DatabaseDriver  string
	DatabaseDSN     string
	LogJSON         bool
	LogLevel        string
	StorageGRPCURL  string
	StorageHTTPURL  string
	TracingURL      string
}

var (
	apiConfig = apiConf{}

	apiFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagAPIPrefix,
			Usage:       "The prefix the api is serving from, default: /",
			Value:       "/",
			Destination: &apiConfig.APIPrefix,
		},
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDriver,
			Usage:       "The database driver to use: memory & postgres",
			Value:       "postgres",
			Destination: &apiConfig.DatabaseDriver,
		},
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDSN,
			EnvVar:      cmd.EnvDatabaseDSN,
			Usage:       "The database connection data",
			Destination: &apiConfig.DatabaseDSN,
		},
		cli.StringFlag{
			Name:        cmd.FlagHTTPAddr,
			Usage:       "The address SourcePods API runs on",
			Value:       ":3020",
			Destination: &apiConfig.HTTPAddr,
		},
		cli.StringFlag{
			Name:        cmd.FlagHTTPPrivateAddr,
			Usage:       "The address SourcePods runs a http server only for internal access",
			Value:       ":3021",
			Destination: &apiConfig.HTTPPrivateAddr,
		},
		cli.BoolFlag{
			Name:        cmd.FlagLogJSON,
			Usage:       "The logger will log json lines",
			Destination: &apiConfig.LogJSON,
		},
		cli.StringFlag{
			Name:        cmd.FlagLogLevel,
			Usage:       "The log level to filter logs with before printing",
			Value:       "info",
			Destination: &apiConfig.LogLevel,
		},
		cli.StringFlag{
			Name:        cmd.FlagStorageGRPCURL,
			Usage:       "The storage's gprc url to connect with",
			Destination: &apiConfig.StorageGRPCURL,
		},
		cli.StringFlag{
			Name:        cmd.FlagStorageHTTPURL,
			Usage:       "The storage's http url to proxy to",
			Destination: &apiConfig.StorageHTTPURL,
		},
		cli.StringFlag{
			Name:        cmd.FlagTracingURL,
			Usage:       "The url to send spans for tracing to",
			Destination: &apiConfig.TracingURL,
		},
	}
)

func apiAction(c *cli.Context) error {
	logger := cmd.NewLogger(apiConfig.LogJSON, apiConfig.LogLevel)
	logger = log.WithPrefix(logger, "app", c.App.Name)

	apiMetrics := apiMetrics()

	if apiConfig.TracingURL != "" {
		traceConfig := config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LocalAgentHostPort: apiConfig.TracingURL,
			},
		}

		traceCloser, err := traceConfig.InitGlobalTracer(c.App.Name)
		if err != nil {
			return err
		}
		defer traceCloser.Close()

		level.Info(logger).Log(
			"msg", "tracing enabled",
			"addr", apiConfig.TracingURL,
		)
	} else {
		level.Info(logger).Log("msg", "tracing is disabled, no url given")
	}

	//
	// Stores
	//
	var (
		repositories repository.Store
		sessions     session.Store
		users        user.Store
	)

	switch apiConfig.DatabaseDriver {
	default:
		db, err := sql.Open("postgres", apiConfig.DatabaseDSN)
		if err != nil {
			return err
		}
		defer db.Close()

		users = user.NewPostgresStore(db)
		sessions = session.NewPostgresStore(db)
		repositories = repository.NewPostgresStore(db)
	}

	//
	// Storage
	//
	storageClient, err := storage.NewClient(apiConfig.StorageGRPCURL)
	if err != nil {
		return err
	}

	githttp, err := NewGitHTTPProxy(apiConfig.StorageHTTPURL)
	if err != nil {
		return err
	}

	//
	// Services
	//
	var ss session.Service
	ss = session.NewService(sessions)
	ss = session.NewMetricsService(ss, apiMetrics.SessionsCreated, apiMetrics.SessionsCleared)
	ss = session.NewTracingService(ss)

	var as authorization.Service
	as = authorization.NewService(users.(authorization.Store), ss)
	as = authorization.NewLoggingService(log.WithPrefix(logger, "service", "authorization"), as)
	as = authorization.NewMetricsService(apiMetrics.LoginAttempts, as)
	as = authorization.NewTracingService(as)

	var us user.Service
	us = user.NewService(users)
	us = user.NewLoggingService(us, api.GetRequestID, log.WithPrefix(logger, "service", "user"))
	us = user.NewTracingService(us, api.GetRequestID)

	var rs repository.Service
	rs = repository.NewService(repositories, storageClient)
	rs = repository.NewLoggingService(rs, api.GetRequestID, log.WithPrefix(logger, "service", "repository"))
	rs = repository.NewTracingService(rs, api.GetRequestID)

	//
	// OpenAPI
	//
	openapi, err := apiv1.New(rs, us)
	if err != nil {
		return err
	}

	//
	// Router
	//
	router := chi.NewRouter()
	{
		router.Use(api.NewRequestID)
		router.Use(api.NewRequestLogger(logger))

		// Wrap the router inside a Router handler to make it possible to listen on / or on /api.
		// Change via APIPrefix.
		router.Route(apiConfig.APIPrefix, func(router chi.Router) {
			router.Mount("/authorize", authorization.NewHandler(as))

			router.Group(func(router chi.Router) {
				router.Use(session.Authorized(ss))
				router.Mount("/sessions", session.NewHandler(ss))
				router.Mount("/v1", middleware.NoCache(openapi.Handler))
			})

			router.Mount("/{owner}/{name}.git", githttp)
		})

		if apiConfig.APIPrefix != "/" {
			router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "API is available at ", apiConfig.APIPrefix)
			})
		}

		router.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"Not Found"}`))
		})
	}

	server := &http.Server{
		Addr:    apiConfig.HTTPAddr,
		Handler: router,
	}

	privateRouter := chi.NewRouter()
	privateRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, http.StatusText(http.StatusOK))
	})
	privateRouter.Mount("/metrics", prom.UninstrumentedHandler())
	privateRouter.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "0.0.0") // TODO: Return json
	})

	privateServer := &http.Server{
		Addr:    apiConfig.HTTPPrivateAddr,
		Handler: privateRouter,
	}

	var gr run.Group
	{
		gr.Add(func() error {
			dur := time.Minute
			level.Info(logger).Log("msg", "starting session cleaner", "interval", dur)
			for {
				if _, err := ss.DeleteExpired(context.TODO()); err != nil {
					return err
				}
				time.Sleep(dur)
			}
		}, func(err error) {
		})
	}
	{
		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "starting SourcePods API",
				"addr", apiConfig.HTTPAddr,
			)
			return server.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				level.Error(logger).Log(
					"msg", "failed to shutdown http server gracefully",
					"err", err,
				)
				return
			}
			level.Info(logger).Log("msg", "http server shutdown gracefully")
		})
	}
	{
		gr.Add(func() error {
			level.Info(logger).Log(
				"msg", "starting internal SourcePods API",
				"addr", apiConfig.HTTPPrivateAddr,
			)
			return privateServer.ListenAndServe()
		}, func(err error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			if err := privateServer.Shutdown(ctx); err != nil {
				level.Error(logger).Log(
					"msg", "failed to shutdown internal http server gracefully",
					"err", err,
				)
				return
			}
			level.Info(logger).Log("msg", "internal http server shutdown gracefully")
		})
	}

	return gr.Run()
}

type APIMetrics struct {
	LoginAttempts   metrics.Counter
	SessionsCreated metrics.Counter
	SessionsCleared metrics.Counter
}

func apiMetrics() *APIMetrics {
	namespace := "sourcepods"

	return &APIMetrics{
		LoginAttempts: prometheus.NewCounterFrom(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "authentication",
			Name:      "login_attempts_total",
			Help:      "Number of login attempts that succeeded and failed",
		}, []string{"status"}),
		SessionsCreated: prometheus.NewCounterFrom(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "sessions",
			Name:      "created_total",
			Help:      "Number of created sessions",
		}, []string{}),
		SessionsCleared: prometheus.NewCounterFrom(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: "sessions",
			Name:      "cleared_total",
			Help:      "Number of cleared sessions",
		}, []string{}),
	}
}

func NewGitHTTPProxy(storageURL string) (*httputil.ReverseProxy, error) {
	backend, err := url.Parse(storageURL)
	if err != nil {
		return nil, err
	}

	return NewSingleHostReverseProxy(backend), nil
}

// NewSingleHostReverseProxy returns a new ReverseProxy that routes
// URLs to the scheme, host, and base path provided in target. If the
// target's path is "/base" and the incoming request was for "/dir",
// the target request will be for /base/dir.
// NewSingleHostReverseProxy does not rewrite the Host header.
// To rewrite Host headers, use ReverseProxy directly with a custom
// Director policy.
func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		// Remove .git from the path before proxying
		req.URL.Path = strings.Replace(req.URL.Path, ".git", "", -1)

		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
