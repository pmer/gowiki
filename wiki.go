package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

var store []byte

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if store == nil {
		http.NotFound(w, r)
		return
	}

	w.Write(store)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    store = data
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	if r.Method == "GET" {
		viewHandler(w, r)
	} else if r.Method == "POST" {
		saveHandler(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	log.Println("Starting :)")

	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
