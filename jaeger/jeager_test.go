package jaeger

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-liam/tracing/config"
	"log"
	"net/http"
	"testing"
)

var sv = new(SvJeager)

//func TestMain(m *testing.M) {
//	//sv.Init(&config.TraceConfig{IsOpen: true, HostPort: "127.0.0.1:6831", SamplerType: "const", SamplerParam: 0.01, LogSpans: true})
//}

func TestJeagerInit(t *testing.T) {
	jeagerInit()
}

func jeagerInit()  {
	sv.Init(&config.TraceConfig{IsOpen: false, HostPort: "127.0.0.1:6831", SamplerType: "const", SamplerParam: 0.01, LogSpans: true})
}

func TestConfig(t *testing.T) {
	jeagerInit()
	v1 := sv.Config()
	log.Printf("config=%+v\n", v1)
}

func TestNewGinMiddlewareHandle(t *testing.T) {
	jeagerInit()
	sv.NewGinMiddlewareHandle("serverName")
}

func TestCreateGinToGRPCConnect(t *testing.T) {
	jeagerInit()
	c := new(gin.Context)
	sv.CreateGinToGRPCConnect("127.0.0.1:6831", c, false)
}

func TestHttpGinTracerRequestInject(t *testing.T) {
	jeagerInit()
	c := &gin.Context{}
	//req := new(http.Header)
	req := &http.Header{
		"Content-Type":   {"text/html; charset=UTF-8"},
		"Content-Length": {"0"},
	}
	sv.HttpGinTracerRequestInject(c,req)
}

func TestNewGRPCServerOption(t *testing.T) {
	jeagerInit()
	sv.NewGRPCServerOption("serviceName")
}

func TestCreateGRPCConnectOpts(t *testing.T) {
	jeagerInit()
	var ctx  context.Context
	sv.CreateGRPCConnectOpts(&ctx)
}

func TestHttpTraceRequestInject(t *testing.T) {
	jeagerInit()
	var ctx  context.Context
	req := &http.Header{
		"Content-Type":   {"text/html; charset=UTF-8"},
		"Content-Length": {"0"},
	}
	sv.HttpTraceRequestInject(&ctx,req)
}
