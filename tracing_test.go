package tracing

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-liam/tracing/config"
	"log"
	"net/http"
	"testing"
)

func TestInit(t *testing.T) {
	tracingInit()
}

func tracingInit()  {
	Init(&config.TraceConfig{IsOpen: false, HostPort: "127.0.0.1:6831", SamplerType: "const", SamplerParam: 0.01, LogSpans: true})
}

func TestConfig(t *testing.T) {
	tracingInit()
	v1 := Config()
	log.Printf("config=%+v\n",v1)
}

func TestNewGinMiddlewareHandle(t *testing.T) {
	tracingInit()
NewGinMiddlewareHandle("serverName")
}

func TestCreateGinToGRPCConnect(t *testing.T) {
	tracingInit()
	c := new(gin.Context)
	CreateGinToGRPCConnect("127.0.0.1:6831",c ,false)
}

func TestHttpGinTracerRequestInject(t *testing.T) {
	tracingInit()
	c := &gin.Context{}
	//req := new(http.Header)
	req := &http.Header{
		"Content-Type":   {"text/html; charset=UTF-8"},
		"Content-Length": {"0"},
	}
	HttpGinTracerRequestInject(c,req)
}

func TestNewGRPCServerOption(t *testing.T) {
	tracingInit()
	NewGRPCServerOption("serviceName")
}

func TestCreateGRPCConnectOpts(t *testing.T) {
	tracingInit()
	var ctx  context.Context
	CreateGRPCConnectOpts(&ctx)
}

func TestHttpTraceRequestInject(t *testing.T) {
	tracingInit()
	var ctx  context.Context
	req := &http.Header{
		"Content-Type":   {"text/html; charset=UTF-8"},
		"Content-Length": {"0"},
	}
	HttpTraceRequestInject(&ctx,req)
}
