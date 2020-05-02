package jaeger

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"log"
	"net/http"
)

// HttpTraceRequestInject : 在http 请求中注入
func (sv *SvJeager) HttpTraceRequestInject(ctx *context.Context, reqHeader *http.Header) {
	if sv.config ==nil || !sv.config.IsOpen {
		return
	}
	// trace: span
	span, _ := opentracing.StartSpanFromContext(
		*ctx,
		"call Http",
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		ext.SpanKindRPCClient,
	)
	span.Finish()
	// trace: trans
	injectErr := sv.Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(*reqHeader))
	if injectErr != nil {
		log.Printf("[ERROR] HttpTraceRequestInject= %+v\n", injectErr)
		//log.Fatalf("%s: Couldn't inject headers", err)
	}
}

