package grpc_client

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go-trace-demo/middleware"
	pb "go-trace-demo/proto"
	"google.golang.org/grpc"
)

func GrpcClient(tracer opentracing.Tracer, span opentracing.Span) {
	middleware.Span = span
	middleware.Tracer = tracer
	dialOption := grpc.WithUnaryInterceptor(middleware.JaegerGrpcClientInterceptor)
	conn, err := grpc.Dial(":8972", grpc.WithInsecure(), dialOption)
	if err != nil {
		fmt.Printf("faild to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	// 调用服务端的SayHello
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "q1mi"})
	if err != nil {
		fmt.Printf("could not greet: %v", err)
	}
	fmt.Printf("Greeting: %s !\n", r.Message)
}
