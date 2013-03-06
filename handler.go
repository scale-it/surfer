package surfer

import (
	"net/http"
	"strings"
	"strconv"
)

type statet int8

var states = struct {
	rendered statet
	finished statet
}{1, 2}


type Handler struct {
	app      *App
	//Template *Template
	//Session  *Session
	template string
	Response http.ResponseWriter
	Request  *http.Request
	Data     map[interface{}]interface{} // data for template
	Params   map[string]string
	state statet
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
	var supported = [3]string{"application/json","text/html","application/xml"}
	accepted := strings.Split(mediatype, ",")

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
	this.state = states.finished
	return true
}

func (this *Handler) Render() {
	this.state = states.rendered
	if this.data == nil {
		// todo: convert status 200 -> 204 No Content
		return;
	}
	// validate template against header accept
	// render based on template (filename, "json"...) and this.Data
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




// HandleFunc registers a new route with a matcher for the URL path.
// See Route.Path() and Route.HandlerFunc().
// func (r *mux.Router) HandleCtx(path string, f func(http.ResponseWriter,
// 	*http.Request)) *Route {
// 	return r.NewRoute().Path(path).HandlerFunc(f)
// }
