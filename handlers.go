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
	"time"
)

type WithHTTPLogger struct {
	Writer io.Writer
	Next   http.Handler
}

func (this WithHTTPLogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	username := ""
	if req.URL.User != nil {
		if name := req.URL.User.Username(); name != "" {
			username = name
		}
	}
	start := time.Now()
	this.Next.ServeHTTP(w, req)
	elapsed := float64(time.Since(start)) / float64(time.Millisecond)

	fmt.Fprintf(this.Writer, "%s - %s \"%s %s %s\". Elapsed: %f ms\n",
		strings.Split(req.RemoteAddr, ":")[0],
		username,
		req.Method,
		req.RequestURI,
		req.Proto,
		// status,
		elapsed)
}

// Handler which check if request is from canonical host and uses https. Otherwise will
// redirect to https://<canonicalhost>/rest/of/the/url
type WithForceHTTPS struct {
	CanonicalHost string
	Next          http.Handler
}

func (this WithForceHTTPS) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	is_http := true
	if h, ok := req.Header["X-Forwarded-Proto"]; ok {
		if h[0] == "https" {
			is_http = false
		}
	}
	hostPort := strings.Split(req.Host, ":")
	if is_http || hostPort[0] != this.CanonicalHost {
		hostPort[0] = this.CanonicalHost
		url := "https://" + strings.Join(hostPort, ":") + req.URL.String()
		http.Redirect(w, req, url, http.StatusMovedPermanently)
		return
	}

	this.Next.ServeHTTP(w, req)
}

type Authenticator func(req *http.Request) bool

// Enusres authentication for handlers. Otherwise call fallback.
type WithAuth struct {
	A        Authenticator
	Next     http.Handler
	Fallback http.Handler
}

func (this WithAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if this.A(req) {
		this.Next.ServeHTTP(w, req)
	} else {
		this.Fallback.ServeHTTP(w, req)
	}
}

// Calls the wrapped handler and on panic calls the specified error handler.
// errH can make some logging or just return:
//   http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
func PanicHandler(h http.Handler, errH func(http.ResponseWriter, *http.Request, interface{})) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errH(w, r, err)
			}
		}()
		h.ServeHTTP(w, r)
	}
}
