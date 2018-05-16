package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/y-okubo/grpc-jwt/user"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/y-okubo/grpc-jwt/awesome"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

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
		tokenstr := md["authorization"][0]
		log.Printf("token string: %v\n", tokenstr)

		bin, err := ioutil.ReadFile("rsa.pub")
		if err != nil {
			log.Println(err)
			return status.Error(codes.Unauthenticated, "")
		}

		key, err := jwt.ParseRSAPublicKeyFromPEM(bin)
		if err != nil {
			log.Println(err)
			return status.Error(codes.Unauthenticated, "")
		}

		// With the Parse method, claims is obtained as a map.
		token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil {
			log.Println(err)
			return status.Error(codes.Unauthenticated, "")
		}

		log.Println(token.Claims, err)

		// Convert claims directly to structure
		u := user.User{}
		token, err = jwt.ParseWithClaims(tokenstr, &u, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil {
			log.Println(err)
			return status.Error(codes.Unauthenticated, "")
		}

		log.Println(token.Valid, u, err)

		return nil
	}

	return status.Error(codes.Internal, "")
}
