package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/y-okubo/grpc-jwt/awesome"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

func main() {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		addr  = fs.String("grpc.addr", ":8002", "Address for gRPC server")
		token = fs.String("grpc.token", "test", "JWT used to gRPC calls")
	)
	flag.Usage = fs.Usage
	fs.Parse(os.Args[1:])

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	// create client and call
	c := awesome.NewAwesomeClient(conn)

	// create context with JWT
	md := metadata.Pairs("Authorization", *token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var header, trailer metadata.MD
	res, err := c.Echo(ctx, &awesome.EchoRequest{Ping: "Hello"}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Println(err)
	}

	log.Printf("Result: %v\n", res.Pong)
}