//go:build wasm

package main

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=handler

import (
	"log"

	"github.com/jpillora/installer/handler"
	"github.com/jpillora/opts"
	"github.com/syumai/workers"
)

var version = "0.0.0-src"

func main() {
	c := handler.DefaultConfig
	opts.New(&c).Repo("github.com/xqbumu/worker-installer").Version(version).Parse()
	log.Printf("default user is '%s'", c.User)
	h := &handler.Handler{Config: c}
	workers.Serve(h)
	log.Print("exiting")
}
