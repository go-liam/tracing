package main

import (
	"context"
	"fmt"
	"github.com/go-liam/tracing/config"
	"github.com/go-liam/tracing/jaeger"
	"github.com/go-liam/util/request"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	port = ":7203"
	serverName ="api-ai"
	hostPort = "127.0.0.1:6831"
)
var sv *jaeger.SvJeager

func main()  {
	println("run server port ", port)
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/api-app", helloHandler)
	http.HandleFunc("/api-app/:id", helloHandler)
	sv = new(jaeger.SvJeager)
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
	runClient(sv.Tracer)

	log.Printf("head: %+v\n",r.Header)
	name, _ := os.Hostname()
	path := r.URL.Path
	st1:="%s: host: %s, path: %s ,token: %s ,UserID: %s ,ip:%s"
	st := fmt.Sprintf(st1,serverName, name,path,r.Header.Get("Token"),r.Header.Get("UserID"),request.ClientIP(r))
	io.WriteString(w, st)
}

func runClient(tracer opentracing.Tracer) {
	// nethttp.Transport from go-stdlib will do the tracing
	c := &http.Client{Transport: &nethttp.Transport{}}

	// create a top-level span to represent full work of the client
	span := tracer.StartSpan(serverName)
	span.SetTag(string(ext.Component), serverName)
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	req, err := http.NewRequest(
		"GET",
		"http://localhost:7201/",
		nil,
	)
	if err != nil {
		onError(span, err)
		return
	}

	req = req.WithContext(ctx)
	// wrap the request in nethttp.TraceRequest
	req, ht := nethttp.TraceRequest(tracer, req)
	defer ht.Finish()

	res, err := c.Do(req)
	if err != nil {
		onError(span, err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		onError(span, err)
		return
	}
	fmt.Printf("Received result: %s\n", string(body))
}

func onError(span opentracing.Span, err error) {
	// handle errors by recording them in the span
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
	log.Print(err)
}