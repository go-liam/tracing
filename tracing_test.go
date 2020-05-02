package tracing

import (
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

}

func TestCreateGinToGRPCConnect(t *testing.T) {

}

func TestHttpGinTracerRequestInject(t *testing.T) {
}

func TestNewGRPCServerOption(t *testing.T) {
}

func TestCreateGRPCConnectOpts(t *testing.T) {
}

func TestHttpTraceRequestInject(t *testing.T) {
}
