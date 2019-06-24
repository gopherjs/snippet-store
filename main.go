// snippet-store is a server for storing GopherJS Playground snippets.
//
// It uses the same mapping scheme as the Go Playground, and when a snippet isn't found locally,
// it defers to fetching it from the Go Playground. This effectively augments our world of available
// snippets with that of the Go Playground.
//
// Newly shared snippets are stored locally in the specified folder and take precedence.
package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/webdav"
)

var (
	storageDirFlag  = flag.String("storage-dir", "", "Storage dir for snippets; if empty, a volatile in-memory store is used.")
	httpFlag        = flag.String("http", ":8080", "Listen for HTTP connections on this address.")
	allowOriginFlag = flag.String("allow-origin", "http://www.gopherjs.org", "Access-Control-Allow-Origin header value.")
)

const maxSnippetSizeBytes = 1024 * 1024

func main() {
	flag.Parse()

	var localStore webdav.FileSystem
	switch *storageDirFlag {
	default:
		if fi, err := os.Stat(*storageDirFlag); os.IsNotExist(err) {
			log.Fatalf("storage directory %q doesn't exist: %v", *storageDirFlag, err)
		} else if err != nil {
			log.Fatalf("error doing stat of directory %q: %v", *storageDirFlag, err)
		} else if !fi.IsDir() {
			log.Fatalf("file %q is not a directory", *storageDirFlag)
		}
		localStore = webdav.Dir(*storageDirFlag)
	case "":
		localStore = webdav.NewMemFS()
	}

	s := &Server{
		Store: &Store{
			LocalFS: localStore,
		},
	}
	http.HandleFunc("/share", s.ShareHandler) // "/share", save snippet and return its id.
	http.HandleFunc("/p/", s.PHandler)        // "/p/{{.SnippetId}}", serve snippet by id.

	log.Println("Started.")

	err := http.ListenAndServe(*httpFlag, nil)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}

// Server is the snippet store HTTP server.
type Server struct {
	Store *Store
}

// ShareHandler is the handler for "/share" requests.
func (s *Server) ShareHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", *allowOriginFlag)
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Needed for Safari.

	if req.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method should be POST.", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(http.MaxBytesReader(w, req.Body, maxSnippetSizeBytes))
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
		return
	}

	id, err := s.Store.StoreSnippet(req.Context(), body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
		return
	}

	_, err = io.WriteString(w, id)
	if err != nil {
		log.Println(err)
		return
	}
}

// PHandler is the handler for "/p/{{.SnippetId}}" requests.
func (s *Server) PHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", *allowOriginFlag)

	if req.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method should be GET.", http.StatusMethodNotAllowed)
		return
	}

	id := req.URL.Path[len("/p/"):]
	err := validateID(id)
	if err != nil {
		http.Error(w, "Unexpected id format.", http.StatusBadRequest)
		return
	}

	// Set a 3 minute timeout to load and serve the snippet.
	ctx, cancel := context.WithTimeout(req.Context(), 3*time.Minute)
	defer cancel()

	snippet, err := s.Store.LoadSnippet(ctx, id)
	if os.IsNotExist(err) {
		// Snippet not found.
		http.Error(w, "Snippet not found.", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
		return
	}
	defer snippet.Close()

	_, err = io.Copy(w, snippet)
	if err != nil {
		log.Println(err)
		return
	}
}
