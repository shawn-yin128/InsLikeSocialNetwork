package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"around/model"
	"around/service"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) { // request uses pointer to save memory cause not copy; ResponseWriter is an interface not a struct, not support pointer
	// Parse from body of request to get a json object.
	fmt.Println("Received one upload request")
	decoder := json.NewDecoder(r.Body)
	var p model.Post
	if err := decoder.Decode(&p); err != nil {
		panic(err)
	}

	fmt.Fprintf(w, "Post received: %s\n", p.Message)
}

func searchHandler(w http.ResponseWriter, r *http.Request) { // go will pack request and response
	fmt.Println("Received one research request")
	w.Header().Set("Content-Type", "application/json") // set the reponse type

	// get value based on the key from URL
	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var posts []model.Post
	var err error
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else { // cover no user no keywords case, and it will go to the ZeroTermsQuery("all") case
		posts, err = service.SearchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to read post from backend %v.\n", err)
		return
	}

	js, err := json.Marshal(posts) // encoding json
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
		return
	}
	w.Write(js)
}
