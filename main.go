package main

import (
	// Import the gorilla/mux library we just installed
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Declare a new router
	r := mux.NewRouter()

	// This is where the router is useful, it allows us to declare methods that
	// this path will be valid for
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/hello", handler).Methods("GET")
	r.HandleFunc("/contact", handler2).Methods("GET")

	staticFileDirectory := http.Dir("./assets/") //point to static file
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET") //parses the url	and returns the file requested

	http.ListenAndServe(":8081", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!!")
}

func handler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Contact Me!")
}
