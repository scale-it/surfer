package main

import (
	"github.com/scale-it/surfer"
)

type SessionHandler struct {
	surfer.Handler
	Session *Session
}

type Index struct {
	SessionHandler
}

func main() {
	app := surfer.New()
	app.Run()
}
