package main

import (
	"context"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/go-liam/tracing"
	"github.com/go-liam/tracing/example/pkg/config"
	"github.com/go-liam/tracing/example/proto/listen"
	"github.com/go-liam/tracing/model"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	engine := gin.New()
	trace.Server.SetAttributes(&model.TraceModel{IsOpen: true, HostPort: "127.0.0.1:6831", SamplerType: "const", SamplerParam: 0.01, LogSpans: true})
	// 中间件
	engine.Use(trace.Server.NewMiddlewareGinHandle("serverName"))
	engine.GET("/test/http", apiHttp)
	engine.GET("/test/grpc", apiGRPC)
	//engine.Run(":7504")
	println("end")
}

func apiHttp(c *gin.Context) {
	httpGetGin("http://localhost/", c, true)
	c.String(http.StatusOK, "Hello,It works.index ")
}

func apiGRPC(c *gin.Context) {
	conn := trace.Server.CreateGRPCConnectFromGin(config.UrlGrpc1Listen, c, true) // grpc_client.CreateServiceListenConn(c)
	grpcListenClient := listen.NewListenClient(conn)
	resListen, _ := grpcListenClient.ListenData(context.Background(), &listen.Request{Name: "listen"})
	c.String(http.StatusOK, "Hello, "+resListen.Message)
}

func httpGetGin(url string, c *gin.Context, openTrace bool) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Second * 5, //默认5秒超时时间
		Transport: tr,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	//[TRACE] 注入 tracer 传输
	if openTrace {
		trace.Server.HttpTracerGinRequestInject(c, &req.Header)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(content), err
}
