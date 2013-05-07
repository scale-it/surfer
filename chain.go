package surfer

import "net/http"

// chain multiple handlers in a single call.
func Chain(h ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range h {
			v.ServeHTTP(w, r)
		}
	})
}

// Helper function to handle http.HandlerFunc
func ChainFuncs(h ...http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range h {
			v(w, r)
		}
	})
}
