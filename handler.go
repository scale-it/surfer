package surfer

import (
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
	//Session  *Session
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
	this.Response = rw
	this.Request = req

	if !valid_methods[req.Method] {
		return // TODO 405
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
		// TODO 405 Error, Method Not Allowed
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
	// TODO 405 Error, Method Not Allowed
}

func (this *Handler) Post() {
	// TODO 405 Error, Method Not Allowed
}

func (this *Handler) Head() {
	// TODO 405 Error, Method Not Allowed
}

func (this *Handler) Delete() {
	// TODO 405 Error, Method Not Allowed
}

func (this *Handler) Put() {
	// TODO 405 Error, Method Not Allowed
}

func (this *Handler) Patch() {
	// TODO 405 Error, Method Not Allowed
}

func (this *Handler) Options() {
	// TODO 405 Error, Method Not Allowed
}
