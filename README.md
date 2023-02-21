snippet-store
=============

[![Go Reference](https://pkg.go.dev/badge/github.com/gopherjs/snippet-store.svg)](https://pkg.go.dev/github.com/gopherjs/snippet-store)

snippet-store is a server for storing GopherJS Playground snippets.

It uses the same mapping scheme as the Go Playground, and when a snippet isn't found locally,
it defers to fetching it from the Go Playground. This effectively augments our world of available
snippets with that of the Go Playground.

Newly shared snippets are stored locally in the specified folder and take precedence.

Installation
------------

```bash
go install github.com/gopherjs/snippet-store@latest
```

License
-------

-	[MIT License](https://opensource.org/licenses/mit-license.php)
