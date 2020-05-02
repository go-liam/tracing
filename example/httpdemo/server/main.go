package main

import (
	"fmt"
	"github.com/go-liam/tracing/config"
	"github.com/go-liam/tracing/jaeger"
	"github.com/go-liam/util/request"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	port = ":7201"
	serverName ="api-app"
	hostPort = "127.0.0.1:6831"
)

func main()  {
	println("run server port ", port)
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/api-app", helloHandler)
	http.HandleFunc("/api-app/:id", helloHandler)
	sv := new(jaeger.SvJeager)
	sv.Init(config.DefaultConfig)
	tracer, closer, _ := sv.NewJaegerTracer(serverName, hostPort)
	defer closer.Close()
	err := http.ListenAndServe(port,
		// use nethttp.Middleware to enable OpenTracing for server
		nethttp.Middleware(tracer, http.DefaultServeMux))
	if err != nil {
		log.Fatal("ListenAndServeTLS:", err.Error())
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("head: %+v\n",r.Header)
	name, _ := os.Hostname()
	path := r.URL.Path
	st1:="%s: host: %s, path: %s ,token: %s ,UserID: %s ,ip:%s"
	st := fmt.Sprintf(st1,serverName, name,path,r.Header.Get("Token"),r.Header.Get("UserID"),request.ClientIP(r))
	io.WriteString(w, st)
}
