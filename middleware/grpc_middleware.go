package middleware

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var Span opentracing.Span
var Tracer opentracing.Tracer

/*
服务端拦截器和相关函数
*/
func JaegerGrpcServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//从context中获取metadata。md.(type) == map[string][]string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		//如果对metadata进行修改，那么需要用拷贝的副本进行修改。（FromIncomingContext的注释）
		md = md.Copy()
	}
	carrier := TextMapReader{md}
	tracer := opentracing.GlobalTracer()
	spanContext, e := tracer.Extract(opentracing.TextMap, carrier)
	if e != nil {
		fmt.Println("Extract err:", e)
	}
	span := tracer.StartSpan("HEEP", opentracing.ChildOf(spanContext))
	defer span.Finish()
	ctx = opentracing.ContextWithSpan(ctx, span)
	return handler(ctx, req)
}

type TextMapReader struct {
	metadata.MD
}

//读取metadata中的span信息
func (t TextMapReader) ForeachKey(handler func(key, val string) error) error { //不能是指针
	for key, val := range t.MD {
		for _, v := range val {
			if err := handler(key, v); err != nil {
				return err
			}
		}
	}
	return nil
}

/*
客户端拦截器和相关函数
*/

type TextMapWriter struct {
	metadata.MD
}

//重写TextMapWriter的Set方法，我们需要将carrier中的数据写入到metadata中，这样grpc才会携带。
func (t TextMapWriter) Set(key, val string) {
	//key = strings.ToLower(key)
	t.MD[key] = append(t.MD[key], val)

}

func JaegerGrpcClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		//如果对metadata进行修改，那么需要用拷贝的副本进行修改。（FromIncomingContext的注释）
		md = md.Copy()
	}
	//定义一个carrier，下面的Inject注入数据需要用到。carrier.(type) == map[string]string
	//carrier := opentracing.TextMapCarrier{}
	carrier := TextMapWriter{md}
	//将span的context信息注入到carrier中
	e := Tracer.Inject(Span.Context(), opentracing.TextMap, carrier)
	if e != nil {
		fmt.Println("tracer Inject err,", e)
	}
	//创建一个新的context，把metadata附带上
	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
