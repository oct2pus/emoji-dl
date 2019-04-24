package main // import "emojo"

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/oct2pus/emoji-dl/emoji"
	"log"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/{instance}", serveInstance).Methods("GET")
	r.HandleFunc("/{instance}/targz", serveTar4Instance).Methods("GET")
	r.HandleFunc("/{instance}/zip", serveZip4Instance).Methods("GET")
	r.HandleFunc("/health", serveHealth).Methods("GET")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8008", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "index.html")
	}
}

func serveInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	emojo, _ := emoji.NewCollection(vars["instance"])
	
	
	println(emojo)
}

func serveTar4Instance(w http.ResponseWriter, r *http.Request) {

}

func serveZip4Instance(w http.ResponseWriter, r *http.Request) {

}

// serveHealth is a simple HealthCheck
func serveHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}