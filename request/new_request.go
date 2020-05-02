package request

import (
	"fmt"
	"github.com/go-liam/tracing/jaeger"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"io"
	"io/ioutil"
	"net/http"
)


func TraceRequest(method string,  url string, body io.Reader,sv *jaeger.SvJeager ,serverName string) ([]byte,error) {
	// nethttp.Transport from go-stdlib will do the tracing
	c := &http.Client{Transport: &nethttp.Transport{}}
	req, err := http.NewRequest(
		method ,//"GET",
		url,//"http://localhost:7201/",
		body ,//nil,
	)
	// trace
	span,req3 := sv.HttpTracerRequestInject(req,serverName)
	if err != nil {
		sv.OnError(span, err)
		return nil, err
	}
	res, err := c.Do(req3)
	if err != nil {
		sv.OnError(span, err)
		return nil, err
	}
	defer res.Body.Close()
	body2, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sv.OnError(span, err)
		return nil,err
	}
	fmt.Printf("Received 222: %s\n", string(body2))
	return body2,nil
}
