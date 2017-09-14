package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gitpods/gitpods/authorization"
	"github.com/gitpods/gitpods/cmd"
	"github.com/gitpods/gitpods/repository"
	"github.com/gitpods/gitpods/resolver"
	"github.com/gitpods/gitpods/session"
	"github.com/gitpods/gitpods/user"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/lib/pq"
	graphql "github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"github.com/oklog/oklog/pkg/group"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
)

type apiConf struct {
	Addr              string
	ListenAddrPrivate string
	APIPrefix         string
	DatabaseDriver    string
	DatabaseDSN       string
	LogJSON           bool
	LogLevel          string
	Secret            string
}

var (
	apiConfig = apiConf{}

	apiFlags = []cli.Flag{
		cli.StringFlag{
			Name:        cmd.FlagAddr,
			EnvVar:      cmd.EnvAddr,
			Usage:       "The address gitpods API runs on",
			Value:       ":3020",
			Destination: &apiConfig.Addr,
		},
		cli.StringFlag{
			Name:        cmd.FlagListenAddrPrivate,
			EnvVar:      cmd.EnvListenAddrPrivate,
			Usage:       "The address gitpods runs a http server only for internal access",
			Value:       ":3021",
			Destination: &apiConfig.ListenAddrPrivate,
		},
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
			Name:        cmd.FlagLogLevel,
			EnvVar:      cmd.EnvLogLevel,
			Usage:       "The log level to filter logs with before printing",
			Value:       "info",
			Destination: &apiConfig.LogLevel,
		},
		cli.BoolFlag{
			Name:        cmd.FlagLogJSON,
			EnvVar:      cmd.EnvLogJSON,
			Usage:       "The logger will log json lines",
			Destination: &apiConfig.LogJSON,
		},
		cli.StringFlag{
			Name:        cmd.FlagSecret,
			EnvVar:      cmd.EnvSecret,
			Usage:       "This secret is going to be used to generate cookies",
			Destination: &apiConfig.Secret,
		},
	}
)

func apiAction(c *cli.Context) error {
	if apiConfig.Secret == "" {
		return errors.New("the secret for the api can't be empty")
	}

	logger := cmd.NewLogger(apiConfig.LogJSON, apiConfig.LogLevel)
	logger = log.WithPrefix(logger, "app", "api")

	apiMetrics := apiMetrics()

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
	// Services
	//
	var ss session.Service
	ss = session.NewService(sessions)
	ss = session.NewMetricsService(ss, apiMetrics.SessionsCreated, apiMetrics.SessionsCleared)

	var as authorization.Service
	as = authorization.NewService(users.(authorization.Store), ss)
	as = authorization.NewLoggingService(log.WithPrefix(logger, "service", "authorization"), as)
	as = authorization.NewMetricsService(apiMetrics.LoginAttempts, as)

	var us user.Service
	us = user.NewService(users)
	us = user.NewLoggingService(log.WithPrefix(logger, "service", "user"), us)

	var rs repository.Service
	rs = repository.NewService(repositories)
	rs = repository.NewLoggingService(log.WithPrefix(logger, "service", "repository"), rs)

	//
	// Resolvers
	//
	res := &resolver.Resolver{
		resolver.NewUser(rs, us),
		resolver.NewRepository(rs, us),
	}

	schema, err := graphql.ParseSchema(resolver.Schema, res)
	if err != nil {
		panic(err)
	}

	//
	// Router
	//
	router := chi.NewRouter()
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

			router.Mount("/query", &relay.Handler{Schema: schema})
			router.Mount("/user", user.NewUserHandler(us))
			router.Mount("/users", user.NewUsersHandler(us))
			router.Mount("/users/{username}/repositories", repository.NewUsersHandler(rs))

			router.Mount("/repositories/{owner}/{name}", repository.NewHandler(rs))
		})
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

	server := &http.Server{
		Addr:    apiConfig.Addr,
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
		Addr:    apiConfig.ListenAddrPrivate,
		Handler: privateRouter,
	}

	var gr group.Group
	{
		gr.Add(func() error {
			dur := time.Minute
			level.Info(logger).Log("msg", "starting session cleaner", "interval", dur)
			for {
				if _, err := ss.ClearSessions(); err != nil {
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
				"addr", apiConfig.Addr,
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
				"addr", apiConfig.ListenAddrPrivate,
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
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.7.8/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.0.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.7.8/graphiql.js"></script>
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
