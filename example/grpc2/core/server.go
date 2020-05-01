package core

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	trace2 "github.com/go-liam/tracing"
	"github.com/go-liam/tracing/example/pkg/config"
	"github.com/go-liam/tracing/example/pkg/request"
	"github.com/go-liam/tracing/example/proto/listen"
	"github.com/go-liam/tracing/example/proto/read"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

var Tracer1 *opentracing.Tracer

type ReadController struct{}

func (s *ReadController) ReadData(ctx context.Context, in *read.Request) (*read.Response, error) {
	// 调用 gRPC 服务
	var conn *grpc.ClientConn
	var err1 error
	if trace2.Server.Attributes().IsOpen {
		conn, err1 = grpc.Dial(config.UrlGrpc1Listen, grpc.WithInsecure(), trace2.Server.CreateGRPCConnectOpts(&ctx)) // 注入
	} else {
		conn, err1 = grpc.Dial(config.UrlGrpc1Listen, grpc.WithInsecure())
	}
	if err1 != nil {
		log.Println("grpc conn err:", err1)
	}
	grpcListenClient := listen.NewListenClient(conn)
	//grpcListenClient := listen.NewListenClient(trace.Server.CreateGRPCConnectFromGRPC(config.UrlGrpc1Listen,&ctx )) // 注入trace
	resListen, _ := grpcListenClient.ListenData(context.Background(), &listen.Request{Name: "listen"})
	// 调用 HTTP 服务
	resHttpGet := ""
	_, err := request.HttpGet(config.UrlRest2, ctx, true)
	if err == nil {
		resHttpGet = "[HttpGetOk]"
	}
	msg := "[" + fmt.Sprintf("%s", in.Name) + "-" +
		resListen.Message + "-" +
		resHttpGet +
		"]"
	return &read.Response{Message: msg}, nil
}

func RegisterServer(s *grpc.Server) {
	// 服务注册
	read.RegisterReadServer(s, &ReadController{})
	//Tracer1 = tracer
}
