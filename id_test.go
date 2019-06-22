package main

import (
	"fmt"
	"testing"
)

func TestValidateID(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want error
	}{
		{
			in:   "D9L6MbPfE4",
			want: nil,
		},
		{
			in:   "ABZdez09-_",
			want: nil,
		},
		{
			in:   "N_M_YelfGeR",
			want: nil,
		},
		{
			in:   "Abc",
			want: fmt.Errorf("id length is 3 instead of 10 or 11"),
		},
		{
			in:   "Abc?q=1235",
			want: fmt.Errorf("id contains unexpected character '?'"),
		},
		{
			in:   "../../file",
			want: fmt.Errorf("id contains unexpected character '.'"),
		},
		{
			in:   "Heya世界",
			want: fmt.Errorf(`id contains unexpected character '\u00e4'`),
		},
	} {
		got := validateID(tc.in)
		if !equalError(got, tc.want) {
			t.Errorf("validateID(%q) error doesn't match:\n got: %v\nwant: %v", tc.in, got, tc.want)
		}
	}
}

// equalError reports whether errors a and b are considered equal.
// They're equal if both are nil, or both are not nil and a.Error() == b.Error().
func equalError(a, b error) bool {
	return a == nil && b == nil || a != nil && b != nil && a.Error() == b.Error()
}
