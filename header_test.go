package hxio

import (
	"net/http"
	"testing"

	"github.com/frk/compare"
)

func TestCookieValue_ReadHeader(t *testing.T) {
	Sp := func(s string) *string { return &s }

	tests := []struct {
		h    http.Header
		dest CookieValue
		want CookieValue
	}{{
		h:    http.Header{"Cookie": {"Cookie-1=v$1; c2=v2"}},
		dest: CookieValue{"Cookie-1": new(string), "c2": new(string)},
		want: CookieValue{"Cookie-1": Sp("v$1"), "c2": Sp("v2")},
	}, {
		h:    http.Header{"Cookie": {"Cookie-1=v$1; c2=v2"}},
		dest: CookieValue{"Cookie-1": new(string)},
		want: CookieValue{"Cookie-1": Sp("v$1")},
	}, {
		h:    http.Header{"Cookie": {"Cookie-1=v$1; c2=v2"}},
		dest: CookieValue{"c2": new(string)},
		want: CookieValue{"c2": Sp("v2")},
	}}

	for _, tt := range tests {
		_ = tt.dest.ReadHeader(tt.h)
		if e := compare.Compare(tt.dest, tt.want); e != nil {
			t.Error(e)
		}
	}
}
