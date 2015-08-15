# snippet-store [![Build Status](https://travis-ci.org/gopherjs/snippet-store.svg?branch=master)](https://travis-ci.org/gopherjs/snippet-store) [![GoDoc](https://godoc.org/github.com/gopherjs/snippet-store?status.svg)](https://godoc.org/github.com/gopherjs/snippet-store)

snippet-store is a server for storing GopherJS Playground snippets.

It uses the same mapping scheme as the Go Playground, and when a snippet isn't found locally,
it defers to fetching it from the Go Playground. This effectively augments our world of available
snippets with that the Go Playground.

Newly shared snippets are stored locally in the specified folder and take precedence.

Installation
------------

```bash
go get -u github.com/gopherjs/snippet-store
```

License
-------

- [MIT License](http://opensource.org/licenses/mit-license.php)
