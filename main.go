package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("✅ [GET] getMovies request sent...")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		log.Fatal("❌ JSON encoding error:", err)
		return
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("✅ [DELETE] deleteMovie request sent...")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			fmt.Println("\tDeleting movie with ID:", item.ID)
			// remove item by appending all other data in its place
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	// return remaining movies
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		log.Fatal("❌ JSON encoding error:", err)
		return
	}
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("✅ [GET] getMovie request sent...")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				log.Fatal("❌ JSON encoding error:", err)
				return
			}
			fmt.Println("\tFound movie with ID:", item.ID)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("✅ [POST] createMovie request sent...")
	w.Header().Set("Content-Type", "application/json")

	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	movie.ID = strconv.Itoa(rand.Intn(100000000))

	movies = append(movies, movie)

	err := json.NewEncoder(w).Encode(movie)
	if err != nil {
		log.Fatal("❌ JSON encoding error:", err)
		return
	}

	fmt.Println("\tCreated movie with ID:", movie.ID)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("✅ [PUT] updateMovie request sent...")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			fmt.Println("\tUpdating movie with ID:", item.ID)
			// delete move with the ID user has sent
			movies = append(movies[:index], movies[index+1:]...)

			createMovie(w, r)
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie 1",
		Director: &Director{Firstname: "John", Lastname: "Doe"}})

	movies = append(movies, Movie{ID: "2", Isbn: "438367", Title: "Movie 2",
		Director: &Director{Firstname: "Dr.", Lastname: "Strange"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server on port :8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
