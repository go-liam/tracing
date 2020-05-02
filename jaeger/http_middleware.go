package jaeger

import (
	"context"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
	"log"
	"net/http"
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
	log.Printf("req %+v\n",req)
	req2 := req.WithContext(ctx)
	log.Printf("req2 %+v\n",req2)
	// wrap the request in nethttp.TraceRequest
	req3, ht := nethttp.TraceRequest(sv.Tracer , req2)
	log.Printf("req3 %+v\n",req3)
	defer ht.Finish()
	return span,req3
}

func (sv *SvJeager) OnError(span opentracing.Span, err error) {
	log.Print("[ERROR]OnError ",err)
	// handle errors by recording them in the span
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
}