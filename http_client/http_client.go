package http_client

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

// http客户端,默认为get请求
func HttpClient(url string, tracer opentracing.Tracer, clientSpan opentracing.Span) {
	req, _ := http.NewRequest("GET", url, nil)
	ext.SpanKindRPCClient.Set(clientSpan)
	ext.HTTPUrl.Set(clientSpan, url)
	ext.HTTPMethod.Set(clientSpan, "GET")
	_ = tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	_, _ = http.DefaultClient.Do(req)
}
