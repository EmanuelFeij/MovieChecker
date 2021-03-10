package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

const (
	listToSee = "Movies to see"
	listSeen  = "Seen Movies"
)

var client *redis.Client
var tmp *template.Template

func initRedis() *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	fmt.Println("Here")
	client.Del(listToSee)
	client.Del(listSeen)
	return client
}

func myNewRouter() *mux.Router {

	r := mux.NewRouter()
	r.HandleFunc("/seen", seenGet).Methods("GET")
	r.HandleFunc("/tosee", toSeePost).Methods("POST")
	r.HandleFunc("/tosee", toSeeGet).Methods("GET")

	return r
}

func seenGet(w http.ResponseWriter, r *http.Request) {
	movies, err := client.LRange(listSeen, 0, 10).Result()
	if err != nil {
		fmt.Println(err)
	}

	tmp.ExecuteTemplate(w, "seen.html", movies)
}

func toSeeGet(w http.ResponseWriter, r *http.Request) {
	movies, err := client.LRange(listToSee, 0, 10).Result()
	if err != nil {
		fmt.Println(err)
	}
	tmp.ExecuteTemplate(w, "to_see.html", movies)

}

func toSeePost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100)

	if r.PostForm.Get("movie") != "" {
		fmt.Println("kkk")
		movie := r.PostForm.Get("movie")
		client.LPush(listToSee, movie)
		http.Redirect(w, r, "/tosee", 302)
		return
	}

	if r.PostForm.Get("movies") != "" {
		value := r.Form["movies"][0]
		client.LPush(listSeen, value)
		client.LRem(listToSee, 1, value)
		http.Redirect(w, r, "/tosee", 302)
		return
	}
	http.Redirect(w, r, "/tosee", 302)
	return

}

func main() {
	tmp = template.Must(template.ParseGlob("templates/*.html"))
	client = initRedis()
	r := myNewRouter()
	fs := http.FileServer(http.Dir("./css/"))
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", fs))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
