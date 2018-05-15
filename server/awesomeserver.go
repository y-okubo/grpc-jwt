package main

import (
	"github.com/y-okubo/grpc-jwt/awesome"
	"golang.org/x/net/context"
)

// AwesomeServer is an awesome server.
type AwesomeServer struct{}

// Echo returns the sent letters.
func (s AwesomeServer) Echo(ctx context.Context, in *awesome.EchoRequest) (*awesome.EchoResponse, error) {
	return &awesome.EchoResponse{Pong: in.Ping + ", OK"}, nil
}
