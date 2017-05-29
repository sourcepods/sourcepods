package main

import (
	"net/http"

	"github.com/AaronO/go-git-http"
	"github.com/gitpods/gitpods/cmd"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pressly/chi"
)

func main() {
	logger := cmd.NewLogger(false, "debug")

	g := githttp.New("./dev/git/")
	g, err := g.Init()
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to initialize githttp",
			"err", err,
		)
	}

	g.EventHandler = eventHandler(logger)

	router := chi.NewRouter()
	router.Use(cmd.NewRequestLogger(logger))
	router.Mount("/", g)

	http.ListenAndServe(":3030", router)
}

func eventHandler(logger log.Logger) func(githttp.Event) {
	return func(ev githttp.Event) {
		level.Debug(logger).Log(
			"type", ev.Type,
			"branch", ev.Branch,
			"commit", ev.Commit,
			"last_commit", ev.Last,
		)
	}
}
