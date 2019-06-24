package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/shurcooL/webdavfs/vfsutil"
	"golang.org/x/net/webdav"
)

const userAgent = "gopherjs.org/play/ playground snippet fetcher"

// Store is the snippet store.
type Store struct {
	// LocalFS is the local store for snippets. Snippets are kept as files,
	// named after the snippet id (with no extension), in the root directory.
	LocalFS webdav.FileSystem
}

// StoreSnippet stores the provided snippet,
// and returns the id assigned to the snippet.
func (s *Store) StoreSnippet(ctx context.Context, body []byte) (id string, err error) {
	// Store the snippet locally.
	id = snippetBodyToID(body)
	err = vfsutil.WriteFile(ctx, s.LocalFS, id, body, 0644)
	return id, err
}

// LoadSnippet loads the snippet with given id.
// It first tries the local store, then the Go Playground.
//
// It returns an error that satisfies os.IsNotExist if snippet is not found.
// If it returns nil error, the ReadCloser must be closed by caller.
func (s *Store) LoadSnippet(ctx context.Context, id string) (io.ReadCloser, error) {
	// Check if we have the snippet locally first.
	if snippet, err := vfsutil.Open(ctx, s.LocalFS, id); err == nil {
		return snippet, nil
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("loadSnippetFromLocalStore: %v", err)
	}

	// If not found locally, try the Go Playground.
	if snippet, err := fetchSnippetFromGoPlayground(ctx, id); err == nil {
		return snippet, nil
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("fetchSnippetFromGoPlayground: %v", err)
	}

	// Not found in either place.
	return nil, os.ErrNotExist
}

// fetchSnippetFromGoPlayground fetches the snippet with given id
// from the Go Playground.
//
// It returns an error that satisfies os.IsNotExist if snippet is not found.
// If it returns nil error, the ReadCloser must be closed by caller.
func fetchSnippetFromGoPlayground(ctx context.Context, id string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, "https://play.golang.org/p/"+id+".go", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", userAgent)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, os.ErrNotExist
	} else if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Go Playground returned unexpected status code %v", resp.StatusCode)
	}
	return resp.Body, nil
}
