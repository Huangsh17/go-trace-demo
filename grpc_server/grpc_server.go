package grpc_server

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go-trace-demo/middleware"
	pb "go-trace-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello" + in.Name}, nil
}

func GrpcServer() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	tracer := opentracing.GlobalTracer()
	s := grpc.NewServer(serverOption(tracer)) // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{})    // 在gRPC服务端注册服务
	reflection.Register(s)                    //在给定的gRPC服务器上注册服务器反射服务
	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们。
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}

func serverOption(tracer opentracing.Tracer) grpc.ServerOption {
	return grpc.UnaryInterceptor(middleware.JaegerGrpcServerInterceptor)
}
