// Command snippet-store is a server for storing GopherJS Playground snippets.
//
// It uses the same mapping scheme as the Go Playground, and when a snippet isn't found locally,
// it defers to fetching it from the Go Playground. This effectively augments our world of available
// snippets with that the Go Playground.
//
// Newly shared snippets are stored locally in the specified folder and take precedence.
//
package main

import (
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var storageDirFlag = flag.String("storage-dir", filepath.Join(os.TempDir(), "gopherjs_snippets"), "Storage dir for snippets.")
var httpFlag = flag.String("http", ":8080", "Listen for HTTP connections on this address.")

const allowOrigin = "http://gopherjs.org"

func pHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method should be GET.", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", allowOrigin)

	id := req.URL.Path[len("/p/"):]

	// TODO: Newly shared snippets are stored locally in the specified folder and take precedence.

	req2, err := http.NewRequest("GET", "http://play.golang.org/p/"+id+".go", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
		return
	}
	req2.Header.Set("User-Agent", "gopherjs.org/play/ playground snippet fetcher")

	resp, err := http.DefaultClient.Do(req2)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	default:
		log.Println("Unexpected StatusCode from Go Playground:", resp.StatusCode)
		http.Error(w, "Server error.", http.StatusInternalServerError)
		return
	case http.StatusOK:
		// Snippet found on Go Playground.
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server error.", http.StatusInternalServerError)
			return
		}
	case http.StatusNotFound:
		// Snippet not found on Go Playground.
		file, err := os.Open(filepath.Join(*storageDirFlag, id))
		if err != nil {
			http.Error(w, "Snippet not found.", http.StatusNotFound)
			return
		}
		defer file.Close()

		_, err = io.Copy(w, file)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server error.", http.StatusInternalServerError)
			return
		}
	}
}

func shareHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Forbidden.", http.StatusForbidden)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", allowOrigin)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
		return
	}

	id := snippetBodyToId(body)

	err = ioutil.WriteFile(filepath.Join(*storageDirFlag, id), body, 0644)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
		return
	}

	_, err = io.WriteString(w, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
		return
	}
}

func main() {
	flag.Parse()

	err := os.MkdirAll(*storageDirFlag, 0755)
	if err != nil {
		log.Fatalf("Error creating directory %q: %v.\n", *storageDirFlag, err)
	}

	http.HandleFunc("/p/", pHandler)        // "/p/{{.SnippetId}}", serve snippet by id.
	http.HandleFunc("/share", shareHandler) // "/share", save snippet and return its id.

	panic(http.ListenAndServe(*httpFlag, nil))
}

// snippetBodyToId mimics the mapping scheme used by the Go Playground.
func snippetBodyToId(body []byte) string {
	const salt = "[replace this with something unique]"

	h := sha1.New()
	io.WriteString(h, salt)
	h.Write(body)
	sum := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)[:10]
}
