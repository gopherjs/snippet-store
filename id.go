package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
)

// snippetBodyToID mimics the mapping scheme used by the Go Playground.
func snippetBodyToID(body []byte) string {
	// This is the actual salt value used by Go Playground, it comes from
	// https://code.google.com/p/go-playground/source/browse/goplay/share.go#18.
	// See https://github.com/gopherjs/snippet-store/pull/1#discussion_r22512198 for more details.
	const salt = "[replace this with something unique]"

	h := sha1.New()
	io.WriteString(h, salt)
	h.Write(body)
	sum := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)[:10]
}

// validateID returns an error if id is of unexpected format.
func validateID(id string) error {
	if len(id) != 10 {
		return fmt.Errorf("id length is %v instead of 10", len(id))
	}

	for _, b := range []byte(id) {
		ok := ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z') || ('0' <= b && b <= '9') || b == '-' || b == '_'
		if !ok {
			return fmt.Errorf("id contains unexpected character %+q", b)
		}
	}

	return nil
}
