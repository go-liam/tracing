package jaeger

import (
	"context"
	"fmt"
	"github.com/go-liam/tracing/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	"strings"

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
	model.TraceModel
}

// 属性
//type TraceModel struct {
//	IsOpen       bool // 开关
//	HostPort     string //  "127.0.0.1:6831"
//	SamplerType  string  //固定采样
//	SamplerParam float64 //1=全采样、0=不采样
//	LogSpans     bool    // 打印日志
//}

func (sv *SvJeager) SetAttributes(info *model.TraceModel) {
	sv.IsOpen = info.IsOpen
	sv.HostPort = info.HostPort
	sv.LogSpans = info.LogSpans
	sv.SamplerParam = info.SamplerParam
	sv.SamplerType = info.SamplerType
}

func (sv *SvJeager) Attributes() *model.TraceModel {
	info := new(model.TraceModel)
	info.IsOpen = sv.IsOpen
	info.HostPort = sv.HostPort
	info.LogSpans = sv.LogSpans
	info.SamplerParam = sv.SamplerParam
	info.SamplerType = sv.SamplerType
	return info
}

// ServerOption grpc server option
//func ServerOption(tracer opentracing.Tracer) grpc.ServerOption {
//	return grpc.UnaryInterceptor(serverInterceptor(tracer))
//}

var Tracer opentracing.Tracer

func (sv *SvJeager) NewJaegerTracer(serviceName string, jaegerHostPort string) (opentracing.Tracer, io.Closer, error) {
	// 采样设置：https://www.jaegertracing.io/docs/1.17/sampling/
	cfg := &jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  sv.SamplerType,  //"const", //固定采样
			Param: sv.SamplerParam, //1,       //1=全采样、0=不采样
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           sv.LogSpans,
			LocalAgentHostPort: jaegerHostPort,
		},

		ServiceName: serviceName,
	}

	var closer io.Closer
	var err error
	Tracer, closer, err = cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(Tracer)
	return Tracer, closer, err
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

// MDReaderWriter :
type MDReaderWriter struct {
	metadata.MD
}

// ForeachKey : implements Foreach Key of opentracing.TextMapReader
func (c MDReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range c.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Set implements Set() of opentracing.TextMapWriter
func (c MDReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	c.MD[key] = append(c.MD[key], val)
}
