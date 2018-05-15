package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/y-okubo/grpc-jwt/awesome"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		addr = fs.String("grpc.addr", ":8002", "Address for gRPC server")
	)
	flag.Usage = fs.Usage
	fs.Parse(os.Args[1:])

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Println(err)
		return
	}

	// Register stream and unary interceptor.
	s := grpc.NewServer(
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
	)
	awesome.RegisterAwesomeServer(s, AwesomeServer{})

	s.Serve(l)
}

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := authorize(stream.Context())
	if err != nil {
		log.Println(err)
		return err
	}
	return handler(srv, stream)
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := authorize(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return handler(ctx, req)
}

func authorize(ctx context.Context) error {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		token := md["authorization"]
		log.Printf("token: %v", token[0])
		return nil
	}
	return status.Error(codes.Internal, "")
}
