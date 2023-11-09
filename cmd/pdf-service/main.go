package main

import (
	"log"
	"net/http"

	"github.com/biancakendall29/pdf-combiner/cmd/pdf-service/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/upload", handlers.FileUploadPage).Methods("GET")
	r.HandleFunc("/upload", handlers.UploadFiles).Methods("POST")
	r.HandleFunc("/download/{filename}", handlers.DownloadLink)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Start the server
	log.Println("Listing for requests at http://localhost:8000/upload")
	log.Fatal(http.ListenAndServe(":8000", r))
}
