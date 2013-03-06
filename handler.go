package surfer

import (
	"net/http"
)

type Handler struct {
	app      *App
	Template *Template
	Session  *Session
	Response http.ResponseWriter
	Request  *http.Request
	Data     map[interface{}]interface{} // data for template
	Params   map[string]string
	finished bool
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
	Prepare()
	Render()
}

func (this *SurferHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	this.Response = rw
	this.Request = req

	this.Prepare()
	switch r.Method {
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

	if !this.finished {
		this.Render()
	}

}

func (this *SurferHandler) Get() {
	// TODO 405 Error, Method Not Allowed
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.Path() and Route.HandlerFunc().
// func (r *mux.Router) HandleCtx(path string, f func(http.ResponseWriter,
// 	*http.Request)) *Route {
// 	return r.NewRoute().Path(path).HandlerFunc(f)
// }
