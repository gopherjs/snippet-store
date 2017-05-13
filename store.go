package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/shurcooL/webdavfs/vfsutil"
	"golang.org/x/net/webdav"
)

// localStore is the local store for snippets. Snippets are kept as files in root directory.
var localStore webdav.FileSystem

// getSnippetFromLocalStore tries to get the snippet with given id from local store.
// If it returns nil error, the ReadCloser must be closed by caller.
func getSnippetFromLocalStore(ctx context.Context, id string) (io.ReadCloser, error) {
	return vfsutil.Open(ctx, localStore, id)
}

const userAgent = "gopherjs.org/play/ playground snippet fetcher"

// getSnippetFromGoPlayground tries to get the snippet with given id from the Go Playground.
// If it returns nil error, the ReadCloser must be closed by caller.
func getSnippetFromGoPlayground(ctx context.Context, id string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", "https://play.golang.org/p/"+id+".go", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Go Playground returned unexpected status code %v", resp.StatusCode)
	}

	return resp.Body, nil
}

// storeSnippet stores snippet in local storage.
// It returns the id assigned to the snippet.
func storeSnippet(ctx context.Context, body []byte) (id string, err error) {
	id = snippetBodyToID(body)
	err = vfsutil.WriteFile(ctx, localStore, id, body, 0644)
	return id, err
}
