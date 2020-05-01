package request

import (
	"context"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/go-liam/tracing/jaeger"
	"io/ioutil"
	"net/http"
	"time"
)

/*
在2个方法是独立的
*/

func HttpGet(url string, ctx context.Context, openTrace bool) (string, error) {
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
		jaeger.Server.HttpTraceRequestInject(&ctx, &req.Header)
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

func HttpGetGin(url string, c *gin.Context, openTrace bool) (string, error) {

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
	//[TRACE] 注入 tracer 传输
	if openTrace {
		jaeger.Server.HttpTracerGinRequestInject(c, &req.Header)
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
