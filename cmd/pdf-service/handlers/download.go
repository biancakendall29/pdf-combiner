package handlers

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func DownloadLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// Specify the directory where the merged files are
	filePath := "./uploads/" + filename

	// Check if file exists and is not a directory before serving
	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/pdf")
		http.ServeFile(w, r, filePath)
	} else {
		http.Error(w, "File not found.", http.StatusNotFound)
	}

	// Optionally, delete the file after serving
	defer func() {
		os.Remove(filePath)
	}()
}
