package main

import (
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

var client *redis.Client
var tmp *template.Template

func myNewRouter() *mux.Router {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r := mux.NewRouter()
	r.HandleFunc("/seen", seenGet).Methods("GET")
	//r.HandleFunc("/seen", seenPost).Methods("POST")
	r.HandleFunc("/tosee", toSeeGet).Methods("GET")

	return r
}

func seenGet(w http.ResponseWriter, r *http.Request) {
	movies, err := client.LRange("movie", 0, 10).Result()
	if err != nil {
		return
	}
	tmp.ExecuteTemplate(w, "seen.html", movies)
}

func toSeeGet(w http.ResponseWriter, r *http.Request) {
	tmp.ExecuteTemplate(w, "to_see.html", nil)

}

func main() {
	tmp = template.Must(template.ParseGlob("templates/*.html"))
	r := myNewRouter()
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
