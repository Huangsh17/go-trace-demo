package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

const JaegerHostPort = "127.0.0.1:6831" // 发送追踪数据到jaeger

func InitTracer(serverName string) (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: serverName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	sender, err := jaeger.NewUDPTransport(JaegerHostPort, 0)
	if err != nil {
		return nil, nil, err
	}
	reporter := jaeger.NewRemoteReporter(sender)
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)
	//tracer = opentracing.GlobalTracer()
	return tracer, closer, err
}

// 获取从上游服务拿到子span,针对http之间的调用
func FromHttpGetSpan(tracer opentracing.Tracer, c *gin.Context) opentracing.Span {
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
	span := tracer.StartSpan(c.Request.Host, ext.RPCServerOption(spanCtx))
	return span
}
