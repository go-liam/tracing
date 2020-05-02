package main

import (
	"context"
	"crypto/tls"
	"github.com/go-liam/tracing"
	"github.com/go-liam/tracing/example/pkg/config"
	"github.com/go-liam/tracing/example/proto/listen"
	"github.com/go-liam/tracing/example/proto/read"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	var s *grpc.Server
	// 注入
	if trace.Server.Config().IsOpen {
		s = grpc.NewServer(trace.Server.NewGRPCServerOption("ServiceName")) // 注入trace
	} else {
		s = grpc.NewServer()
	}
	// 服务注册
	read.RegisterReadServer(s, &ReadController{})
	println("end")
}

type ReadController struct{}

func (s *ReadController) ReadData(ctx context.Context, in *read.Request) (*read.Response, error) {
	// 调用 HTTP 服务
	resHttpGet := ""
	_, err := httpGet(config.UrlRest2, ctx, true)
	if err == nil {
		resHttpGet = "[HttpGetOk]"
	}
	//grpc
	var conn *grpc.ClientConn
	var err1 error
	if trace.Server.Config().IsOpen {
		conn, err1 = grpc.Dial(config.UrlGrpc1Listen, grpc.WithInsecure(), trace.Server.CreateGRPCConnectOpts(&ctx)) // 注入
	} else {
		conn, err1 = grpc.Dial(config.UrlGrpc1Listen, grpc.WithInsecure())
	}
	if err1 != nil {
		log.Println("grpc conn err:", err1)
	}
	grpcListenClient := listen.NewListenClient(conn)
	resListen, _ := grpcListenClient.ListenData(context.Background(), &listen.Request{Name: "listen"})
	log.Printf("grpc:%+v\n", resListen)
	return &read.Response{Message: resHttpGet}, nil
}

func httpGet(url string, ctx context.Context, openTrace bool) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Second * 5, //默认5秒超时时间
		Transport: tr,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	// [TRACE] trace: 注入
	if openTrace {
		trace.Server.HttpTraceRequestInject(&ctx, &req.Header)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	return string(content), err
}
