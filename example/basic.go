/* Example application
 * Demonstration Basic usage of surfer
 * It runs a server which handle / path and stores the number of visits in an context view
 * run by 'go run function_pattern.go' and launch in browser: http://127.0.0.1:8000/
 */

package main

import (
	"github.com/scale-it/surfer"
)

type Counter struct {
	surfer.Handler
	num int
}

func (this *Counter) Get() {
	this.App.Log.Debug("Handled index")
	this.num++
	//fmt.Fprintf(this.Response, "Hello World, counter: %d", this.num)
}

func main() {
	app := surfer.New()
	app.Router.Handle("/", &Counter{}) // Handle a simple function
	app.Run()
}
