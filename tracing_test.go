package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/go-liam/tracing/config"
	"log"
	"testing"
)

func TestInit(t *testing.T) {
	Init( &config.TraceConfig{IsOpen: true,HostPort: "127.0.0.1:6831",SamplerType: "const",SamplerParam: 0.01,LogSpans: true})
}

func TestConfig(t *testing.T) {
	v1 := Config()
	log.Printf("config=%+v\n",v1)
}

func TestNewGinMiddlewareHandle(t *testing.T) {
NewGinMiddlewareHandle("serverName")
}

func TestCreateGinToGRPCConnect(t *testing.T) {
	c := new(gin.Context)
	CreateGinToGRPCConnect("127.0.0.1:6831",c ,false)
}

func TestHttpGinTracerRequestInject(t *testing.T) {
	//c := &gin.Context{}
	////req := new(http.Header)
	//req := &http.Header{
	//	"Content-Type":   {"text/html; charset=UTF-8"},
	//	"Content-Length": {"0"},
	//}
	//HttpGinTracerRequestInject(c,req)
}

func TestNewGRPCServerOption(t *testing.T) {
}

func TestCreateGRPCConnectOpts(t *testing.T) {
}

func TestHttpTraceRequestInject(t *testing.T) {
}
