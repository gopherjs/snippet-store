package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
)

// snippetBodyToID mimics the mapping scheme used by the Go Playground.
func snippetBodyToID(body []byte) string {
	// This was the actual salt value used by Go Playground at some time in the past.
	// It came from https://code.google.com/p/go-playground/source/browse/goplay/share.go#18.
	// See https://github.com/gopherjs/snippet-store/pull/1#discussion_r22512198 for more details.
	//
	// It has since changed. We continue to use the same value for now to keep things consistent.
	// See https://github.com/golang/playground/blob/a72214bb7a8781349b57129256dc0c64d233ef08/share.go#L17-L19
	// for the current status. It's possible to change our hashing algorithm to be in sync with
	// the Go Playground.
	const salt = "[replace this with something unique]"

	h := sha1.New()
	io.WriteString(h, salt)
	h.Write(body)
	sum := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)[:10]
}

// validateID returns an error if id is of unexpected format.
// ID of length 10 and 11 are both supported, so that historical IDs continue to work.
func validateID(id string) error {
	if len(id) != 10 && len(id) != 11 {
		return fmt.Errorf("id length is %v instead of 10 or 11", len(id))
	}

	for _, b := range []byte(id) {
		ok := ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z') || ('0' <= b && b <= '9') || b == '-' || b == '_'
		if !ok {
			return fmt.Errorf("id contains unexpected character %+q", b)
		}
	}

	return nil
}
