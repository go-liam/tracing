package jaeger

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// NewMiddlewareGinHandle : AppName = "go-gin-api"
func (sv *SvJeager) NewMiddlewareGinHandle(appName string) gin.HandlerFunc {
	if !sv.IsOpen {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	return func(c *gin.Context) {
		var parentSpan opentracing.Span
		tracer, closer, _ := sv.NewJaegerTracer(appName, sv.HostPort)
		defer closer.Close()
		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			parentSpan = tracer.StartSpan(c.Request.URL.Path)
			defer parentSpan.Finish()
		} else {
			parentSpan = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
			defer parentSpan.Finish()
		}
		c.Set("Tracer", tracer)
		c.Set("ParentSpanContext", parentSpan.Context())
		c.Next()
	}
}