package config

const (
	// data
	JaegerHostPort = "127.0.0.1:6831"

	PortGateWay     = ":7410"
	PortRest1       = ":7421"
	PortRest2       = ":7422"
	PortGrpc1Listen = ":9901"
	PortGrpc2Read   = ":9902"
	PortGrpc3Speak  = ":9903"
	PortGrpc4Write  = ":9904"

	// http
	UrlGateWay = "http://localhost" + PortGateWay
	UrlRest1   = "http://localhost" + PortRest1
	UrlRest2   = "http://localhost" + PortRest2
	//grpc
	UrlGrpc1Listen = "localhost" + PortGrpc1Listen //listen
	UrlGrpc2Read   = "localhost" + PortGrpc2Read   // read
	UrlGrpc3Speak  = "localhost" + PortGrpc3Speak  // speak
	UrlGrpc4Write  = "localhost" + PortGrpc4Write  // write

)
