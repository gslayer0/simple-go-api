package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var Articles []Article

func generateErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Message: message}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello! welcome to my homepage")
	fmt.Println("Endpoint hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	id, err := strconv.Atoi(key)
	if err != nil {
		errorResp := generateErrorResponse("Wrong request format")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	found := false

	for _, article := range Articles {
		if article.ID == id {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(article)
			found = true
			return
		}
	}
	if !found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(generateErrorResponse("No matched article found"))
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)

	newID := len(Articles) + 1
	article.ID = newID

	Articles = append(Articles, article)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Articles)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(generateErrorResponse("Invalid ID"))
		return
	}

	for key, val := range Articles {
		if val.ID == ID {
			Articles = append(Articles[:key], Articles[key+1:]...)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Articles)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(generateErrorResponse("No record deleted"))

}

func handleRequest() {

	srvRouter := mux.NewRouter().StrictSlash(true)
	srvRouter.HandleFunc("/", homePage)
	srvRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")
	srvRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	srvRouter.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
	srvRouter.HandleFunc("/article", createNewArticle).Methods("POST")

	log.Fatal(http.ListenAndServe(":10000", srvRouter))
}

func main() {
	fmt.Println("REST API v5.0 - New Delete Article Handler")

	Articles = []Article{
		Article{ID: 1, Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		{ID: 2, Title: "Hello 2", Desc: "Article 2 Description", Content: "Article 2 Content"},
	}
	handleRequest()
}
