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

package surfer

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type WithHTTPLogger struct {
	Writer io.Writer
	Next   http.Handler
}

func (this WithHTTPLogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	username := "-"
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}

	fmt.Fprintf(this.Writer, "%s - %s \"%s %s %s\"\n",
		strings.Split(req.RemoteAddr, ":")[0],
		username,
		req.Method,
		req.RequestURI,
		req.Proto,
		// status,
	)

	this.Next.ServeHTTP(w, req)
}
