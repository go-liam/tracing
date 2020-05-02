package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	trace2 "github.com/go-liam/tracing"
	"github.com/go-liam/tracing/example/pkg/config"
	"github.com/go-liam/tracing/example/pkg/request"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	part       = config.PortRest2 //":7422"
	serverName = "liam-server-rest2"
)

func main() {
	println("run api part ", part)
	engine := gin.New()
	// 设置路由
	SetupRouter(engine)
	engine.Run(part)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
}

// SetupRouter :
func SetupRouter(engine *gin.Engine) {
	//设置路由中间件
	engine.Use(trace2.Server.NewGinMiddlewareHandle(serverName))
	engine.GET("/", Index)
	engine.GET("/api/trace", JaegerTest)

	//404
	engine.NoRoute(func(c *gin.Context) {
		c.String(404, "请求方法不存在 ")
	})
}

// Index :
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello,It works.index ")
}

// JaegerTest : get gin http
func JaegerTest(c *gin.Context) {
	// 调用 HTTP 服务
	resHttpGet := ""
	v, err := request.HttpGetGin(config.UrlRest1+"/", c, true)
	if err == nil {
		println("getBack", v)
		resHttpGet = "[HttpGetOk]"
	}
	c.String(http.StatusOK, fmt.Sprintf("[API] %s : %s ", serverName, resHttpGet))
}
