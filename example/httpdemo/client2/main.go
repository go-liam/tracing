package main

import (
	"context"
	"fmt"
	config2 "github.com/go-liam/tracing/config"
	"github.com/go-liam/tracing/example/pkg/config"
	"github.com/go-liam/tracing/example/proto/listen"
	"github.com/go-liam/tracing/jaeger"
	"github.com/go-liam/tracing/request"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	port = ":7203"
	serverName ="api-ai2"
	hostPort = "127.0.0.1:6831"
)
var sv *jaeger.SvJeager

func main()  {
	println("run server port ", port)
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/api-app", helloHandler)
	http.HandleFunc("/api-app/:id", helloHandler)
	sv = new(jaeger.SvJeager)
	sv.Init(config2.DefaultConfig)
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
	//runClient(sv.Tracer)
	request.TraceRequest("GET","http://localhost:7201/",nil,sv,serverName)
	//connectGrpc()
	log.Printf("head: %+v\n",r.Header)
	name, _ := os.Hostname()
	path := r.URL.Path
	st1:="%s: 33host: %s, path: %s ,token: %s ,UserID: %s "
	st := fmt.Sprintf(st1,serverName, name,path,r.Header.Get("Token"),r.Header.Get("UserID"))
	io.WriteString(w, st)
}

// 不成功，待处理；
func connectGrpc()  {
	println("connectGrpc --------------- ")
	conn := sv.CreateToGRPCConnect(config.UrlGrpc1Listen,  true,serverName) // grpc_client.CreateServiceListenConn(c)
	grpcListenClient := listen.NewListenClient(conn)
	resListen, _ := grpcListenClient.ListenData(context.Background(), &listen.Request{Name: "listen"})
	//conn := sv.CreateToGRPCConnect(config.UrlGrpc2Read, true, serverName)
	//grpcReadClient := read.NewReadClient(conn)
	//resRead, _ := grpcReadClient.ReadData(context.Background(), &read.Request{Name: "read"})
	log.Printf("grpc:result: %+v\n",resListen)
}