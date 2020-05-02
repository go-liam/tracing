package jaeger

import (
	"google.golang.org/grpc/metadata"
	"strings"
)

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
