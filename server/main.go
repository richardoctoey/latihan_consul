package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	helloworld "latihan_consul/proto_files"
	"log"
	"net"
	"os"
)

const (
	port = ":50051"
)

var node string
var ipaddr string

// server is used to implement helloworld.GreeterServer.
type server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	log.Printf("Received: %v", req)
	return &helloworld.HelloReply{Response: fmt.Sprintf("Hello %s %d from %s", req.Name, req.Age, node)}, nil
}

func (s *server) Health(ctx context.Context, req *helloworld.HealthRequest) (*helloworld.HealthReply, error) {
	log.Printf("Received: %v", req)
	return &helloworld.HealthReply{Message: "pong"}, nil
}

func init() {
	if len(os.Args) > 2 {
		node = os.Args[1]
		ipaddr = os.Args[2]
	} else {
		node = "N1"
		ipaddr = "192.168.22.12"
	}
}

func registerConsul() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	grpcService := &api.AgentServiceRegistration{
		ID:                "greeter",
		Name:              "greeter",
		Port:              50051,
		Address:           ipaddr,
		TaggedAddresses:   nil,
		EnableTagOverride: false,
		Meta:              nil,
		Weights:           nil,
		Check: &api.AgentServiceCheck{
			TCP:                            fmt.Sprintf("%s%s", ipaddr, port),
			Timeout:                        "2s",
			Interval:                       "4s",
			DeregisterCriticalServiceAfter: "15s",
		},
	}
	err = client.Agent().ServiceRegister(grpcService)
	if err != nil {
		log.Fatal("Error", err)
	}
	fmt.Println("Consul registered..")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	registerConsul()
	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
