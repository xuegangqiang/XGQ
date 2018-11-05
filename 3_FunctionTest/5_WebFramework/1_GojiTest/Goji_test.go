package main

import (
	"fmt"
	"net/http"
	"testing"

	"goji.io"
	"goji.io/pat"
)

func HelloInMain(w http.ResponseWriter, r *http.Request) {
	name := pat.Param(r, "name")
	fmt.Fprintf(w, "Hello, %s!", name)
}

func TestGoji(t *testing.T) {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/hello/:name"), HelloInMain)

	http.ListenAndServe("localhost:8000", mux)
}
