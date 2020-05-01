package model

// 属性
type TraceModel struct {
	IsOpen       bool // 开关
	HostPort     string //  "127.0.0.1:6831"
	SamplerType  string  //固定采样
	SamplerParam float64 //1=全采样、0=不采样, 建议值 0.01
	LogSpans     bool    // true 打印日志
}