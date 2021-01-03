package main

import (
	"io/ioutil"
	"github.com/DataDog/zstd"
	"log"
	"net/http"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	compressed, err := StoreGet()
	if err != nil {
		log.Fatal(err)
		return
	}
	if compressed == nil {
		http.NotFound(w, r)
		return
	}

	decompressed, err := zstd.Decompress(nil, compressed)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Decompressing %d -> %d bytes", len(compressed), len(decompressed))

	w.Write(decompressed)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	compressed, err := zstd.Compress(nil, body)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Compressing %d -> %d bytes", len(body), len(compressed))

	err = StoreSet(compressed)
	if err != nil {
		log.Fatal(err)
		return
	}

	//io.Write(w, "OK")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
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

	err := StoreConstruct()
	if err != nil {
		log.Fatal(err)
	}

	defer StoreDestroy()

	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
