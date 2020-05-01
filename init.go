package trace

import (
	"github.com/go-liam/tracing/jaeger"
)

var Server InTrace

func init() {
	// 接入 jaeger
	Server = jaeger.Server
}
