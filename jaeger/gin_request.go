package jaeger

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"log"
	"net/http"
)

// HttpGinTracerRequestInject :  在gin 请求中注入
func (sv *SvJeager) HttpTracerGinRequestInject(c *gin.Context, reqHeader *http.Header) {
	if !sv.IsOpen {
		return
	}
	tracer, _ := c.Get("Tracer")
	parentSpanContext, _ := c.Get("ParentSpanContext")
	span := opentracing.StartSpan(
		"call Gin Http",
		opentracing.ChildOf(parentSpanContext.(opentracing.SpanContext)),
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		ext.SpanKindRPCClient,
	)
	span.Finish()
	injectErr := tracer.(opentracing.Tracer).Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(*reqHeader))
	if injectErr != nil {
		log.Fatalf("[ERROR] %s: Couldn't inject headers", injectErr)
	}
}
