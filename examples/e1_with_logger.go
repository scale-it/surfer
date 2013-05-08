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

// open http://0.0.0.0:8000/ and http://0.0.0.0:8000/log and http://0.0.0.0:8000/log/other
// to see efects
package main

import (
	"fmt"
	"github.com/fiorix/go-web/httpxtra"
	"github.com/scale-it/go-log"
	"github.com/scale-it/surfer"
	"net/http"
	"os"
	"time"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	counter += 1
	fmt.Fprintln(w, "Hello, world", counter)
}

var logger *log.Logger = log.New()
var counter int

func main() {
	logger.AddHandler(os.Stderr, log.Levels.Trace, log.TimeFormatter{"Request"})

	// here we use httpxtra, which preserve status code and support logger function.
	http.Handle("/log", httpxtra.Handler{
		Logger: func(req *http.Request, created time.Time, status, bytes int) {
			surfer.LogRequest(logger, req, created, status, bytes)
		},
		Handler: http.HandlerFunc(IndexHandler)})
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8000", nil)
}
