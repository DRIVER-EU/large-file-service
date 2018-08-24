package main

import (
	"path/filepath"
	"time"
	//"encoding/json"
	"crypto/rand"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	log.Printf("File uploader service started!")

	// handles upload requests
	router.HandleFunc("/public/upload", upload).Methods("GET", "POST")

	// handles upload requests that should be private
	router.HandleFunc("/private/upload", privateUpload).Methods("GET", "POST")

	// handles private file retrieval by ID
	router.HandleFunc("/private/{fileId}", getFile).Methods("GET")

	// public file server
	fs := http.FileServer(http.Dir("./public/"))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	server := &http.Server{
		Handler: router,
		Addr:    ":9090",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func getFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileId"]

	path := "./private/" + fileID
	log.Println(path)
	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		http.ServeFile(w, r, path)
		return
	}

	http.NotFound(w, r)
}

func privateUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, "private")
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println("Unable to open Form file: ", err)
			return
		}
		defer file.Close()

		uuid := generateUUID()
		fileExt := filepath.Ext(handler.Filename)
		targetFileName := uuid + fileExt
		f, err := os.OpenFile("./private/"+targetFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("Could not create private file: ", err)
			return
		}
		defer f.Close()

		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, targetFileName)
	}
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, "public")
	} else {

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./public/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)
	}
}

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
