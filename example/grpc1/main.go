package main

import (
	"fmt"
	trace2 "github.com/go-liam/tracing"
	"github.com/go-liam/tracing/example/pkg/config"
	"github.com/go-liam/tracing/example/proto/listen"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	ServiceName     = "liam-gRPC1-Listen"
	ServiceHostPort = config.PortGrpc1Listen //"0.0.0.0:9901"
)

func main() {
	l, err2 := net.Listen("tcp", ServiceHostPort)
	if err2 != nil {
		log.Fatalf("Failed to listen: %v", err2)
	}
	var s *grpc.Server
	if trace2.Server.Config().IsOpen {
		s = grpc.NewServer(trace2.Server.NewGRPCServerOption(ServiceName)) // 注入trace
	} else {
		s = grpc.NewServer()
	}
	// 服务注册
	listen.RegisterListenServer(s, &ListenController{})
	log.Println(ServiceName + " Listen on " + ServiceHostPort)
	reflection.Register(s)
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// ListenController :
type ListenController struct{}

func (l *ListenController) ListenData(ctx context.Context, in *listen.Request) (*listen.Response, error) {
	return &listen.Response{Message: fmt.Sprintf("[%s]", in.Name)}, nil
}
