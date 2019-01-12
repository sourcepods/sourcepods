package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gitpods/gitpods/authorization"
	"github.com/gitpods/gitpods/cmd"
	apiv1 "github.com/gitpods/gitpods/internal/api/v1"
	"github.com/gitpods/gitpods/repository"
	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/storage"
	"github.com/gitpods/gitpods/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/lib/pq"
	"github.com/oklog/run"
	prom "github.com/prometheus/client_golang/prometheus"
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
	Secret          string
	StorageGRPCURL  string
	StorageHTTPURL  string
	TracingURL      string
}

var (
	apiConfig = apiConf{}

	apiFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagAPIPrefix,
			EnvVar:      cmd.EnvAPIPrefix,
			Usage:       "The prefix the api is serving from, default: /",
			Value:       "/",
			Destination: &apiConfig.APIPrefix,
		},
		cli.StringFlag{
			Name:        cmd.FlagDatabaseDriver,
			EnvVar:      cmd.EnvDatabaseDriver,
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
			EnvVar:      cmd.EnvHTTPAddr,
			Usage:       "The address gitpods API runs on",
			Value:       ":3020",
			Destination: &apiConfig.HTTPAddr,
		},
		cli.StringFlag{
			Name:        cmd.FlagHTTPPrivateAddr,
			EnvVar:      cmd.EnvHTTPPrivateAddr,
			Usage:       "The address gitpods runs a http server only for internal access",
			Value:       ":3021",
			Destination: &apiConfig.HTTPPrivateAddr,
		},
		cli.BoolFlag{
			Name:        cmd.FlagLogJSON,
			EnvVar:      cmd.EnvLogJSON,
			Usage:       "The logger will log json lines",
			Destination: &apiConfig.LogJSON,
		},
		cli.StringFlag{
			Name:        cmd.FlagLogLevel,
			EnvVar:      cmd.EnvLogLevel,
			Usage:       "The log level to filter logs with before printing",
			Value:       "info",
			Destination: &apiConfig.LogLevel,
		},
		cli.StringFlag{
			Name:        cmd.FlagSecret,
			EnvVar:      cmd.EnvSecret,
			Usage:       "This secret is going to be used to generate cookies",
			Destination: &apiConfig.Secret,
		},
		cli.StringFlag{
			Name:        cmd.FlagStorageGRPCURL,
			EnvVar:      cmd.EnvStorageGRPCURL,
			Usage:       "The storage's gprc url to connect with",
			Destination: &apiConfig.StorageGRPCURL,
		},
		cli.StringFlag{
			Name:        cmd.FlagStorageHTTPURL,
			EnvVar:      cmd.EnvStorageHTTPURL,
			Usage:       "The storage's http url to proxy to",
			Destination: &apiConfig.StorageHTTPURL,
		},
		cli.StringFlag{
			Name:        cmd.FlagTracingURL,
			EnvVar:      cmd.EnvTracingURL,
			Usage:       "The url to send spans for tracing to",
			Destination: &apiConfig.TracingURL,
		},
	}
)

func apiAction(c *cli.Context) error {
	if apiConfig.Secret == "" {
		return errors.New("the secret for the api can't be empty")
	}

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
	us = user.NewLoggingService(log.WithPrefix(logger, "service", "user"), us)
	us = user.NewTracingService(us)

	var rs repository.Service
	rs = repository.NewService(repositories, storageClient)
	rs = repository.NewLoggingService(log.WithPrefix(logger, "service", "repository"), rs)
	rs = repository.NewTracingService(rs)

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
		router.Use(cmd.NewRequestLogger(logger))

		// Wrap the router inside a Router handler to make it possible to listen on / or on /api.
		// Change via APIPrefix.
		router.Route(apiConfig.APIPrefix, func(router chi.Router) {
			router.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write(page)
			})

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
				"msg", "starting gitpods api",
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
				"msg", "starting internal gitpods api",
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
	namespace := "gitpods"

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

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.0.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.6.1/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.6.1/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				graphQLParams.variables = graphQLParams.variables ? JSON.parse(graphQLParams.variables) : null;
				return fetch("/api/query", {
					method: "post",
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}

			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)

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
