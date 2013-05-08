package surfer

import (
	"encoding/json"
	"fmt"
	"github.com/hinasssan/msgpack-go"
	"github.com/scale-it/go-log"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const (
	r_error = iota
	r_html
	r_json
	r_msgpack
	r_unknown
)

var r_map = map[string]int{
	"error":   r_error,
	"json":    r_json,
	"msgpack": r_msgpack,
}

/* Type of function to handle by renderer. It should return
* template name which is a fielname or renderer type.
  Renderer type is one of "error", "json", "msgpack". If template name doesn't match renderer
  type, then Renderer will read Accept HTTP Request header to get one. Otherwise Renderer type
  will be forced by template name. */
type RendererHandler func(w http.ResponseWriter, r *http.Request) (string, interface{}, int)

type WithRenderer struct {
	T   *template.Template
	Log *log.Logger
	H   RendererHandler
}

func (this WithRenderer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tname, data, status := this.H(w, r)
	switch negotiateRenderer(r, tname) {
	case r_error:
		http.Error(w, data.(string), status)
	case r_json:
		w.Header().Set("Content-Type", "application/json")
		content, err := json.Marshal(data)
		this.write(w, content, err, status)
	case r_msgpack:
		w.Header().Set("Content-Type", "application/x-msgpack")
		content, err := msgpack.Marshal(data)
		this.write(w, content, err, status)
	case r_html:
		w.Header().Set("Content-Type", "text/html")
		if err := this.T.ExecuteTemplate(w, tname, data); err != nil {
			this.write(w, nil, err, status)
		}
	default:
		w.Header().Set("Content-Type", "text/plain")
		this.write(w, []byte(fmt.Sprint(data)), nil, status)
	}
}

func (this WithRenderer) write(w http.ResponseWriter, data []byte, err error, status int) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		this.Log.Error(err.Error())
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

func negotiateRenderer(r *http.Request, tname string) int {
	if v, present := r_map[tname]; present {
		return v
	}
	for _, a := range strings.Split(r.Header.Get("Accept"), ",") {
		if strings.Contains(a, "json") {
			return r_json
		}
		if strings.Contains(a, "html") {
			return r_html
		}
		if strings.Contains(a, "msgpack") {
			return r_msgpack
		}
	}
	if tname[:4] == "html" || tname[:3] == "htm" {
		return r_html
	}
	return r_unknown
}
