package jaeger

import (
	"context"
	grpcMiddeware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
)

func (sv *SvJeager) NewMiddlewareHandle(serverName string) http.Handler {
	if sv.config == nil || !sv.config.IsOpen {
		return nil
	}
	tracer, closer, _ := sv.NewJaegerTracer(serverName, sv.config.HostPort)
	defer closer.Close()
	return nethttp.Middleware(tracer, http.DefaultServeMux)
}

func (sv *SvJeager) HttpTracerRequestInject( req *http.Request,serverName string) (opentracing.Span,*http.Request) {
	if sv.config == nil || !sv.config.IsOpen {
		//return nil
	}
	// create a top-level span to represent full work of the client
	span := sv.Tracer.StartSpan(serverName)
	span.SetTag(string(ext.Component), serverName)
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	req2 := req.WithContext(ctx)
	// wrap the request in nethttp.TraceRequest
	req3, ht := nethttp.TraceRequest(sv.Tracer , req2)
	defer ht.Finish()
	return span,req3
}

func (sv *SvJeager) OnError(span opentracing.Span, err error) {
	log.Print("[ERROR]OnError ",err)
	// handle errors by recording them in the span
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
}

func (sv *SvJeager) CreateToGRPCConnect(serviceAddress string, openTrace bool,serverName string) *grpc.ClientConn {
	var conn *grpc.ClientConn
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	if openTrace && sv.config !=nil && sv.config.IsOpen {
		// tracer
		span := sv.Tracer.StartSpan(serverName)
		span.SetTag(string(ext.Component), serverName)
		defer span.Finish()

		//tracer := req.Header.Get("Tracer")
		//parentSpanContext := req.Header.Get("ParentSpanContext")
		conn, err = grpc.DialContext(
			ctx,
			serviceAddress,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(
				grpcMiddeware.ChainUnaryClient(
					// trace grpcMiddeware
					//clientInterceptorFromHttp(tracer.(opentracing.Tracer), parentSpanContext.(opentracing.SpanContext)),
					//clientInterceptorFromHttp(sv.Tracer, span.(opentracing.SpanContext)),
					clientInterceptor(sv.Tracer,ctx ),
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
