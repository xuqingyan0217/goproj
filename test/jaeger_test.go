package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"testing"

	jaegerCfg "github.com/uber/jaeger-client-go/config"
)

func Test_Jaeger(t *testing.T) {
	// 配置信息
	cfg := jaegerCfg.Configuration{
		Sampler:             &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter:            &jaegerCfg.ReporterConfig{
			LogSpans:                   true,
			CollectorEndpoint:          fmt.Sprintf("http://%s/api/traces", "10.0.0.17:14268"),
		},
	}
	// 创建客户端，设置到了全局Opentracing中
	Jaeger, err := cfg.InitGlobalTracer("client test", jaegerCfg.Logger(jaeger.StdLogger))
	if err != nil {
		t.Log(err)
		return
	}
	defer Jaeger.Close()
	// 任务执行
	// 获取tracer，也就是上面设置到全局里面的Jaeger
	tracer := opentracing.GlobalTracer()
	// 任务节点span
	parentSpan := tracer.StartSpan("A")
	// 将其刷新到服务上去
	defer parentSpan.Finish()

	B(tracer, parentSpan)
}

func B(tracer opentracing.Tracer, parentSpan opentracing.Span) {
	// 创建一个子集的span
	childSpan := tracer.StartSpan("B", opentracing.ChildOf(parentSpan.Context()))
	defer childSpan.Finish()
}