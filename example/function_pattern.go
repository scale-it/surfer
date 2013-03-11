/* Example application
 * Demonstration of Function Router pattern
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

func Counter(app surfer.App, ctx *Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Log.Debug("Handled index")
		ctx.num++
		fmt.Fprintf(w, "Hello World, counter: %d", ctx.num)
	}
}

type Context struct {
	num int
	db  map[string]int
}

func main() {
	app := surfer.New()
	ctx := new(Context)
	app.Router.HandleFunc("/", Home)                   // Handle a simple function
	app.Router.HandleFunc("/index", Counter(app, ctx)) // Handle function which needs a context
	app.Run()
}
