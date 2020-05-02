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

func (sv *SvJeager) HttpTracerRequestInject( req *http.Request,serverName string) {
	if sv.config == nil || !sv.config.IsOpen {
		return
	}
	// create a top-level span to represent full work of the client
	span := sv.Tracer.StartSpan(serverName)
	span.SetTag(string(ext.Component), serverName)
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	req = req.WithContext(ctx)

	// wrap the request in nethttp.TraceRequest
	req, ht := nethttp.TraceRequest(sv.Tracer , req)
	defer ht.Finish()
}

func (sv *SvJeager) OnError(span opentracing.Span, err error) {
	// handle errors by recording them in the span
	span.SetTag(string(ext.Error), true)
	span.LogKV(otlog.Error(err))
	log.Print("[ERROR]",err)
}