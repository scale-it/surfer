package surfer

import (
	"github.com/gorilla/mux"
	//"github.com/gorilla/sessions"
	"github.com/kylelemons/go-gypsy/yaml" // yaml parser
	"github.com/scale-it/clog"            // logger, http://godoc.org/github.com/cratonica/clog
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
)

type App struct {
	Config   *yaml.File
	Router   *mux.Router
	Log   *clog.Clog
	listener net.Listener
}

func New() App {
	config, err := yaml.ReadFile("config.yaml")
	if err != nil {
		path, _ := os.Getwd()
		println("Can't read configuration file `config.yaml` from from", path)
		panic(err)
	}
	log_level, err := config.Get("log_level")
	level, err := clog.String2Level(log_level)
	if err != nil {
		println("WARNING, ", err.Error())
	}
	log := clog.NewClog()
	log.AddOutput(os.Stdout, level)

	log.Debug("App initialized")
	return App{config, mux.NewRouter(), log, nil}
}

func (this App) Run() {
	addr, err := this.Config.Get("server.addr")
	if err != nil {
		this.Log.Fatal("`server.addr` not specified in config file")
		panic(err)
	}
	port, err := this.Config.Get("server.port")
	if err != nil {
		this.Log.Fatal("`server.port` not specified in config file")
		panic(err)
	}
	typ, err := this.Config.Get("server.type")
	if err != nil {
		this.Log.Warning("`server.type` not specified in config file")
	}
	listenAddr := net.JoinHostPort(addr, port)
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic("Listen error: " + err.Error())
	}
	this.listener = l

	switch typ {
	case "http":
		this.Log.Info("HTTP server, %s:%s", addr, port)
		http.Serve(l, this.Router)
	case "fcgi":
		this.Log.Info("FCGI server, %s:%s", addr, port)
		fcgi.Serve(l, this.Router)
	default:
		if typ != "" {
			this.Log.Warning("Wrong server type: %s. Using HTTP", typ)
		}
		this.Log.Info("HTTP server, %s:%s", addr, port)
		http.Serve(l, this.Router)
	}
}

func (this App) Stop() {
	this.listener.Close()
}
