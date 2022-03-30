package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/upload", uploadFileHandler)
	mux.HandleFunc("/serve/", serveFileHandler)

	// Prepare upload directory
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating the uploads directory:", err.Error())
		return
	}

	// Start Webserver
	if err := http.ListenAndServe(":8090", mux); err != nil {
		fmt.Println("Error starting the webserver:", err.Error())
		return
	}
}

func pingHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Pong")
}

func uploadFileHandler(w http.ResponseWriter, req *http.Request) {
	file, fileHandler, err := req.FormFile("file")
	if err != nil {
		http.Error(w, "Error when parsing the uploaded file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a new file in the uploads directory
	destFile, err := os.Create(fmt.Sprintf("./uploads/%s", fileHandler.Filename))
	if err != nil {
		http.Error(w, "Error when trying to write the uploaded file to the upload directory:"+err.Error(), http.StatusInternalServerError)
		return
	}

	defer destFile.Close()

	// Write uploaded file to filesystem
	_, err = io.Copy(destFile, file)
	if err != nil {
		http.Error(w, "Error when trying to write the uploaded file to the upload directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File successfully uploaded!")
}

func serveFileHandler(w http.ResponseWriter, req *http.Request) {
	requestedFile := strings.TrimPrefix(req.URL.Path, "/serve/")
	http.ServeFile(w, req, "uploads/"+requestedFile)
}
