package surfer

import (
	"github.com/gorilla/mux"
	"github.com/kylelemons/go-gypsy/yaml" // yaml parser
	"github.com/scale-it/go-log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
)

type App struct {
	Config   *yaml.File
	Router   *mux.Router
	Log      *log.Logger
	listener net.Listener
}

// Creates new surfer application.
func New() App {
	config, err := yaml.ReadFile("config.yaml")
	if err != nil {
		path, _ := os.Getwd()
		println("Can't read configuration file `config.yaml` from from", path)
		panic(err)
	}
	log_level, err := config.Get("log_level")
	level, err := log.String2Level(log_level)
	if err != nil {
		println("WARNING, ", err.Error())
	}
	l := log.NewStd(os.Stderr, level, "", log.Ldate|log.Lmicroseconds, true)

	l.Debug("App initialized")
	return App{config, mux.NewRouter(), l, nil}
}

// Blocks current gorotine and runs web server.
func (this App) Run() {
	addr, err := this.Config.Get("server.addr")
	if err != nil {
		this.Log.Fatal("'server.addr' not specified in config file")
		panic(err)
	}
	port, err := this.Config.Get("server.port")
	if err != nil {
		this.Log.Fatal("'server.port' not specified in config file")
		panic(err)
	}
	typ, err := this.Config.Get("server.type")
	if err != nil {
		this.Log.Warning("'server.type' not specified in config file. Using default type - http")
	}
	listenAddr := net.JoinHostPort(addr, port)
	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic("Listen error: " + err.Error())
	}
	this.listener = l

	switch typ {
	case "fcgi":
		this.Log.Info("FCGI server, %s:%s", addr, port)
		fcgi.Serve(l, this.Router)
	default:
		if typ != "" && typ != "http" {
			this.Log.Warning("Wrong server type: %s. Should be one of: (http, fcgi). Using default type - http", typ)
		}
		this.Log.Info("HTTP server, http://%s:%s", addr, port)
		http.Serve(l, this.Router)
	}
}

// Stops web server.
func (this App) Stop() {
	this.listener.Close()
	this.Log.Info("server stopped")
}
