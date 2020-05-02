package jaeger

import (
	"context"
	"fmt"
	"github.com/go-liam/tracing/config"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	//"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"io"
)

// SvJeager ：
type SvJeager struct {
	config *config.TraceConfig
	Tracer opentracing.Tracer
}

func (sv *SvJeager) Init(info *config.TraceConfig) {
	sv.config = info
}

func (sv *SvJeager) Config() *config.TraceConfig {
	return sv.config
}

//var Tracer opentracing.Tracer

func (sv *SvJeager) NewJaegerTracer(serviceName string, jaegerHostPort string) (opentracing.Tracer, io.Closer, error) {
	// 采样设置：https://www.jaegertracing.io/docs/1.17/sampling/
	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  sv.config.SamplerType,  //"const", //固定采样
			Param: sv.config.SamplerParam, //1,       //1=全采样、0=不采样
		},
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           sv.config.LogSpans,
			LocalAgentHostPort: jaegerHostPort,
		},
		ServiceName: serviceName,
	}

	var closer io.Closer
	var err error
	sv.Tracer, closer, err = cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(sv.Tracer)
	return sv.Tracer, closer, err
}

// ServerInterceptor grpc server
func serverInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		var parentContext context.Context

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		spanContext, err := tracer.Extract(opentracing.TextMap, MDReaderWriter{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			grpclog.Errorf("extract from metadata err: %v", err)
		} else {
			span := tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()

			parentContext = opentracing.ContextWithSpan(ctx, span)
		}

		return handler(parentContext, req)
	}
}

// clientInterceptor grpc client
func clientInterceptor(tracer opentracing.Tracer, ctx1 context.Context) grpc.UnaryClientInterceptor {
	span, _ := opentracing.StartSpanFromContext(
		ctx1,
		"call gRPC",
		opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
		ext.SpanKindRPCClient,
	)
	return clientGRPCInterceptor(tracer, span)
}

// clientInterceptorFromHttp grpc client
func clientInterceptorFromHttp(tracer opentracing.Tracer, spanContext opentracing.SpanContext) grpc.UnaryClientInterceptor {
	span := opentracing.StartSpan(
		"call gRPC",
		opentracing.ChildOf(spanContext),
		opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
		ext.SpanKindRPCClient,
	)
	return clientGRPCInterceptor(tracer, span)
}

// clientGRPCInterceptor :
func clientGRPCInterceptor(tracer opentracing.Tracer, span opentracing.Span) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string,

		req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		defer span.Finish()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		mdWriter := MDReaderWriter{md}
		err := tracer.Inject(span.Context(), opentracing.TextMap, mdWriter)
		if err != nil {
			span.LogFields(traceLog.String("inject-error", err.Error()))
		}

		newCtx := metadata.NewOutgoingContext(ctx, md)
		err = invoker(newCtx, method, req, reply, cc, opts...)
		if err != nil {
			span.LogFields(traceLog.String("call-error", err.Error()))
		}
		return err
	}
}
