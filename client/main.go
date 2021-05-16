package main

import (
	"context"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	helloworld "latihan_consul/proto_files"
	"log"
	"time"
)

const (
	address     = "consul://192.168.22.95:8500/greeter"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := helloworld.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &helloworld.HelloRequest{Name: "Richie", Age: 30})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetResponse())
}