package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

const maxUploadSize = 50 << 20

func FileUploadPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/upload.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}

func UploadFiles(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in UploadFiles: %v", r)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	// Limit the size of the request to prevent out of memory issues
	if err := r.ParseMultipartForm(maxUploadSize); err != nil { // Limit to 10 MB files.
		if err == http.ErrNotMultipart {
			http.Error(w, "Request body must be multipart/form-data", http.StatusBadRequest)
			return
		}
		http.Error(w, "The uploaded file is too large. Please upload a file less than 50MB.", http.StatusBadRequest)
		return
	}

	// Check if the parsed form is nil
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		http.Error(w, "No file data", http.StatusBadRequest)
		return
	}

	files, ok := r.MultipartForm.File["files"]
	if !ok {
		http.Error(w, "No files found in form", http.StatusBadRequest)
		return
	}

	uploadDir := "./uploads/"
	var pdfPaths []string

	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("Error creating directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Process each file
	for i, fileHeader := range files {
		if fileHeader.Size > maxUploadSize {
			http.Error(w, fmt.Sprintf("The uploaded file '%s' exceeds the size limit of %d MB.", fileHeader.Filename, maxUploadSize/(1<<20)), http.StatusBadRequest)
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			continue
		}
		defer file.Close()

		timestamp := time.Now().Format("20060102-150405")
		filePath := filepath.Join(uploadDir, fmt.Sprintf("%s_%d_%s", timestamp, i, fileHeader.Filename))
		newFile, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating file: %v", err)
			continue
		}
		defer newFile.Close()

		if _, err := io.Copy(newFile, file); err != nil {
			log.Printf("Error saving file: %v", err)
			continue
		}

		pdfPaths = append(pdfPaths, filePath)
	}
	if len(pdfPaths) == 0 {
		log.Printf("No PDFs to merge")
		http.Error(w, "No PDFs were uploaded successfully, cannot merge", http.StatusInternalServerError)
		return
	}

	// Now let's merge the uploaded PDFs
	mergedFileName := "merged_" + time.Now().Format("20060102150405") + ".pdf"
	outputPath := filepath.Join(uploadDir, mergedFileName)
	if err := api.MergeCreateFile(pdfPaths, outputPath, nil); err != nil {
		log.Printf("Error merging PDF files: %v", err)
		http.Error(w, "Error merging PDF files", http.StatusInternalServerError)
		return
	}

	downloadLink := "/download/" + mergedFileName

	// Respond with a page that includes the download link
	tmpl, err := template.ParseFiles("templates/download.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, map[string]string{"DownloadLink": downloadLink, "UploadPath": "/upload"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	defer func() {
		// Delete the uploaded files
		for _, path := range pdfPaths {
			if err := os.Remove(path); err != nil {
				log.Printf("Failed to delete uploaded file '%s': %v", path, err)
			}
		}
	}()

}
