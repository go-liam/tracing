package trace

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-liam/tracing/config"
	"google.golang.org/grpc"
	"net/http"
)

type InTrace interface {
	Init(info *config.TraceConfig)
	Config() *config.TraceConfig // 获取Config
	//gin 服务器
	NewGinMiddlewareHandle(serviceName string) gin.HandlerFunc
	CreateGinToGRPCConnect(serviceAddress string, c *gin.Context, openTrace bool) *grpc.ClientConn // grpc trace
	HttpGinTracerRequestInject(c *gin.Context, reqHeader *http.Header)                             // http trace
	// grpc 服务器
	NewGRPCServerOption(serviceName string) grpc.ServerOption
	CreateGRPCConnectOpts(ctx *context.Context) grpc.DialOption          // grpc trace
	HttpTraceRequestInject(ctx *context.Context, reqHeader *http.Header) // http trace
}
