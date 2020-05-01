package trace

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-liam/tracing/model"
	"google.golang.org/grpc"
	"net/http"
)

type InTrace interface {
	//gin 服务器
	NewMiddlewareGinHandle(appName string) gin.HandlerFunc
	CreateGRPCConnectFromGin(serviceAddress string, c *gin.Context, openTrace bool) *grpc.ClientConn // grpc trace
	HttpTracerGinRequestInject(c *gin.Context, reqHeader *http.Header)                               // http trace
	// grpc 服务器
	NewGRPCServerOption(serviceName string) grpc.ServerOption
	CreateGRPCConnectOpts(ctx *context.Context) grpc.DialOption          // grpc trace
	HttpTraceRequestInject(ctx *context.Context, reqHeader *http.Header) // http trace
	// other
	Attributes() *model.TraceModel        // 获取属性
	SetAttributes(info *model.TraceModel) // 设置属性
}
