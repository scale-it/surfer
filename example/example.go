package main

import (
	"github.com/scale-it/surfer"
	"time"
)

func main() {
	app := surfer.New()
	app.Run()
	time.Sleep(time.Second / 4)

}
