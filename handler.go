package surfer

import (
	"mime"
	"net/http"
)

type statet int8

var states = struct {
	rendered statet
	finished statet
}{1, 2}

var valid_methods = map[string]bool{"GET": true, "POST": true, "HEAD": true, "DELETE": true, "PUT": true, "PATCH": true, "OPTIONS": true}

type Handler struct {
	app *App
	//Template *Template
	template string
	Response http.ResponseWriter
	Request  *http.Request
	Data     map[interface{}]interface{} // data for template
	Params   map[string]string
	state    statet
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

func (this *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	this.RW = rw
	this.Request = req

	if !valid_methods[req.Method] {
		return http.Error(c.Ctx.Response, "Method Not Allowed", 405)
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
		return http.Error(this.Response, "Method Not Allowed", 405)
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
	return true
}

func (this *Handler) Render() {
	// validate template against header accept
	// render based on template (filename, "json"...) and this.Data
	// set state
}

func (this *Handler) Get() {
	http.Error(this.Response, "Method Not Allowed", 405)
}

func (this *Handler) Post() {
	http.Error(this.Response, "Method Not Allowed", 405)
}

func (this *Handler) Head() {
	http.Error(this.Response, "Method Not Allowed", 405)
}

func (this *Handler) Delete() {
	http.Error(this.Response, "Method Not Allowed", 405)
}

func (this *Handler) Put() {
	http.Error(this.Response, "Method Not Allowed", 405)
}

func (this *Handler) Patch() {
	http.Error(this.Response, "Method Not Allowed", 405)
}

func (this *Handler) Options() {
	http.Error(this.Response, "Method Not Allowed", 405)
}

// HELPER FUNCTIONS

// Redirect to different `url` with `status`
func (this *Handler) Redirect(status int, url string) {
	this.Response.Header().Set("Location", url_)
	ctx.ResponseWriter.WriteHeader(status)
}

// Redirect to different `url` with standard HTTP status code: 404
func (this *Handler) RedirectUrl(url string) {
	this.Redirect(302, url)
}

//Sets the content type by extension, as defined in the mime package.
//For example, xgoContext.ContentType("json") sets the content-type to "application/json"
func (this *Handler) SetContentType(ctype string) {
	if !strings.HasPrefix(ctype, ".") {
		ctype = "." + ext
	}
	ctype := mime.TypeByExtension(ctype)
	if ctype != "" {
		this.Response.Header().Set("Content-Type", ctype)
	}
}
