package jaeger

import (
	"context"
	"google.golang.org/grpc"
)

// CreateGRPCConnectOpts :注入的 opts
func (sv *SvJeager) CreateGRPCConnectOpts(ctx *context.Context) grpc.DialOption {
	if !sv.config.IsOpen {
		return nil
	}
	return grpc.WithUnaryInterceptor(clientInterceptor(Tracer, *ctx))
}
