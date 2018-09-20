package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type uploadResponse struct {
	FileURL string
}

var publicDir string
var privateDir string
var URL string
var port string
var writeTimeout int
var readTimeout int

func main() {
	// first set settings from ENV variables with a default
	getEnvVariables()

	// then override with flags if set
	getFlags()

	// create a server router for relative paths
	router := mux.NewRouter()

	// refer root of the service to the upload form
	router.HandleFunc("/", upload).Methods("GET")

	// handles public and private upload requests and shows the upload form
	router.HandleFunc("/upload", upload).Methods("GET", "POST")

	// handles private file retrieval by ID
	router.HandleFunc("/private/{fileId}", getFile).Methods("GET")

	// public file server
	fs := http.FileServer(http.Dir(publicDir))
	router.PathPrefix("/public").Handler(http.StripPrefix("/public", fs))

	// file server for the Swagger UI static resources
	swaggerServer := http.FileServer(http.Dir("swagger"))
	router.PathPrefix("/api/resources").Handler(http.StripPrefix("/api/resources", swaggerServer))

	// function handler to generate an serve the Swagger UI template so it is dynamic based on provided hostname
	router.HandleFunc("/api", swaggerUI).Methods("GET")

	server := &http.Server{
		Handler:      router,
		Addr:         ":80",
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
	}

	makeDirs()

	log.Printf("File uploader service started! Listening at: " + URL)
	log.Printf("Public upload dir is: " + publicDir)
	log.Printf("Private upload dir is: " + privateDir)
	log.Fatal(server.ListenAndServe())
}

// gets and sets the environment variables if present, otherwise falls back to defaults
func getEnvVariables() {
	URL = strings.TrimSuffix(getEnv("URL", "http://localhost"), "/")
	publicDir = strings.TrimSuffix(getEnv("PUBLIC_UPLOAD_DIR", "public"), "/")
	privateDir = strings.TrimSuffix(getEnv("PRIVATE_UPLOAD_DIR", "private"), "/")
	writeTimeout, _ = strconv.Atoi(getEnv("WRITE_TIMEOUT_SECS", "120"))
	readTimeout, _ = strconv.Atoi(getEnv("READ_TIMEOUT_SECS", "120"))
}

// overwrites the configuration with command line parameters if provided
func getFlags() {
	flag.StringVar(&port, "port", port, "port of the file server")
	flag.StringVar(&URL, "url", URL, "URL of the file server (including http://)")
	flag.StringVar(&publicDir, "publicDir", publicDir, "path of the public upload directory")
	flag.StringVar(&privateDir, "privateDir", privateDir, "path of the private upload directory")
	readTimeoutPtr := flag.Int("writeTimeout", readTimeout, "HTTP Server write timeout in seconds.")
	writeTimeoutPtr := flag.Int("readTimeout", writeTimeout, "HTTP Server read timeout in seconds.")
	flag.Parse()
	writeTimeout = *writeTimeoutPtr
	readTimeout = *readTimeoutPtr
}

// serves a private file if it exists, returns 404 otherwise
func getFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileId"]

	path := privateDir + "/" + fileID
	log.Println("Private file requested: " + path)
	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		http.ServeFile(w, r, path)
		return
	}

	http.NotFound(w, r)
}

// renders the Swagger UI template with the configured hostname
func swaggerUI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("swagger/index.gtpl")
		t.Execute(w, nil)
	}
}

// Main upload functionaly. Serves a simple HTML form for get, and handles public and private file uploads via a form post.
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/upload.gtpl")
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

		// check private checkbox, and generate a filepath for a private or public file
		private := r.FormValue("private")
		var targetFileName string
		if len(private) != 0 && !strings.EqualFold(private, "false") {
			uuid := generateUUID()
			fileExt := filepath.Ext(handler.Filename)
			targetFileName = privateDir + "/" + uuid + fileExt
		} else {
			targetFileName = publicDir + "/" + handler.Filename
		}

		// create the file on the filesystem
		f, err := os.OpenFile(targetFileName, os.O_WRONLY|os.O_CREATE, 0744)
		if err != nil {
			log.Println("Could not create file: ", err)
			return
		}
		defer f.Close()

		// copy file from form to filesystem
		io.Copy(f, file)

		// write the JSON response
		response := uploadResponse{URL + "/" + targetFileName}
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

// creates the public and private dirs that were configured
func makeDirs() {
	os.MkdirAll(publicDir, 0744)
	os.MkdirAll(privateDir, 0744)
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
