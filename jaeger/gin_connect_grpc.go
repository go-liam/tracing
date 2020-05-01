package jaeger

import (
	"context"
	"github.com/gin-gonic/gin"
	grpcMiddeware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"log"
	"time"
)

// CreateGRPCConnectFromGin :创建grpc客户端链接
func (sv *SvJeager) CreateGRPCConnectFromGin(serviceAddress string, c *gin.Context, openTrace bool) *grpc.ClientConn {
	var conn *grpc.ClientConn
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	if openTrace && sv.IsOpen {
		// tracer
		tracer, _ := c.Get("Tracer")
		parentSpanContext, _ := c.Get("ParentSpanContext")
		conn, err = grpc.DialContext(
			ctx,
			serviceAddress,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(
				grpcMiddeware.ChainUnaryClient(
					// trace grpcMiddeware
					clientInterceptorFromHttp(tracer.(opentracing.Tracer), parentSpanContext.(opentracing.SpanContext)),
					//grpc_log.clientInterceptor(),
				),
			),
		)
	} else {
		conn, err = grpc.DialContext(
			ctx,
			serviceAddress,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(
				grpcMiddeware.ChainUnaryClient(
					//grpc_log.clientInterceptor(),
				),
			),
		)
	}
	if err != nil {
		log.Println(serviceAddress, "[ERROR] grpc conn err:", err)
	}
	return conn
}
