package config

// 属性
type TraceConfig struct {
	IsOpen       bool    // 开关
	HostPort     string  //  "127.0.0.1:6831"
	SamplerType  string  //固定采样
	SamplerParam float64 //1=全采样、0=不采样, 建议值 0.01
	LogSpans     bool    // true 打印日志
}

var DefaultConfig *TraceConfig

func init() {
	DefaultConfig = &TraceConfig{IsOpen: true, HostPort: "127.0.0.1:6831", SamplerType: "const", SamplerParam: 0.01, LogSpans: true}
}
