package main

import (
	trace2 "github.com/go-liam/tracing"
	"github.com/go-liam/tracing/example/grpc2/core"
	"github.com/go-liam/tracing/example/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	ServiceName     = "liam-gRPC2-Read"
	ServiceHostPort = config.PortGrpc2Read // "0.0.0.0:9902"
)

func main() {
	// grpc
	l, err := net.Listen("tcp", ServiceHostPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	var s *grpc.Server
	if trace2.Config().IsOpen {
		s = grpc.NewServer(trace2.NewGRPCServerOption(ServiceName)) // 注入trace
	} else {
		s = grpc.NewServer()
	}
	// 服务注册
	//read.RegisterReadServer(s, &core.ReadController{})
	core.RegisterServer(s)

	log.Println(ServiceName + " Listen on " + ServiceHostPort)
	reflection.Register(s)
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
