package surfer

import (
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

func (this *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	this.Response = rw
	this.Request = req

	if !valid_methods[req.Method] {
		http.Error(this.Response, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if !this.Prepare() {
		return
	}
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
		panic("Implementation Error. This shouldn't be accessiable")
	}

	if this.state < states.rendered {
		this.Render()
	}
	if this.state < states.finished {
		this.Finish()
	}

}

func (this *Handler) Prepare() bool {
	return true
}

func (this *Handler) Finish() bool {
	this.state = states.finished
	return true
}

func (this *Handler) Render() {
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
		// TODO: send out json
	case "text/html":
		if this.template == "" {
			// TODO: fall back to json? or crash in production?
		}
		// TODO: render template
	default:
		// We need to look at what's a resonable default
		// output, probably HTML.
	}
}

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

// HELPER FUNCTIONS

// Redirect to different `url` with `status`
func (this *Handler) Redirect(status int, url string) {
	this.Response.Header().Set("Location", url)
	this.Response.WriteHeader(status)
}

// Redirect to different `url` with standard HTTP status code: 404
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
