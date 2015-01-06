package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// getSnippetFromLocalStore tries to get the snippet with given id from local store.
// If it returns nil error, the ReadCloser must be closed by caller.
func getSnippetFromLocalStore(id string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(*storageDirFlag, id))
}

const userAgent = "gopherjs.org/play/ playground snippet fetcher"

// getSnippetFromGoPlayground tries to get the snippet with given id from the Go Playground.
// If it returns nil error, the ReadCloser must be closed by caller.
func getSnippetFromGoPlayground(id string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", "https://play.golang.org/p/"+id+".go", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Go Playground returned unexpected status code %v", resp.StatusCode)
	}

	return resp.Body, nil
}

// Store snippet in local storage.
func storeSnippet(body []byte) (id string, err error) {
	id = snippetBodyToId(body)
	err = ioutil.WriteFile(filepath.Join(*storageDirFlag, id), body, 0644)
	return id, err
}
