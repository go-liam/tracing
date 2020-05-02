package jaeger

import (
	"github.com/go-liam/tracing/example/pkg/config"
	"google.golang.org/grpc"
	"log"
)

// NewGRPCServerOption :
func (sv *SvJeager) NewGRPCServerOption(serviceName string) grpc.ServerOption {
	if sv.config ==nil || !sv.config.IsOpen {
		return nil
	}
	tracer, _, err := sv.NewJaegerTracer(serviceName, config.JaegerHostPort)
	if err != nil {
		log.Printf("new tracer err: %+v\n", err)
		//os.Exit(-1)
		return nil
	}
	return grpc.UnaryInterceptor(serverInterceptor(tracer))
}
