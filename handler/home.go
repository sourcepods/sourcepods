package handler

import (
	"net/http"

	"github.com/gobuffalo/packr"
)

func HomeHandler(box packr.Box) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(box.Bytes("index.html"))
	}
}
