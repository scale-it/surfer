/* Example application
 * Demonstration of Function Context pattern
 * It runs a server which handle /index path and stores the number of visits in a context view
 * run by 'go run function_pattern.go' and launch in browser: http://127.0.0.1:8000/index
 */

package main

import (
	"fmt"
	"github.com/scale-it/surfer"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html><body>
  check <a href="/index">/index</a> to handle some context
</body></html>`))
}

type Context struct {
	num int
	db  map[string]int
}

func counter(app surfer.App, ctx *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Log.Debug("Handled index")
		ctx.num++
		fmt.Fprintf(w, "Hello World, counter: %d", ctx.num)
	}
}

// The idea of above function we can put in more general form

type ContextFunc func(surfer.App, *Context, http.ResponseWriter, *http.Request)

func WithContext(app surfer.App, ctx *Context, f ContextFunc) http.HandlerFunc {
	// init code

	return func(w http.ResponseWriter, r *http.Request) {
		// prepare code
		app.Log.Debug("Preparing")

		f(app, ctx, w, r)

		// finish code
		app.Log.Debug("Finalize")
	}
}

func counter2(app surfer.App, ctx *Context, w http.ResponseWriter, r *http.Request) {
	app.Log.Debug("Handled index")
	ctx.num++
	fmt.Fprintf(w, "Hello World, counter: %d", ctx.num)
}

func main() {
	app := surfer.New()
	ctx := new(Context)
	app.Router.HandleFunc("/", Home)                                  // Handle a simple function
	app.Router.HandleFunc("/index", counter(app, ctx))                // Handle function which needs a context
	app.Router.HandleFunc("/index2", WithContext(app, ctx, counter2)) // Handle function which needs a context, with common init / prepare / finish
	app.Run()
}
