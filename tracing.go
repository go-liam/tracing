package tracing

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-liam/tracing/config"
	"github.com/go-liam/tracing/jaeger"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"net/http"
)

var Sv InTrace

func init() {
	// 接入 jaeger
	Sv = new(jaeger.SvJeager)
	if Sv.Config() == nil || Sv.Config().HostPort == "" {
		println("[WARNING]Must set configuration information")
		Sv.Init(config.DefaultConfig)
	}
}

// Init :
func Init(info *config.TraceConfig) {
	println("[INFO]Set Tracing configuration information")
	Sv.Init(info)
}

func Config() *config.TraceConfig {
	return Sv.Config()
}

// NewGinMiddlewareHandle :
func NewGinMiddlewareHandle(serviceName string) gin.HandlerFunc {
	return Sv.NewGinMiddlewareHandle(serviceName)
}

// CreateGinToGRPCConnect :
func CreateGinToGRPCConnect(serviceAddress string, c *gin.Context, openTrace bool) *grpc.ClientConn {
	return Sv.CreateGinToGRPCConnect(serviceAddress, c, openTrace)
}

// HttpGinTracerRequestInject :
func HttpGinTracerRequestInject(c *gin.Context, reqHeader *http.Header) {
	Sv.HttpGinTracerRequestInject(c, reqHeader)
}

// NewGRPCServerOption :
func NewGRPCServerOption(serviceName string) grpc.ServerOption {
	return Sv.NewGRPCServerOption(serviceName)
}

// CreateGRPCConnectOpts :
func CreateGRPCConnectOpts(ctx *context.Context) grpc.DialOption {
	return Sv.CreateGRPCConnectOpts(ctx)
}

// HttpTraceRequestInject :
func HttpTraceRequestInject(ctx *context.Context, reqHeader *http.Header) {
	Sv.HttpTraceRequestInject(ctx, reqHeader)
}

func NewMiddlewareHandle(serverName string) http.Handler {
	return Sv.NewMiddlewareHandle(serverName)
}

func HttpTracerRequestInject(req *http.Request, serverName string) (opentracing.Span, *http.Request) {
	return Sv.HttpTracerRequestInject(req, serverName)
}

func OnError(span opentracing.Span, err error) {
	Sv.OnError(span, err)
}
