package main

import (
	"context"
	"github.com/gin-gonic/gin"
	trace2 "github.com/go-liam/tracing"
	"github.com/go-liam/tracing/example/pkg/config"
	"github.com/go-liam/tracing/example/pkg/request"
	"github.com/go-liam/tracing/example/proto/listen"
	"github.com/go-liam/tracing/example/proto/read"
	"github.com/go-liam/tracing/example/proto/speak"
	"github.com/go-liam/tracing/example/proto/write"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	port       = ":7410"
	serverName = "liam-api-gateway"
)

func main() {
	println("run api port ", port)
	engine := gin.New()
	// 设置路由
	SetupRouter(engine)
	engine.Run(port)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
}

// SetupRouter :
func SetupRouter(engine *gin.Engine) {
	//设置路由jaeger中间件
	engine.Use(trace2.NewGinMiddlewareHandle(serverName))
	engine.GET("/", Index)
	engine.GET("/api/trace", JaegerTest)               // 多微服trace
	engine.GET("/api/trace/gin", JaegerTestGinRequest) // 简单http 跟踪demo
	engine.GET("/api/trace/grpc", JaegerTestGrpc)      // 简单grpc 跟踪demo
}

// JaegerTestGinRequest : 测试 gin 中请求其它http资源
func JaegerTestGinRequest(c *gin.Context) {
	resHttpGet := ""
	rest1, err := request.HttpGetGin(config.UrlRest1+"/", c, true)
	if err == nil {
		println("Rest1 =", rest1)
		resHttpGet += "-[Rest1 Ok-/api/trace]"
	} else {
		log.Printf("Rest1 err= %+v\n", err)
	}
	c.String(http.StatusOK, "Msg : "+resHttpGet)
}

func JaegerTestGrpc(c *gin.Context) {
	conn := trace2.CreateGinToGRPCConnect(config.UrlGrpc2Read, c, true)
	grpcReadClient := read.NewReadClient(conn)
	resRead, _ := grpcReadClient.ReadData(context.Background(), &read.Request{Name: "read"})

	c.String(http.StatusOK, "Msg : "+resRead.Message)
}

// JaegerTest :
func JaegerTest(c *gin.Context) {
	// 调用 gRPC 服务 listen
	conn := trace2.CreateGinToGRPCConnect(config.UrlGrpc1Listen, c, true) // grpc_client.CreateServiceListenConn(c)
	grpcListenClient := listen.NewListenClient(conn)
	resListen, _ := grpcListenClient.ListenData(context.Background(), &listen.Request{Name: "listen"})

	// 调用 gRPC 服务 speak
	conn = trace2.CreateGinToGRPCConnect(config.UrlGrpc3Speak, c, true)
	grpcSpeakClient := speak.NewSpeakClient(conn)
	resSpeak, _ := grpcSpeakClient.SpeakData(context.Background(), &speak.Request{Name: "speak"})

	// 调用 gRPC 服务 read
	conn = trace2.CreateGinToGRPCConnect(config.UrlGrpc2Read, c, true)
	grpcReadClient := read.NewReadClient(conn)
	resRead, _ := grpcReadClient.ReadData(context.Background(), &read.Request{Name: "read"})

	// 调用 gRPC 服务 write
	conn = trace2.CreateGinToGRPCConnect(config.UrlGrpc4Write, c, true)
	grpcWriteClient := write.NewWriteClient(conn)
	resWrite, _ := grpcWriteClient.WriteData(context.Background(), &write.Request{Name: "write"})

	resHttpGet := ""
	rest1, err := request.HttpGetGin(config.UrlRest1+"/api/trace", c, true)
	if err == nil {
		println("Rest1 =", rest1)
		resHttpGet += "-[Rest1 Ok-/api/trace]"
	} else {
		log.Printf("Rest1 err= %+v\n", err)
	}

	rest2, err2 := request.HttpGetGin(config.UrlRest2+"/api/trace", c, true)
	if err2 == nil {
		println("rest2 =", rest2)
		resHttpGet += "-[rest2 Ok-7422]"
	} else {
		log.Printf("rest1 err= %+v\n", err2)
	}
	msg := resListen.Message + "-" +
		resSpeak.Message + "-" +
		resRead.Message + "-" +
		resWrite.Message + "-" +
		resHttpGet
	c.String(http.StatusOK, "Msg : "+msg)
}

// Index :
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello,It works.index ")
}
