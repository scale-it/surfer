/* Copyright 2013 Robert Zaremba
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Run this, and from console:
//    curl -i -H "Accept: application/msgpack"  http://localhost:8000
//    curl -i -H "Accept: application/json"  http://localhost:8000
//    curl -i -H "Accept: text/html"  http://localhost:8000
//    curl  http://localhost:8000
package main

import (
	"github.com/scale-it/go-log"
	"github.com/scale-it/surfer"
	"html/template"
	"net/http"
	"os"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) (string, interface{}, int) {
	counter += 1
	return "simple.html", counter, 200
}

var Log = log.NewStd(os.Stderr, log.Levels.Debug, log.Ldate|log.Lmicroseconds, true)
var counter int

func main() {
	t := template.Must(template.ParseGlob("./templates/*.html"))
	http.Handle("/", surfer.WithRenderer{t, Log, IndexHandler})
	http.ListenAndServe(":8000", nil)
}
