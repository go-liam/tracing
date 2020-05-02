package main

import (
	"fmt"
	trace2 "github.com/go-liam/tracing"
	"github.com/go-liam/tracing/example/pkg/config"
	"github.com/go-liam/tracing/example/proto/write"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	ServiceName     = "liam-grpc4-write"
	ServiceHostPort = config.PortGrpc4Write //"0.0.0.0:9904"
)

func main() {
	l, err2 := net.Listen("tcp", ServiceHostPort)
	if err2 != nil {
		log.Fatalf("Failed to listen: %v", err2)
	}
	var s *grpc.Server
	if trace2.Config().IsOpen {
		s = grpc.NewServer(trace2.NewGRPCServerOption(ServiceName)) // 注入trace
	} else {
		s = grpc.NewServer()
	}
	// 服务注册
	write.RegisterWriteServer(s, &WriteController{})
	log.Println("Listen on " + ServiceHostPort)
	reflection.Register(s)
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type WriteController struct{}

func (l *WriteController) WriteData(ctx context.Context, in *write.Request) (*write.Response, error) {
	// 调用 gRPC 服务
	//grpcListenClient := listen.NewListenClient(CreateServiceListenConn(ctx))
	//resListen, _ := grpcListenClient.ListenData(context.Background(), &listen.Request{Name: "listen"})

	return &write.Response{Message: fmt.Sprintf("[%s]", in.Name)}, nil
}
