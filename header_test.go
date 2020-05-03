package hiox

import (
	"net/http"
	"testing"

	"github.com/frk/compare"
)

func TestCookieValues_ReadHeader(t *testing.T) {
	Sp := func(s string) *string { return &s }

	tests := []struct {
		h    http.Header
		dest CookieValues
		want CookieValues
	}{{
		h:    http.Header{"Cookie": {"Cookie-1=v$1; c2=v2"}},
		dest: CookieValues{"Cookie-1": new(string), "c2": new(string)},
		want: CookieValues{"Cookie-1": Sp("v$1"), "c2": Sp("v2")},
	}, {
		h:    http.Header{"Cookie": {"Cookie-1=v$1; c2=v2"}},
		dest: CookieValues{"Cookie-1": new(string)},
		want: CookieValues{"Cookie-1": Sp("v$1")},
	}, {
		h:    http.Header{"Cookie": {"Cookie-1=v$1; c2=v2"}},
		dest: CookieValues{"c2": new(string)},
		want: CookieValues{"c2": Sp("v2")},
	}}

	for _, tt := range tests {
		_ = tt.dest.ReadHeader(tt.h)
		if e := compare.Compare(tt.dest, tt.want); e != nil {
			t.Error(e)
		}
	}
}
