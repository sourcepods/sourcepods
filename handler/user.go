package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func UserList(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
	fmt.Fprintln(w, "users")
}

func UserCreate(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
	fmt.Fprintln(w, "user created")
}

func User(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
	fmt.Fprintln(w, "user", id)
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
	fmt.Fprintf(w, "user %s updated", id)
}

func UserDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
	fmt.Fprintf(w, "user %s deleted", id)
}
