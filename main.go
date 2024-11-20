//go:build !wasm

package main

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jpillora/installer/handler"
	"github.com/jpillora/opts"
)

var version = "0.0.0-src"

func main() {
	c := handler.DefaultConfig
	opts.New(&c).Repo("github.com/xqbumu/worker-installer").Version(version).Parse()
	log.Printf("default user is '%s'", c.User)
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	h := &handler.Handler{Config: c}
	s := &http.Server{
		Addr:           addr,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Ready on %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
	log.Print("exiting")
}
