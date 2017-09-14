package main

// Much of this code originates from https://github.com/AaronO/go-git-http
// Licensed under Apache-2.0

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gitpods/gitpods/cmd"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
)

type GitHTTP struct {
	root string
	git  string
}

func NewGitHTTP(root string) *GitHTTP {
	return &GitHTTP{
		root: root,
		git:  "/usr/bin/git",
	}
}

func (gh *GitHTTP) Handler(logger log.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cmd.NewRequestLogger(logger))

	r.Get("/{owner}/{name}/HEAD", NoCaching(gh.textFileHandler("HEAD", "text/plain")))
	r.Get("/{owner}/{name}/info/refs", NoCaching(gh.infoRefsHandler))
	r.Get("/{owner}/{name}/objects/{folder:[0-9a-f]{2}}/{file:[0-9a-f]{38}}", CacheForever(gh.looseObjectHandler))
	r.Get("/{owner}/{name}/objects/info/{thing:[^/]*}", NoCaching(gh.infoHandler)) // TODO
	r.Get("/{owner}/{name}/objects/info/alternates", NoCaching(gh.textFileHandler("/objects/info/alternates", "text/plain")))
	r.Get("/{owner}/{name}/objects/info/http-alternates", NoCaching(gh.textFileHandler("/objects/info/http-alternates", "text/plain")))
	r.Get("/{owner}/{name}/objects/info/packs", implementHandler)
	r.Get("/{owner}/{name}/objects/pack/pack-[0-9a-f]{40}\\.idx", implementHandler)
	r.Get("/{owner}/{name}/objects/pack/pack-[0-9a-f]{40}\\.pack", implementHandler)
	r.Post("/{owner}/{name}/git-receive-pack", serviceHandler)
	r.Post("/{owner}/{name}/git-upload-pack", serviceHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		println("not found", r.URL.String())
	})

	return r
}

func implementHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusInternalServerError)
}

func NoCaching(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
		next.ServeHTTP(w, r)
	}
}

func CacheForever(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		expires := now.AddDate(1, 0, 0)
		w.Header().Set("Date", fmt.Sprintf("%d", now.Unix()))
		w.Header().Set("Expires", fmt.Sprintf("%d", expires.Unix()))
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		next.ServeHTTP(w, r)
	}
}

func (gh *GitHTTP) textFileHandler(path string, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Stat(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", f.Size()))
		w.Header().Set("Last-Modified", f.ModTime().Format(http.TimeFormat))
		http.ServeFile(w, r, path)
	}
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	fmt.Fprintf(w, "%s/%s", owner, name)
}

func (gh *GitHTTP) infoRefsHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	service := serviceQuery(r)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	args := []string{service, "--stateless-rpc", "--advertise-refs", "."}
	cmd := exec.CommandContext(ctx, gh.git, args...)
	cmd.Dir = filepath.Join(gh.root, owner, name)

	refs, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", service))
	w.Write(packetWrite(fmt.Sprintf("# service=git-%s\n", service)))
	w.Write(packetFlush())
	w.Write(refs)
}

func (gh *GitHTTP) looseObjectHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	folder := chi.URLParam(r, "folder")
	file := chi.URLParam(r, "file")
	path := filepath.Join(gh.root, owner, name, folder, file)

	gh.textFileHandler(path, "application/x-git-loose-object")
}

func (gh *GitHTTP) infoHandler(w http.ResponseWriter, r *http.Request) {
	owner, name := ownerName(r)
	thing := chi.URLParam(r, "thing")
	path := filepath.Join(gh.root, owner, name, "objects", "info", thing)

	fmt.Println(r.URL.String())
	gh.textFileHandler(path, "text/plain")

}

func ownerName(r *http.Request) (string, string) {
	return chi.URLParam(r, "owner"), chi.URLParam(r, "name")
}

func serviceQuery(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Query().Get("service"), "git-")
}

func packetFlush() []byte {
	return []byte("0000")
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)

	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}

	return []byte(s + str)
}
