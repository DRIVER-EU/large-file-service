package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type uploadResponse struct {
	FileURL string
}

var publicDir string
var privateDir string
var hostname string

func main() {
	// gets settings or sets defaults
	hostname = getEnv("HOST", "localhost:9090")
	publicDir = strings.TrimSuffix(getEnv("PUBLIC_UPLOAD_DIR", "public"), "/")
	privateDir = strings.TrimSuffix(getEnv("PRIVATE_UPLOAD_DIR", "private"), "/")
	writeTimeout, _ := strconv.Atoi(getEnv("WRITE_TIMEOUT_SECS", "120"))
	readTimeout, _ := strconv.Atoi(getEnv("READ_TIMEOUT_SECS", "120"))

	// create a server router for relative paths
	router := mux.NewRouter()

	// handles public and private upload requests
	router.HandleFunc("/upload", upload).Methods("GET", "POST")

	// handles private file retrieval by ID
	router.HandleFunc("/private/{fileId}", getFile).Methods("GET")

	// public file server
	fs := http.FileServer(http.Dir(publicDir))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	server := &http.Server{
		Handler:      router,
		Addr:         hostname,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}

	makeDirs()

	log.Printf("File uploader service started! Listening at: http://" + hostname)
	log.Fatal(server.ListenAndServe())
}

func getFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileId"]

	path := privateDir + "/" + fileID
	log.Println(path)
	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		http.ServeFile(w, r, path)
		return
	}

	http.NotFound(w, r)
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseMultipartForm(32 << 20)

		// get the uploaded file from the form values
		file, handler, err := r.FormFile("uploadFile")
		if err != nil {
			log.Println("Unable to open Form file: ", err)
			return
		}
		defer file.Close()

		private := r.FormValue("private")
		var targetFileName string
		if len(private) != 0 {
			uuid := generateUUID()
			fileExt := filepath.Ext(handler.Filename)
			targetFileName = privateDir + "/" + uuid + fileExt
		} else {
			targetFileName = publicDir + "/" + handler.Filename
		}

		f, err := os.OpenFile(targetFileName, os.O_WRONLY|os.O_CREATE, 0744)
		if err != nil {
			log.Println("Could not create file: ", err)
			return
		}
		defer f.Close()

		io.Copy(f, file)

		response := uploadResponse{getHostURL() + "/" + targetFileName}
		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

// Gets an Environment variable with a default as fallback
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func makeDirs() {
	os.MkdirAll(publicDir, 0744)
	os.MkdirAll(privateDir, 0744)
}

func getHostURL() string {
	return "http://" + hostname
}

// Generates a UUID that can be used as an obfuscated file name
func generateUUID() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)

	if err != nil {
		log.Fatalln("Could not generate UUID for private file upload!", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
