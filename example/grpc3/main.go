package main

import (
	"fmt"
	trace2 "github.com/go-liam/tracing"
	"github.com/go-liam/tracing/example/pkg/config"
	"github.com/go-liam/tracing/example/pkg/request"
	"github.com/go-liam/tracing/example/proto/speak"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	ServiceName     = "liam-grpc3-speak"
	ServiceHostPort = config.PortGrpc3Speak // "0.0.0.0:9903"
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
	speak.RegisterSpeakServer(s, &SpeakController{})
	log.Println(ServiceName + " Listen on " + ServiceHostPort)
	reflection.Register(s)
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type SpeakController struct{}

func (l *SpeakController) SpeakData(ctx context.Context, in *speak.Request) (*speak.Response, error) {
	// 调用 HTTP 服务
	resHttpGet := ""
	_, err := request.HttpGet(config.UrlRest2, ctx, true)
	if err == nil {
		resHttpGet = "[HttpGetOk]"
	}
	msg := "[" + fmt.Sprintf("%s", in.Name) + "-" +
		resHttpGet +
		"]"
	println("msg:", msg)
	return &speak.Response{Message: fmt.Sprintf("[%s]", in.Name)}, nil
}
