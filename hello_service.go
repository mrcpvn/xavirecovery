package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

const cl = "Content-Length"

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, true)
		log.Printf("request: %v", string(dump))
		h.ServeHTTP(w, r)
		//w.Write([]byte("ksdghfkjsdgfkjzdgghkszzxkjhcb######"))
		log.Printf("%v", w)
	})
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{Hello XAVI version1!}")
	}).Methods("GET").Headers("Accept", "application/json+v1")

	r.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{Hello XAVI version2!}")
	}).Methods("GET").Headers("Accept", "application/json+v2")

	r.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{no version header} panic!")
		//panic("No version specified!")
	}).Methods("GET")

	http.Handle("/hello/", Logger(r))

	http.ListenAndServe(":8080", nil)
}
