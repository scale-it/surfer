package surfer

import (
	"encoding/json"
	"github.com/gorilla/securecookie"
	"mime"
	"net/http"
	"strconv"
	"strings"
)


type statet int8

var states = struct {
	rendered statet
	finished statet
}{1, 2}

var valid_methods = map[string]bool{"GET": true, "POST": true, "HEAD": true, "DELETE": true, "PUT": true, "PATCH": true, "OPTIONS": true}

type Handler struct {
	App *App
	//Template *Template
	template  string // template name to render, when accept is text/html
	Response  http.ResponseWriter
	Request   *http.Request
	Data      map[interface{}]interface{} // data for template
	state     statet
	SecCookie *securecookie.SecureCookie
}

type SurferHandler interface {
	http.Handler
	Get()
	Post()
	Head()
	Delete()
	Put()
	Patch()
	Options()
	Prepare() bool
	Finish()
	Render()
}

// Implement any of the following methods to handle the corresponding HTTP method.

func (this *Handler) Get() {
	http.Error(this.Response, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (this *Handler) Post() {
	http.Error(this.Response, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (this *Handler) Head() {
	http.Error(this.Response, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (this *Handler) Delete() {
	http.Error(this.Response, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (this *Handler) Put() {
	http.Error(this.Response, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (this *Handler) Patch() {
	http.Error(this.Response, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (this *Handler) Options() {
	http.Error(this.Response, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// Implementation of html.Handler interface
func (this *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	this.Response = rw
	this.Request = req

	if !valid_methods[req.Method] {
		http.Error(this.Response, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if this.Prepare() {
		switch req.Method {
		case "GET":
			this.Get()
		case "POST":
			this.Post()
		case "HEAD":
			this.Head()
		case "DELETE":
			this.Delete()
		case "PUT":
			this.Put()
		case "PATCH":
			this.Patch()
		case "OPTIONS":
			this.Options()
		default:
			panic("Implementation Error. This shouldn't be accessible")
		}
	}

	if this.state < states.rendered {
		this.Render()
	}
	if this.state < states.finished {
		if !this.Finish() {
			// If finish fails we need to ensure that we don't leave anything hanging
			this.App.Log.Error("%s, Handler.Finish returned false.", this.Request.URL.Path)
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

// This is most useful in a your base Handler. Override this method to perform common initialization regardless of the request method.
// Prepare is called no matter which HTTP method is used. prepare may produce output. If it calls Finish, processing stops here.
func (this *Handler) Prepare() bool {
	return true
}

// Override this method to perform cleanup, logging, etc.
// This method is a counterpart to prepare. Finish may not produce any output, as it is called after the response has been sent to the client.
func (this *Handler) Finish() bool {
	this.state = states.finished
	return true
}

// You probably won't overwrite this method. Better thing about Finish method.
func (this *Handler) Render() {
	// TODO check if html.Error wasn't called
	this.state = states.rendered
	if this.Data == nil {
		this.Response.WriteHeader(http.StatusNoContent)
		return
	}

	// Find the best media-type by looking at the first item of
	// intersection of the Accept and Supported sets.
	mt := selectMediaType(this.Request.Header.Get("Accept"))

	// Set the content type based on the accept header
	this.SetContentType(mt)

	// render based on template (filename, "json"...) and this.Data
	switch mt {
	case "application/json":
		content, err := json.Marshal(this.Data)
		if err != nil {
			this.App.Log.Error("%s, %v", this.Request.URL.Path, err.Error())
			http.Error(this.Response, err.Error(), http.StatusInternalServerError)
			return
		}
		this.Response.Header().Set("Content-Length", strconv.Itoa(len(content)))
		this.Response.Write(content)
		this.SetContentType(mt)
		// TODO: check if it is OK

	case "text/html":
		if this.template == "" {
			this.App.Log.Fatal("%s, Trying to render template, but it was not set", this.Request.URL.Path)
			// TODO: fall back to json? or crash in production?
		}
		// TODO: render template
	default:
		// We need to look at what's a resonable default
		// output, probably HTML.
	}
}

// HELPER FUNCTIONS ------------------------------------------

// Redirect to different `url` with `status`
func (this *Handler) Redirect(status int, url string) {
	this.Response.Header().Set("Location", url)
	this.Response.WriteHeader(status)
}

// Redirect to different `url` with Found HTTP status code: 302
func (this *Handler) RedirectUrl(url string) {
	this.Redirect(http.StatusFound, url)
}

//Sets the content type by extension, as defined in the mime package.
//For example, xgoContext.ContentType("json") sets the content-type to "application/json"
func (this *Handler) SetContentType(ctype string) {
	if !strings.HasPrefix(ctype, ".") {
		ctype = "." + ctype
	}
	ctype = mime.TypeByExtension(ctype)
	if ctype != "" {
		this.Response.Header().Set("Content-Type", ctype)
	}
}

// Gets the quality of media-type of the form foo/bar;q=XYZ
func getMtQual(mediatype string) (string, float64, error) {
	tp, par, err := mime.ParseMediaType(mediatype)
	if err != nil {
		return "", 0, err
	}
	q, err := strconv.ParseFloat(par["q"], 32)
	if err != nil {
		q = 1
	}
	// we can't get an error here: omitting the quality means q=1
	return tp, q, nil
}

// prelimimary implementation of mediatype parsing, this _should_ work
func selectMediaType(accept string) string {
	// TODO: this shouldn't be hard coded
	var supported = [3]string{"application/json", "text/html", "application/xml"}
	accepted := strings.Split(accept, ",")
	bestmime := ""
	bestqual := 1.0
	for _, acc := range accepted {
		tp, q, err := getMtQual(acc)
		if err != nil {
			continue
		}

		if bestmime == "" || q > bestqual {
			for _, mt := range supported {
				if mt == acc {
					bestmime = tp
					bestqual = q
					break
				}
			}
		}
	}
	if bestmime != "" {
		return bestmime
	}
	// We should consider returning nil here
	return supported[0]
}
