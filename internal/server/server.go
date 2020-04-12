package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	//pb "github.com/maddevsio/grpc-rest-api-example/pb"
	pb "gitlab.com/andrsolo21/courseproject/internal/consignment"

	"google.golang.org/grpc"
)

// Greeter ...
type Greeter struct {
	wg sync.WaitGroup
}

// New creates new server greeter
func New() *Greeter {
	return &Greeter{}
}

// Start starts server
func (g *Greeter) Start() {
	g.wg.Add(1)
	go func() {
		log.Fatal(g.startGRPC())
		g.wg.Done()
	}()
	g.wg.Add(1)
	go func() {
		log.Fatal(g.startREST())
		g.wg.Done()
	}()
}
func (g *Greeter) startGRPC() error {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterYourServiceServer(grpcServer, g)
	_ = grpcServer.Serve(lis)
	return nil
}
func (g *Greeter) startREST() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterYourServiceHandlerFromEndpoint(ctx, mux, ":8080", opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8090", mux)
}

// SayHello says hello
func (g *Greeter) GetUser(ctx context.Context, r *pb.StringMessage) (*pb.StringMessage, error) {
	return &pb.StringMessage{
		Id: fmt.Sprintf("Hello, %s!", r.Id),
	}, nil
}
