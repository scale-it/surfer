package surfer

import (
	"bytes"
	"net/http"
	"testing"
)

func TestChain(t *testing.T) {
	var buf bytes.Buffer

	a := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf.WriteRune('a')
	})
	b := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf.WriteRune('b')
	})
	c := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf.WriteRune('c')
	})
	handler := Chain(a, b, c)
	handler.ServeHTTP(nil, nil)

	if buf.String() != "abc" {
		t.Errorf("expected 'abc', got %s", buf.String())
	}
}

func TestChainFuncs(t *testing.T) {
	var buf bytes.Buffer

	a := func(w http.ResponseWriter, r *http.Request) {
		buf.WriteRune('a')
	}
	b := func(w http.ResponseWriter, r *http.Request) {
		buf.WriteRune('b')
	}
	handler := ChainFuncs(a, b)
	handler.ServeHTTP(nil, nil)

	if buf.String() != "ab" {
		t.Errorf("expected 'ab', got %s", buf.String())
	}
}
