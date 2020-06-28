package main

import (
	"golang-mongodb-crud/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/people/", peopleHandler)
	log.Println("Server started")
	http.ListenAndServe(":9090", mux)
}

func peopleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlers.GetHandler(w, r)
	case "POST":
		handlers.PostHandler(w, r)
	case "DELETE":
		handlers.DeleteHandler(w, r)
	}
}
