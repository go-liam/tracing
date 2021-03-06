package jaeger

import (
	"context"
	"google.golang.org/grpc"
)

// CreateGRPCConnectOpts :注入的 opts
func (sv *SvJeager) CreateGRPCConnectOpts(ctx *context.Context) grpc.DialOption {
	if sv.config ==nil || !sv.config.IsOpen {
		return nil
	}
	return grpc.WithUnaryInterceptor(clientInterceptor(sv.Tracer, *ctx))
}
