package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/y-okubo/grpc-jwt/awesome"
	"github.com/y-okubo/grpc-jwt/user"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	fs := flag.NewFlagSet("", flag.ExitOnError)
	var (
		addr  = fs.String("grpc.addr", ":7830", "Address for gRPC server")
		token = fs.String("grpc.token", "test", "JWT used to gRPC calls")
	)
	flag.Usage = fs.Usage
	fs.Parse(os.Args[1:])

	// Create the client TLS credentials
	var cert = "server.crt"
	creds, err := NewClientTLSFromFileSkipVerify(cert, "")
	if err != nil {
		log.Println(err)
	}

	conn, err := grpc.Dial(
		*addr,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	// create client and call
	c := awesome.NewAwesomeClient(conn)

	// Create JWT
	token = user.Authenticate("name", "pass")

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

func NewClientTLSFromFileSkipVerify(certFile, serverNameOverride string) (credentials.TransportCredentials, error) {
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return credentials.NewTLS(&tls.Config{
		ServerName:         serverNameOverride,
		RootCAs:            cp,
		InsecureSkipVerify: true,
	}), nil
}
