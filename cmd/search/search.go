package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"Presto.Sonata/lib/search"
	"Presto.Sonata/search/protos"
)

var (
	grpcReady = sync.WaitGroup{}
)

func runGrpcService(ctx context.Context, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen %v", err)
	}
	server := grpc.NewServer()
	if service, err := search.NewService(); err != nil {
		return err
	} else {
		searchpb.RegisterSearchServiceServer(server, service)
	}

	fmt.Printf("gRPC server is listening on %d", port)
	grpcReady.Done()
	return server.Serve(lis)
}

func runHTTPServer(ctx context.Context, port int) error {
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := searchpb.RegisterSearchServiceHandlerFromEndpoint(ctx, gwmux, "localhost:8080", opts)
	if err != nil {
		return fmt.Errorf("failed to register grpc handler %v", err)
	}

	//mux := http.NewServeMux()
	//mux.Handle("/presto/api/", gwmux)

	fmt.Printf("listening on http port %d", port)
	return http.ListenAndServe(fmt.Sprintf("localhost:%d", port), gwmux)
}

func main() {

	flag.Parse()

	ctx := context.Background()

	grpcReady.Add(1)
	go func() {
		if err := runGrpcService(ctx, 8080); err != nil {
			fmt.Printf("failed to run grpc server %v\n", err)
		}
	}()

	grpcReady.Wait()
	if err := runHTTPServer(ctx, 8081); err != nil {
		fmt.Printf("failed to run HTTP server %v\n", err)
	}
}
